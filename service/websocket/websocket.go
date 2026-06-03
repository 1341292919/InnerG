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
	manager *manager.ConnectionManager
	// 消息队列
	sender *rabbitmq.Producer
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
			manager: manager.NewConnectionManager(config.Service.WebsocketShardNum),
			sender:  producer,
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

func (ws *WebSocketSrv) NewConnection(ctx context.Context, connect *websocket.Conn) {
	uid := ctl.GetUserInfo(ctx).Id
	conn := manager.UserConnection{
		UserID:     uid,
		Conn:       connect,
		LastActive: time.Now(),
	}

	// 先推送离线消息，再注册连接，确保消息顺序：离线 → 实时
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

	// 离线消息推完后再注册，后续 RouteMessage 推送的实时消息都排在后面
	ws.manager.AddConnection(&conn)
	defer ws.manager.RemoveConnection(&conn)

	for {
		t, body, err := connect.ReadMessage()
		if err != nil { // 连接中断
			break
		}
		//  解析消息
		switch t {
		case websocket.TextMessage:
			var m message.Message
			err = json.Unmarshal(body, &m)
			if err != nil {
				logger.Log.Error("NewConnection:ReadMessage: ", err)
				continue
			}
			// 拒绝非本人信息
			if m.UserID != uid {
				continue
			}
			if err = canChat(ctx, m.UserID, m.TargetID); err != nil {
				logger.Log.Error("NewConnection:canChat: ", err)
				continue
			}
			//  路由消息
			// -- 失败需要中断连接
			err = ws.RouteMessage(m)
			if err != nil {
				return
			}
		default:
			logger.Log.Error("NewConnection:ReadMessage: ", err)
		}
	}
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
