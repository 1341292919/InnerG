package _interface

import (
	"InnerG/dao/db/model"
	"context"
)

type WebSocketDB interface {
	InsertMessage(ctx context.Context, msg *model.Message) error
	GetOfflineMessages(ctx context.Context, toUser int64, limit int) ([]*model.Message, error)
	GetMessagesAfterTimestamp(ctx context.Context, user1, user2 int64, timestamp int64) ([]*model.Message, error)
	UpdateMessageStatus(ctx context.Context, msgID string, status int8) error
	BatchUpdateMessageStatus(ctx context.Context, msgIDs []string, status int8) error
}
