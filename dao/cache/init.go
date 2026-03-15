package cache

import (
	"InnerG/config"
	_interface "InnerG/dao/interface"
	"context"
	"github.com/redis/go-redis/v9"
)

var _Ca *redis.Client
var RedisContext = context.Background()

func InitCache() {
	rConfig := config.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     rConfig.Addr,
		Username: rConfig.Username,
		Password: rConfig.Password,
	})
	_, err := client.Ping(RedisContext).Result()
	if err != nil {
		panic(err)
	}
	_Ca = client
}

func NewRedisClient() _interface.UserCache {
	return NewUserCache(_Ca)
}
