package websocket

import (
	"InnerG/pkg/ctl"
	"InnerG/pkg/logger"
	"InnerG/service/websocket/manager"
	"InnerG/service/websocket/message"
	"context"
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
}

func NewWebSocketSrv() *WebSocketSrv {
	WebSocketSrvOnce.Do(func() {
		WebSocketSrvIns = &WebSocketSrv{
			manager: manager.NewConnectionManager(),
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
			ws.RouteMessage(m)
		default:
			logger.Log.Error("NewConnection:ReadMessage: ", err)
		}
	}
	// 连接中断下线
	ws.manager.RemoveConnection(uid)
}

// 从消息结构创建高频的特点 这里选择传入 message.Message 而不是 *message.Message
// GC 友好？
func (ws *WebSocketSrv) RouteMessage(m message.Message) {
	//pushed := false
	if ws.manager.IsConnected(m.TargetID) {
		c := ws.manager.GetConnection(m.TargetID)
		content, err := m.JsonContent()
		if err != nil {
			logger.Log.Error("RouteMessage:JsonContent: ", err)
			return
		}
		err = c.Conn.WriteMessage(websocket.TextMessage, []byte(content))
		if err != nil {
			logger.Log.Error("RouteMessage:WriteMessage: ", err)
		}

	}
	// 加入消息队列

}
