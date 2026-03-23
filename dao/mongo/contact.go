package mongo

import (
	_interface "InnerG/dao/interface"
	"InnerG/dao/mongo/model"
	"InnerG/pkg/constants"
	"InnerG/pkg/errno"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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
	if _, err := m.client.Collection(constants.ChatSessionCollection).InsertOne(ctx, session); err != nil {
		return errno.NewErr(errno.MongoDBErrorCode, "NewChatSession: "+err.Error())
	}
	return nil
}

func (m *contactMongoDB) IsQuerySessionExist(ctx context.Context, sessionId string) (bool, *model.ChatSession, error) {
	filter := bson.M{
		"sessionId": sessionId,
		"status": bson.M{
			"$nin": []interface{}{"0", 0},
		},
	}
	var session model.ChatSession
	err := m.client.Collection(constants.ChatSessionCollection).FindOne(ctx, filter).Decode(&session)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil, fmt.Errorf("IsQuerySessionExist:Session not exist")
		}
		return false, nil, errno.NewErr(errno.MongoDBErrorCode, "IsQuerySessionExist: "+err.Error())
	}
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
	if err != nil {
		return errno.NewErr(errno.MongoDBErrorCode, "InsertMessageToSession: "+err.Error())
	}
	return nil
}

func (m *contactMongoDB) GetSessionByUserId(ctx context.Context, userId string) ([]*model.ChatSession, int, error) {
	filter := bson.M{
		"userId": userId,
		"status": bson.M{
			"$nin": []interface{}{"0", 0},
		},
	}
	var sessionList []*model.ChatSession
	cursor, err := m.client.Collection(constants.ChatSessionCollection).Find(ctx, filter)
	defer cursor.Close(ctx)
	if err != nil {
		return nil, -1, errno.NewErr(errno.MongoDBErrorCode, "GetSessionByUserId: "+err.Error())
	}
	err = cursor.All(ctx, &sessionList)
	if err != nil {
		return nil, -1, errno.NewErr(errno.MongoDBErrorCode, "GetSessionByUserId: "+err.Error())
	}
	return sessionList, len(sessionList), nil
}

func (m *contactMongoDB) UpdateSessionTitle(ctx context.Context, sessionId string, title string) error {
	filter := bson.M{"sessionId": sessionId}
	update := bson.M{
		"$set": bson.M{
			"title":     title,
			"updatedAt": time.Now(),
		},
	}
	_, err := m.client.Collection(constants.ChatSessionCollection).UpdateOne(ctx, filter, update)
	if err != nil {
		return errno.NewErr(errno.MongoDBErrorCode, "UpdateSessionTitle: "+err.Error())
	}
	return nil
}

func (m *contactMongoDB) DeleteSession(ctx context.Context, sessionId string) error {
	filter := bson.M{"sessionId": sessionId}
	update := bson.M{
		"$set": bson.M{
			"status":    constants.CommonDeletedStatus,
			"updatedAt": time.Now(),
		},
	}
	_, err := m.client.Collection(constants.ChatSessionCollection).UpdateOne(ctx, filter, update)
	if err != nil {
		return errno.NewErr(errno.MongoDBErrorCode, "DeleteSession: "+err.Error())
	}
	return err
}

func (m *contactMongoDB) GetSessionByUserIdWithPagination(ctx context.Context, userId string, pageNum, pageSize int) ([]*model.ChatSession, int, error) {
	filter := bson.M{
		"userId": userId,
		"status": bson.M{
			"$nin": []interface{}{"0", 0},
		},
	}

	// 获取总数
	total, err := m.client.Collection(constants.ChatSessionCollection).CountDocuments(ctx, filter)
	if err != nil {
		return nil, -1, errno.NewErr(errno.MongoDBErrorCode, "GetSessionByUserIdWithPagination CountDocuments: "+err.Error())
	}

	// 计算跳过数量
	skip := int64((pageNum - 1) * pageSize)
	limit := int64(pageSize)

	var sessionList []*model.ChatSession

	// v2 版本使用 options.Find()
	opts := options.Find()
	opts.SetSort(bson.D{{Key: "updatedAt", Value: -1}})
	opts.SetSkip(skip)
	opts.SetLimit(limit)

	cursor, err := m.client.Collection(constants.ChatSessionCollection).Find(ctx, filter, opts)
	if err != nil {
		return nil, -1, errno.NewErr(errno.MongoDBErrorCode, "GetSessionByUserIdWithPagination Find: "+err.Error())
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &sessionList)
	if err != nil {
		return nil, -1, errno.NewErr(errno.MongoDBErrorCode, "GetSessionByUserIdWithPagination All: "+err.Error())

	}

	return sessionList, int(total), nil
}
