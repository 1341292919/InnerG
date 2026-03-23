package cache

import (
	_interface "InnerG/dao/interface"
	"InnerG/pkg/constants"
	"InnerG/pkg/errno"
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
		return errno.NewErr(errno.RedisDBErrorCode, fmt.Sprintf("SetEmailCode: %v", err))
	}
	return nil
}

func (ca *userCache) GetEmailCode(ctx context.Context, key string) (string, error) {
	code, err := ca.client.Get(ctx, key).Result()
	if err != nil {
		return "", errno.NewErr(errno.RedisDBErrorCode, fmt.Sprintf("GetEmailCode: %v", err))
	}
	return code, nil
}

func (ca *userCache) BlockToken(ctx context.Context, key string) error {
	// 仅把token作为key即可
	if err := _Ca.Set(ctx, key, "", constants.AccessTokenTTL).Err(); err != nil {
		return errno.NewErr(errno.RedisDBErrorCode, fmt.Sprintf("BlockToken: %v", err))
	}
	return nil
}
