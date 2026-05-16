package _interface

import (
	"InnerG/dao/db/model"
	"InnerG/service/websocket/message"
	"context"
)

type WebSocketDB interface {
	InsertMessage(ctx context.Context, msg *model.Message) error
	GetOfflineMessages(ctx context.Context, toUser int64, limit int) ([]*model.Message, error)
	GetMessagesAfterTimestamp(ctx context.Context, user1, user2 int64, timestamp int64) ([]*model.Message, error)
	UpdateMessageStatus(ctx context.Context, msgID int64, status int8) error
	BatchUpdateMessageStatus(ctx context.Context, msgIDs []int64, status int8) error
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
}
