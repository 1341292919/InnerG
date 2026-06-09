package cache

import (
	"InnerG/pkg/constants"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type GroupCache struct {
	rdb *redis.Client
}

func NewGroupCache(rdb *redis.Client) *GroupCache {
	return &GroupCache{rdb: rdb}
}

func (g *GroupCache) SetGroupMembers(ctx context.Context, groupID int64, userIDs []int64, ttl time.Duration) error {
	key := fmt.Sprintf("%s%d", constants.GroupMembersCacheKeyPrefix, groupID)

	members := make([]interface{}, len(userIDs))
	for i, uid := range userIDs {
		members[i] = strconv.FormatInt(uid, 10)
	}

	pipe := g.rdb.Pipeline()
	pipe.Del(ctx, key)
	if len(members) > 0 {
		pipe.SAdd(ctx, key, members...)
	}
	pipe.Expire(ctx, key, ttl)

	_, err := pipe.Exec(ctx)
	return err
}

func (g *GroupCache) GetGroupMembers(ctx context.Context, groupID int64) ([]int64, error) {
	key := fmt.Sprintf("%s%d", constants.GroupMembersCacheKeyPrefix, groupID)

	members, err := g.rdb.SMembers(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	userIDs := make([]int64, 0, len(members))
	for _, member := range members {
		uid, err := strconv.ParseInt(member, 10, 64)
		if err != nil {
			continue
		}
		userIDs = append(userIDs, uid)
	}

	return userIDs, nil
}

func (g *GroupCache) DeleteGroupMembers(ctx context.Context, groupID int64) error {
	key := fmt.Sprintf("%s%d", constants.GroupMembersCacheKeyPrefix, groupID)
	return g.rdb.Del(ctx, key).Err()
}
