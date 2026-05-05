package metrics

import (
	"InnerG/config"
	"InnerG/pkg/logger"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

type QpsCounter struct {
	c *redis.Client
}

// 另启redis客户端
func NewQpsCounter() *QpsCounter {
	rConfig := config.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     rConfig.Addr,
		Username: rConfig.Username,
		Password: rConfig.Password,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	return &QpsCounter{c: client}
}

func (qc *QpsCounter) Increment() {
	index := time.Now().Unix()
	key := fmt.Sprintf("qps:%d", index)
	minuteKey := fmt.Sprintf("qps:minute:%d", index/(5*60))

	pipe := qc.c.Pipeline()
	pipe.SetNX(context.Background(), key, 0, 5*time.Minute)
	pipe.SetNX(context.Background(), minuteKey, 0, 25*time.Minute)
	pipe.Incr(context.Background(), key)
	pipe.Incr(context.Background(), minuteKey)

	_, err := pipe.Exec(context.Background())
	if err != nil {
		logger.Log.Errorf("Failed to increment QPS counter: %v", err)
	}
}

func (qc *QpsCounter) GetCurrentQPS() int64 {
	return qc.GetTargetSecondQPS(time.Now().Unix())
}

func (qc *QpsCounter) GetTargetSecondQPS(time int64) int64 {
	key := fmt.Sprintf("qps:%d", time)
	if qc.c.Exists(context.Background(), key).Val() != 1 {
		return 0
	}
	val, err := qc.c.Get(context.Background(), key).Result()
	if err != nil {
		logger.Log.Errorf("Failed to get current QPS counter: %v", err)
		return 0
	}
	res, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		logger.Log.Errorf("Failed to parse current QPS counter: %v", err)
	}
	return res
}

// 获取近5分钟的QPS总和
func (qc *QpsCounter) GetMinuteQps(time int64) int64 {
	MinuteKey := fmt.Sprintf("qps:minute:%d", time/(5*60))
	if qc.c.Exists(context.Background(), MinuteKey).Val() != 1 {
		return 0
	}
	val, err := qc.c.Get(context.Background(), MinuteKey).Result()
	if err != nil {
		logger.Log.Errorf("Failed to get current QPS counter: %v", err)
		return 0
	}
	res, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		logger.Log.Errorf("Failed to parse current QPS counter: %v", err)
	}
	return res
}
