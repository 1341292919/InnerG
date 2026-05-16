package cache

import (
	_interface "InnerG/dao/interface"
	"InnerG/pkg/constants"
	"InnerG/pkg/errno"
	"InnerG/service/websocket/message"
	"context"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
)

type websocketCache struct {
	client *redis.Client
}

func NewWebsocketCache(cache *redis.Client) _interface.WebSocketCache {
	return &websocketCache{
		client: cache,
	}
}

// IsKeyExist 检查 key 是否存在
func (wc *websocketCache) IsKeyExist(ctx context.Context, key string) bool {
	return wc.client.Exists(ctx, key).Val() == 1
}

// AddOfflineMessage 追加离线消息到用户的 Set 结构中
func (wc *websocketCache) AddOfflineMessage(ctx context.Context, key string, m *message.Message) error {
	messageJson, err := json.Marshal(m)
	if err != nil {
		return errno.NewErr(errno.RedisDBErrorCode, fmt.Sprintf("AddOfflineMessage Marshal: %v", err))
	}

	// 使用 Pipeline 批量执行
	pipe := wc.client.Pipeline()
	pipe.SAdd(ctx, key, messageJson)
	pipe.Expire(ctx, key, constants.WebsocketKeyExpire)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return errno.NewErr(errno.RedisDBErrorCode, fmt.Sprintf("AddOfflineMessage Pipeline Exec: %v", err))
	}

	return nil
}

// GetOfflineMessages 获取用户的所有离线消息
func (wc *websocketCache) GetOfflineMessages(ctx context.Context, key string) ([]*message.Message, error) {
	messagesJson, err := wc.client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, errno.NewErr(errno.RedisDBErrorCode, fmt.Sprintf("GetOfflineMessages SMembers: %v", err))
	}

	if len(messagesJson) == 0 {
		return []*message.Message{}, nil
	}

	messages := make([]*message.Message, 0, len(messagesJson))
	for _, msgJson := range messagesJson {
		var msg message.Message
		if err := json.Unmarshal([]byte(msgJson), &msg); err != nil {
			return nil, errno.NewErr(errno.RedisDBErrorCode, fmt.Sprintf("GetOfflineMessages Unmarshal: %v", err))
		}
		messages = append(messages, &msg)
	}

	return messages, nil
}

// DeleteOfflineMessages 删除用户的离线消息 key
func (wc *websocketCache) DeleteOfflineMessages(ctx context.Context, key string) error {
	return wc.client.Del(ctx, key).Err()
}
