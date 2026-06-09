package websocket

import (
	"InnerG/config"
	"InnerG/dao"
	"InnerG/dao/rabbitmq"
	"InnerG/pkg/constants"
	"InnerG/pkg/ctl"
	"InnerG/pkg/logger"
	rootservice "InnerG/service"
	"InnerG/service/websocket/manager"
	"InnerG/service/websocket/message"
	"context"
	"fmt"
	"sync"
	"time"

	json "github.com/goccy/go-json"
	"github.com/gorilla/websocket"
)

var WebSocketSrvIns *WebSocketSrv
var WebSocketSrvOnce sync.Once

var isFriend = func(ctx context.Context, userID, targetID int64) (bool, error) {
	return rootservice.GetFriendSrv().IsFriend(ctx, userID, targetID)
}

type WebSocketSrv struct {
	manager      *manager.ConnectionManager
	groupManager *GroupManager
	sender       *rabbitmq.Producer
}

func NewWebSocketSrv() *WebSocketSrv {
	if WebSocketSrvIns != nil {
		return WebSocketSrvIns
	}
	WebSocketSrvOnce.Do(func() {
		producer, err := rabbitmq.NewProducer(constants.WebsocketService)
		if err != nil {
			panic(fmt.Errorf("init WebSocketService: %w", err))
		}
		WebSocketSrvIns = &WebSocketSrv{
			manager:      manager.NewConnectionManager(config.Service.WebsocketShardNum),
			groupManager: NewGroupManager(),
			sender:       producer,
		}
	})
	return WebSocketSrvIns
}

func canChat(ctx context.Context, userID, targetID int64) error {
	ok, err := isFriend(ctx, userID, targetID)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("双方不是好友")
	}
	return nil
}

func shouldAcceptMessageType(messageType int8) bool {
	return constants.IsValidMessageType(messageType)
}

func (ws *WebSocketSrv) NewConnection(ctx context.Context, connect *websocket.Conn) {
	uid := ctl.GetUserInfo(ctx).Id
	conn := manager.UserConnection{
		UserID:     uid,
		Conn:       connect,
		LastActive: time.Now(),
	}

	// 先推送离线消息，再注册连接，确保消息顺序：离线 -> 实时
	websocketDao := dao.NewWebsocketDao()
	mList, err := websocketDao.Db.GetOfflineMessages(ctx, uid, constants.OnceOffMessagePushNum)
	if err != nil {
		logger.Log.Error("GetOfflineMessages: ", err)
	}
	if len(mList) > 0 {
		err = conn.WriteJSONData(mList)
		if err != nil {
			logger.Log.Error("Failed to push offline messages: ", err)
			return
		}
		ids := make([]int64, len(mList))
		for i := range mList {
			ids[i] = mList[i].ID
		}
		err = websocketDao.Db.BatchUpdateMessageStatus(ctx, ids, constants.MessagePushedStatus)
		if err != nil {
			logger.Log.Error("Failed to push offline messages: ", err)
			return
		}
	}

	groupDao := dao.NewGroupDao()
	groupIDs, err := groupDao.Db.GetUserGroupIDs(ctx, uid)
	if err != nil {
		logger.Log.Errorf("Failed to get user group IDs: %v", err)
	} else if len(groupIDs) > 0 {
		ws.groupManager.SubscribeUser(uid, groupIDs)
	}

	// 离线消息推完后再注册，后续 RouteMessage 推送的实时消息都排在后面
	ws.manager.AddConnection(&conn)
	defer func() {
		ws.manager.RemoveConnection(&conn)
		ws.groupManager.UnsubscribeUser(uid)
	}()

	messageLimiter := newWebsocketMessageLimiter(constants.WebsocketMessageRateLimit, constants.WebsocketMessageRateBurst)
	closeIfOverLimit := func(overLimit bool) bool {
		if !overLimit {
			return false
		}
		closeRateLimitedConnection(connect, uid)
		return true
	}

	for {
		t, body, err := connect.ReadMessage()
		if err != nil {
			break
		}
		messageOverLimit := !messageLimiter.Allow()
		switch t {
		case websocket.TextMessage:
			var m message.Message
			err = json.Unmarshal(body, &m)
			if err != nil {
				logger.Log.Error("NewConnection:ReadMessage: ", err)
				if closeIfOverLimit(messageOverLimit) {
					return
				}
				continue
			}
			if m.UserID != uid {
				if closeIfOverLimit(messageOverLimit) {
					return
				}
				continue
			}
			if !shouldAcceptMessageType(m.Type) {
				logger.Log.Error("NewConnection:invalid message type: ", m.Type)
				if closeIfOverLimit(messageOverLimit) {
					return
				}
				continue
			}

			if m.ChatType == constants.ChatTypeGroup {
				groupMsg := &message.GroupMessage{
					ID:        m.ID,
					UserID:    m.UserID,
					GroupID:   m.TargetID,
					Content:   m.Content,
					Type:      m.Type,
					Status:    m.Status,
					ChatType:  m.ChatType,
					CreatedAt: m.CreatedAt,
				}
				err = ws.RouteGroupMessage(ctx, groupMsg)
				if err != nil {
					logger.Log.Error("NewConnection:RouteGroupMessage: ", err)
					if closeIfOverLimit(messageOverLimit) {
						return
					}
					continue
				}
			} else {
				if err = canChat(ctx, m.UserID, m.TargetID); err != nil {
					logger.Log.Error("NewConnection:canChat: ", err)
					if closeIfOverLimit(messageOverLimit) {
						return
					}
					continue
				}
				err = ws.RouteMessage(m)
				if err != nil {
					return
				}
			}

			if closeIfOverLimit(messageOverLimit) {
				return
			}
		default:
			logger.Log.Error("NewConnection:ReadMessage: ", err)
			if closeIfOverLimit(messageOverLimit) {
				return
			}
		}
	}
}

