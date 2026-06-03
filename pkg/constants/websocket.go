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

	OfflineMessageTopic            = "offline"
	OfflineConsumeQueueTopic       = "offline_queue"
	StoreConsumeQueueTopic         = "store_queue"
	StoreMessageTopic              = "store"
	FriendRequestMessageTopic      = "friend.request"
	FriendRequestConsumeQueueTopic = "friend_request_queue"
	FriendRequestEventType         = "friend_request"
	FriendRequestAcceptedEventType = "friend_request_accepted"
	FriendRequestAcceptedMessage   = "我通过了你的好友验证，现在我们开始聊天吧！"

	OfflineConsumerNum       = 3
	StoreConsumerNum         = 3
	FriendRequestConsumerNum = 1
)

const (
	WebsocketKeyExpire = 24 * ONE_HOUR
)

const (
	WebsocketImageOssOrigin      = "websocket/image"
	WebsocketVideoOssOrigin      = "websocket/video"
	WebsocketImageFileNamePrefix = "ws_image"
	WebsocketVideoFileNamePrefix = "ws_video"
)
