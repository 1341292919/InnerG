package constants

const (
	OnceOffMessagePushNum = 50
)

const (
	MessagePushedStatus   = 0
	MessageUnPushedStatus = 1
	MessageRecalledStatus = 2
	MessageDeletedStatus  = 4
)

const (
	WebsocketService = "websocket"

	OfflineMessageTopic      = "offline"
	OfflineConsumeQueueTopic = "offline_queue"
	StoreConsumeQueueTopic   = "store_queue"
	StoreMessageTopic        = "store"

	OfflineConsumerNum = 3
	StoreConsumerNum   = 3
)

const (
	WebsocketKeyExpire = 24 * ONE_HOUR
)
