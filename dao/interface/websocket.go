package _interface

import (
	"InnerG/dao/db/model"
	"InnerG/service/websocket/message"
	"context"
	"time"
)

type WebSocketDB interface {
	InsertMessage(ctx context.Context, msg *model.Message) error
	GetOfflineMessages(ctx context.Context, toUser int64, limit int) ([]*model.Message, error)
	GetMessagesByTimeRange(ctx context.Context, user1, user2 int64, after, before int64, pageSize, pageNum int) ([]*model.Message, int64, error)
	UpdateMessageStatus(ctx context.Context, msgID int64, status int8) error
	BatchUpdateMessageStatus(ctx context.Context, msgIDs []int64, status int8) error
	AckMessages(ctx context.Context, toUser int64, msgIDs []string) error
	GetMessagesAfterID(ctx context.Context, userID, lastID int64, limit int) ([]*model.Message, error)
}

type WebSocketCache interface {
	// IsKeyExist 检查 key 是否存在
	IsKeyExist(ctx context.Context, key string) bool

	// AddOfflineMessage 追加离线消息到用户的 Set 结构中
	AddOfflineMessage(ctx context.Context, key string, m *message.Message) error

	// GetOfflineMessages 获取用户的所有离线消息
	GetOfflineMessages(ctx context.Context, key string) ([]*message.Message, error)

	// DeleteOfflineMessages 删除用户的离线消息 key
	DeleteOfflineMessages(ctx context.Context, key string) error

	// SetTicket 保存 WebSocket 一次性连接 ticket
	SetTicket(ctx context.Context, key, userID string, ttl time.Duration) error

	// ConsumeTicket 原子读取并删除 WebSocket 一次性连接 ticket
	ConsumeTicket(ctx context.Context, key string) (string, error)
}
