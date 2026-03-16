package _interface

import (
	"InnerG/dao/mongo/model"
	"context"
)

type ContactMongoDB interface {
	NewChatSession(ctx context.Context, session *model.ChatSession) error
	IsQuerySessionExist(ctx context.Context, sessionId string) (bool, *model.ChatSession, error)
	InsertMessageToSession(ctx context.Context, sessionId string, message []model.Message) error
}
