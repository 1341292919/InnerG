package cache

import (
	_interface "InnerG/dao/interface"
	"InnerG/pkg/constants"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type userCache struct {
	client *redis.Client
}

func NewUserCache(cache *redis.Client) _interface.UserCache {
	return &userCache{
		client: cache,
	}
}
func (ca *userCache) IsKeyExist(ctx context.Context, key string) bool {
	return ca.client.Exists(ctx, key).Val() == 1
}
func (ca *userCache) SetEmailCode(ctx context.Context, key string, code string) error {
	if err := _Ca.Set(ctx, key, code, constants.EmailCodeKeyExpire).Err(); err != nil {
		return fmt.Errorf("SetEmailCode: Set cache failed: %w", err)
	}
	return nil
}

func (ca *userCache) GetEmailCode(ctx context.Context, key string) (string, error) {
	code, err := ca.client.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("GetEmailCode: %w", err)
	}
	return code, nil
}