func closeRateLimitedConnection(connect *websocket.Conn, uid int64) {
	if logger.Log != nil {
		logger.Log.Errorf("websocket message rate limit exceeded: user_id=%d remote_addr=%s limit=%d/s burst=%d", uid, connect.RemoteAddr().String(), constants.WebsocketMessageRateLimit, constants.WebsocketMessageRateBurst)
	}
	_ = connect.WriteControl(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "message rate limit exceeded"),
		time.Now().Add(time.Second),
	)
}

// 从消息结构创建高频的特点 这里选择传入 message.Message 而不是 *message.Message
// GC 友好？
func (ws *WebSocketSrv) RouteMessage(m message.Message) error {
	pushed := false
	// 在线直接推送
	if ws.manager.IsConnected(ws.manager.WithConnectionId(m.TargetID)) {
		c := ws.manager.GetConnection(ws.manager.WithConnectionId(m.TargetID))
		content, err := m.JsonContent()
		if err != nil {
			logger.Log.Error("RouteMessage:JsonContent: ", err)
			return err
		}
		err = c.WriteData(content)
		if err != nil {
			logger.Log.Error("RouteMessage:WriteMessage: ", err)
			return err
		}
		pushed = true
	}
	var err error
	m.Status = constants.MessagePushedStatus
	// 加入消息队列
	// 离线消息
	if !pushed {
		m.Status = constants.MessageUnPushedStatus
		if err = ws.sender.SendMessage(constants.OfflineMessageTopic, m); err != nil {
			logger.Log.Error("RouteMessage:TaskQueue:SendMessage: ", err)
			return err
		}
	}
	//  持久化消息
	if err = ws.sender.SendMessage(constants.StoreMessageTopic, m); err != nil {
		logger.Log.Error("RouteMessage:TaskQueue:SendMessage: ", err)
		return err
	}
	return nil
}
