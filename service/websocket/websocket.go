package websocket

import (
	"InnerG/dao/rabbitmq"
	"InnerG/pkg/constants"
	"InnerG/pkg/ctl"
	"InnerG/pkg/logger"
	"InnerG/service/websocket/manager"
	"InnerG/service/websocket/message"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
)

var WebSocketSrvIns *WebSocketSrv
var WebSocketSrvOnce sync.Once

type WebSocketSrv struct {
	manager *manager.ConnectionManager
	// 消息队列
	sender *rabbitmq.Producer
}

func NewWebSocketSrv() *WebSocketSrv {
	WebSocketSrvOnce.Do(func() {
		producer, err := rabbitmq.NewProducer(constants.WebsocketService)
		if err != nil {
			panic(fmt.Errorf("Init WebSocketService : ", err.Error()))
		}
		WebSocketSrvIns = &WebSocketSrv{
			manager: manager.NewConnectionManager(),
			sender:  producer,
		}
	})
	return WebSocketSrvIns
}

func (ws *WebSocketSrv) NewConnection(ctx context.Context, connect *websocket.Conn) {
	uid := ctl.GetUserInfo(ctx).Id
	// 注册该连接
	conn := manager.UserConnection{
		UserID:     uid,
		Conn:       connect,
		LastActive: time.Now(),
	}
	ws.manager.AddConnection(&conn)
	defer ws.manager.RemoveConnection(uid)
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
	if ws.manager.IsConnected(m.TargetID) {
		c := ws.manager.GetConnection(m.TargetID)
		content, err := m.JsonContent()
		if err != nil {
			logger.Log.Error("RouteMessage:JsonContent: ", err)
			return err
		}
		err = c.Conn.WriteMessage(websocket.TextMessage, content)
		if err != nil {
			logger.Log.Error("RouteMessage:WriteMessage: ", err)
			return err
		}
		pushed = true
	}
	var err error
	// 加入消息队列
	// 离线消息
	if !pushed {
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
