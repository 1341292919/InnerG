package websocket

import (
	"InnerG/dao"
	"InnerG/dao/db/model"
	mq "InnerG/dao/rabbitmq"
	"InnerG/pkg/constants"
	"InnerG/pkg/logger"
	"InnerG/service/websocket/message"
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
	consumerList      []*mq.Consumer
	offlineConsumeNum int
	storeConsumeNum   int
}

func NewWebsocketConsume() *WebsocketConsume {
	WebSocketConsumeOnce.Do(func() {
		WsConsume = &WebsocketConsume{
			offlineConsumeNum: constants.OfflineConsumerNum,
			storeConsumeNum:   constants.StoreConsumerNum,
		}
		WsConsume.consumerList = make([]*mq.Consumer,
			WsConsume.offlineConsumeNum+
				WsConsume.storeConsumeNum,
		)
		for i := 0; i < WsConsume.offlineConsumeNum; i++ {
			WsConsume.consumerList[i], _ = mq.NewConsumer(
				constants.WebsocketService,
				constants.OfflineConsumeQueueTopic,
				constants.OfflineMessageTopic)
		}
		for i := WsConsume.offlineConsumeNum; i < WsConsume.offlineConsumeNum+WsConsume.storeConsumeNum; i++ {
			WsConsume.consumerList[i], _ = mq.NewConsumer(
				constants.WebsocketService,
				constants.StoreConsumeQueueTopic,
				constants.StoreMessageTopic)
		}
		TaskMapping = map[string]mq.Handler{
			constants.OfflineConsumeQueueTopic: OutlineMessageHandler,
			constants.StoreConsumeQueueTopic:   StoreMessageHandler,
		}
	})
	return WsConsume
}

// 执行函数
func OutlineMessageHandler(d rabbitmq.Delivery) rabbitmq.Action {
	fmt.Println("OutlineMessageHandler:", string(d.Body))
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
		Status:    0,
	}
	websocketDao := dao.NewWebsocketDao(context.Background())
	err = websocketDao.Db.InsertMessage(context.Background(), dbM)
	if err != nil {
		logger.Log.Error("StoreMessageHandler: ", err)
		return rabbitmq.NackRequeue
	}
	return rabbitmq.Ack
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
