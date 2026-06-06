package _interface

import (
	"InnerG/dao/mongo/model"
	"context"
)

type ContactMongoDB interface {
	NewChatSession(ctx context.Context, session *model.ChatSession) error
	IsQuerySessionExist(ctx context.Context, sessionId string) (bool, *model.ChatSession, error)
	InsertMessageToSession(ctx context.Context, sessionId string, message []model.Message) error
	GetSessionByUserId(ctx context.Context, userId int64) ([]*model.ChatSession, int, error)
	UpdateSessionTitle(ctx context.Context, sessionId string, title string) error
	DeleteSession(ctx context.Context, sessionId string) error
	GetSessionByUserIdWithPagination(ctx context.Context, userId int64, pageNum, pageSize int) ([]*model.ChatSession, int, error)
}
