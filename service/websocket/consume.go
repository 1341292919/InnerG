package websocket

import (
	"InnerG/dao"
	"InnerG/dao/db/model"
	mq "InnerG/dao/rabbitmq"
	"InnerG/pkg/constants"
	"InnerG/pkg/logger"
	"InnerG/service/websocket/message"
	"InnerG/types"
	"context"
	"fmt"
	"sync"

	"github.com/goccy/go-json"
	"github.com/wagslane/go-rabbitmq"
)

var WsConsume *WebsocketConsume
var WebSocketConsumeOnce sync.Once
var TaskMapping map[string]mq.Handler

type WebsocketConsume struct {
	consumerList            []*mq.Consumer
	offlineConsumeNum       int
	storeConsumeNum         int
	friendRequestConsumeNum int
}

func NewWebsocketConsume() *WebsocketConsume {
	WebSocketConsumeOnce.Do(func() {
		WsConsume = &WebsocketConsume{
			offlineConsumeNum:       constants.OfflineConsumerNum,
			storeConsumeNum:         constants.StoreConsumerNum,
			friendRequestConsumeNum: constants.FriendRequestConsumerNum,
		}
		WsConsume.consumerList = make([]*mq.Consumer,
			WsConsume.offlineConsumeNum+
				WsConsume.storeConsumeNum+
				WsConsume.friendRequestConsumeNum,
		)
		for i := 0; i < WsConsume.offlineConsumeNum; i++ {
			c, err := mq.NewConsumer(
				constants.WebsocketService,
				constants.OfflineConsumeQueueTopic,
				constants.OfflineMessageTopic)
			if err != nil {
				panic(fmt.Errorf("Init WebsocketConsume offline[%d]: %w", i, err))
			}
			WsConsume.consumerList[i] = c
		}
		for i := WsConsume.offlineConsumeNum; i < WsConsume.offlineConsumeNum+WsConsume.storeConsumeNum; i++ {
			c, err := mq.NewConsumer(
				constants.WebsocketService,
				constants.StoreConsumeQueueTopic,
				constants.StoreMessageTopic)
			if err != nil {
				panic(fmt.Errorf("Init WebsocketConsume store[%d]: %w", i, err))
			}
			WsConsume.consumerList[i] = c
		}
		start := WsConsume.offlineConsumeNum + WsConsume.storeConsumeNum
		for i := start; i < start+WsConsume.friendRequestConsumeNum; i++ {
			c, err := mq.NewConsumer(
				constants.WebsocketService,
				constants.FriendRequestConsumeQueueTopic,
				constants.FriendRequestMessageTopic)
			if err != nil {
				panic(fmt.Errorf("Init WebsocketConsume friendRequest[%d]: %w", i-start, err))
			}
			WsConsume.consumerList[i] = c
		}
		TaskMapping = map[string]mq.Handler{
			constants.OfflineConsumeQueueTopic:       OfflineMessageHandler,
			constants.StoreConsumeQueueTopic:         StoreMessageHandler,
			constants.FriendRequestConsumeQueueTopic: FriendRequestEventHandler,
		}
	})
	return WsConsume
}

// 执行函数
func OfflineMessageHandler(d rabbitmq.Delivery) rabbitmq.Action {
	var m message.Message
	err := json.Unmarshal(d.Body, &m)
	if err != nil {
		logger.Log.Error("OfflineMessageHandler: ", err)
		return rabbitmq.NackRequeue
	}
	websocketDao := dao.NewWebsocketDao()
	err = websocketDao.Cache.AddOfflineMessage(context.Background(), fmt.Sprintf("offlineM:%d", m.TargetID), &m)
	if err != nil {
		logger.Log.Error("OfflineMessageHandler: ", err)
		return rabbitmq.NackRequeue
	}
	return rabbitmq.Ack
}
func StoreMessageHandler(d rabbitmq.Delivery) rabbitmq.Action {
	var m message.Message
	err := json.Unmarshal(d.Body, &m)
	if err != nil {
		logger.Log.Error("StoreMessageHandler: ", err)
		return rabbitmq.NackRequeue
	}
	dbM := &model.Message{
		Content:   m.Content,
		FromUser:  m.UserID,
		ToUser:    m.TargetID,
		MsgID:     m.ID,
		Type:      m.Type,
		CreatedAt: m.CreatedAt,
		Status:    m.Status,
	}
	websocketDao := dao.NewWebsocketDao()
	err = websocketDao.Db.InsertMessage(context.Background(), dbM)
	if err != nil {
		logger.Log.Error("StoreMessageHandler: ", err)
		return rabbitmq.NackRequeue
	}
	return rabbitmq.Ack
}

func FriendRequestEventHandler(d rabbitmq.Delivery) rabbitmq.Action {
	var event types.FriendRequestEvent
	if err := json.Unmarshal(d.Body, &event); err != nil {
		logWebsocketConsumeError("FriendRequestEventHandler: ", err)
		return rabbitmq.Ack
	}
	ws := NewWebSocketSrv()
	connID := ws.manager.WithConnectionId(event.ToUser)
	if !ws.manager.IsConnected(connID) {
		return rabbitmq.Ack
	}
	conn := ws.manager.GetConnection(connID)
	if conn == nil {
		return rabbitmq.Ack
	}
	if err := conn.WriteJSONData(event); err != nil {
		logWebsocketConsumeError("FriendRequestEventHandler:WriteJSONData: ", err)
	}
	return rabbitmq.Ack
}

func logWebsocketConsumeError(v ...interface{}) {
	if logger.Log == nil {
		return
	}
	logger.Log.Error(v...)
}

func (w *WebsocketConsume) Run() {
	for i := 0; i < len(w.consumerList); i++ {
		go func() {
			err := w.consumerList[i].Run(TaskMapping[w.consumerList[i].Queue])
			if err != nil {
				logger.Log.Error("Run: ", err)
			}
		}()
	}
}
