package mongo

import (
	_interface "InnerG/dao/interface"
	"InnerG/dao/mongo/model"
	"InnerG/pkg/constants"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"log"
	"time"
)

type contactMongoDB struct {
	client *mongo.Database
}

func NewContactMongoDB(db *mongo.Database) _interface.ContactMongoDB {
	return &contactMongoDB{
		client: db,
	}
}

func (m *contactMongoDB) NewChatSession(ctx context.Context, session *model.ChatSession) error {
	_, err := m.client.Collection(constants.ChatSessionCollection).InsertOne(ctx, session)
	return err
}

func (m *contactMongoDB) IsQuerySessionExist(ctx context.Context, sessionId string) (bool, *model.ChatSession, error) {
	filter := bson.M{"sessionId": sessionId}
	var session model.ChatSession
	err := m.client.Collection(constants.ChatSessionCollection).FindOne(ctx, filter).Decode(&session)
	if err != nil {
		log.Println(err)
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil, fmt.Errorf("IsQuerySessionExist:Session not exist")
		}
		return false, nil, err
	}
	log.Print(1)
	return true, &session, nil
}

// InsertMessageToSession 用于插入会话新的聊天内容
func (m *contactMongoDB) InsertMessageToSession(ctx context.Context, sessionId string, messages []model.Message) error {
	filter := bson.M{"sessionId": sessionId}
	update := bson.M{
		"$push": bson.M{
			"messages": bson.M{
				"$each": messages, // 使用 $each 批量插入
			},
		},
		"$set": bson.M{
			"updatedAt": time.Now(),
		},
	}
	_, err := m.client.Collection(constants.ChatSessionCollection).UpdateOne(ctx, filter, update)
	return err
}
