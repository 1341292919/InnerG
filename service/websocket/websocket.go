package websocket

import (
	"InnerG/config"
	"InnerG/dao"
	"InnerG/dao/rabbitmq"
	"InnerG/pkg/constants"
	"InnerG/pkg/ctl"
	"InnerG/pkg/logger"
	"InnerG/pkg/utils"
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
var messageSnowflake *utils.Snowflake // 消息专用雪花ID生成器
var messageSnowflakeOnce sync.Once    // 保证只初始化一次

// initMessageSnowflake 延迟初始化雪花ID生成器
func initMessageSnowflake() {
	messageSnowflakeOnce.Do(func() {
		var err error
		messageSnowflake, err = utils.NewSnowflake(
			config.Snowflake.DatancenterID,
			config.Snowflake.WorkerID,
		)
		if err != nil {
			panic(fmt.Errorf("init message snowflake: %w", err))
		}
	})
}

// GenerateMessageID 生成消息雪花ID
func GenerateMessageID() (int64, error) {
	if messageSnowflake == nil {
		initMessageSnowflake()
	}
	return messageSnowflake.NextVal()
}

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

	// 订阅用户所属的所有群组
	groupDao := dao.NewGroupDao()
	groupIDs, err := groupDao.Db.GetUserGroupIDs(ctx, uid)
	if err != nil {
		logger.Log.Errorf("Failed to get user group IDs: %v", err)
	} else if len(groupIDs) > 0 {
		ws.groupManager.SubscribeUser(uid, groupIDs)
	}

	// 注册连接
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
			// 先尝试解析为同步请求
			var syncReq SyncRequest
			if err := json.Unmarshal(body, &syncReq); err == nil && syncReq.Action == "sync" {
				// 处理游标同步请求
				resp, err := ws.SyncMessages(ctx, uid, syncReq.LastID, syncReq.Limit)
				if err != nil {
					logger.Log.Errorf("NewConnection:SyncMessages error: %v", err)
					if closeIfOverLimit(messageOverLimit) {
						return
					}
					continue
				}
				if err := conn.WriteJSONData(resp); err != nil {
					logger.Log.Errorf("NewConnection:WriteSyncResponse error: %v", err)
					if closeIfOverLimit(messageOverLimit) {
						return
					}
					continue
				}
				if closeIfOverLimit(messageOverLimit) {
					return
				}
				continue
			}

			// 否则按照普通消息处理
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

// RouteMessage 路由私聊消息（简化版：持久化 + 尽力推送）
// 从消息结构创建高频的特点 这里选择传入 message.Message 而不是 *message.Message
// GC 友好？
func (ws *WebSocketSrv) RouteMessage(m message.Message) error {
	// 1. 先持久化消息（发送到存储队列）
	m.Status = constants.MessageStatusNormal
	if err := ws.sender.SendMessage(constants.StoreMessageTopic, m); err != nil {
		logger.Log.Error("RouteMessage: send to store queue error: ", err)
		return err
	}

	// 2. 尽力而为的实时推送（如果接收方在线）
	if ws.manager.IsConnected(ws.manager.WithConnectionId(m.TargetID)) {
		c := ws.manager.GetConnection(ws.manager.WithConnectionId(m.TargetID))
		if c != nil {
			content, err := m.JsonContent()
			if err != nil {
				logger.Log.Error("RouteMessage:JsonContent: ", err)
				// 持久化成功，推送失败不影响整体流程
				return nil
			}
			err = c.WriteData(content)
			if err != nil {
				logger.Log.Error("RouteMessage:WriteMessage: ", err)
				// 推送失败，客户端会通过游标同步补齐
			}
		}
	}
	// 不再发送离线消息到队列，客户端重连时通过游标同步获取
	return nil
}
