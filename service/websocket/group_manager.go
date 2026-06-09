package websocket

import (
	"InnerG/dao"
	"InnerG/pkg/constants"
	"InnerG/pkg/logger"
	"context"
	"sync"
)

const defaultGroupShardNum = 32

type groupShard struct {
	mu            sync.RWMutex
	onlineMembers map[int64]map[int64]bool
}

type GroupManager struct {
	shards []*groupShard
}

func NewGroupManager() *GroupManager {
	return NewGroupManagerWithShards(defaultGroupShardNum)
}

func NewGroupManagerWithShards(shardNum int) *GroupManager {
	if shardNum <= 0 {
		shardNum = defaultGroupShardNum
	}
	shards := make([]*groupShard, shardNum)
	for i := range shards {
		shards[i] = &groupShard{
			onlineMembers: make(map[int64]map[int64]bool),
		}
	}
	return &GroupManager{shards: shards}
}

func (gm *GroupManager) getShard(groupID int64) *groupShard {
	return gm.shards[groupID%int64(len(gm.shards))]
}

func (gm *GroupManager) SubscribeUser(userID int64, groupIDs []int64) {
	for _, groupID := range groupIDs {
		s := gm.getShard(groupID)
		s.mu.Lock()
		if s.onlineMembers[groupID] == nil {
			s.onlineMembers[groupID] = make(map[int64]bool)
		}
		s.onlineMembers[groupID][userID] = true
		s.mu.Unlock()
	}

	logger.Log.Infof("user %d subscribed to %d groups", userID, len(groupIDs))
}

func (gm *GroupManager) UnsubscribeUser(userID int64) {
	for _, shard := range gm.shards {
		shard.mu.Lock()
		for groupID := range shard.onlineMembers {
			delete(shard.onlineMembers[groupID], userID)
			if len(shard.onlineMembers[groupID]) == 0 {
				delete(shard.onlineMembers, groupID)
			}
		}
		shard.mu.Unlock()
	}

	logger.Log.Infof("user %d unsubscribed from all groups", userID)
}

func (gm *GroupManager) SubscribeGroup(userID int64, groupID int64) {
	s := gm.getShard(groupID)
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.onlineMembers[groupID] == nil {
		s.onlineMembers[groupID] = make(map[int64]bool)
	}
	s.onlineMembers[groupID][userID] = true
}

func (gm *GroupManager) UnsubscribeGroup(userID int64, groupID int64) {
	s := gm.getShard(groupID)
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.onlineMembers[groupID] != nil {
		delete(s.onlineMembers[groupID], userID)
		if len(s.onlineMembers[groupID]) == 0 {
			delete(s.onlineMembers, groupID)
		}
	}
}

func (gm *GroupManager) GetOnlineMembers(groupID int64) []int64 {
	s := gm.getShard(groupID)
	s.mu.RLock()
	defer s.mu.RUnlock()

	members := make([]int64, 0)
	if s.onlineMembers[groupID] != nil {
		for userID := range s.onlineMembers[groupID] {
			members = append(members, userID)
		}
	}
	return members
}

func (gm *GroupManager) GetOnlineMembersWithCache(ctx context.Context, groupID int64) []int64 {
	onlineMembers := gm.GetOnlineMembers(groupID)
	if len(onlineMembers) > 0 {
		return onlineMembers
	}

	groupDao := dao.NewGroupDao()
	cachedMembers, err := groupDao.Cache.GetGroupMembers(ctx, groupID)
	if err == nil && len(cachedMembers) > 0 {
		s := gm.getShard(groupID)
		s.mu.RLock()
		onlineMembers = make([]int64, 0)
		if s.onlineMembers[groupID] != nil {
			for _, memberID := range cachedMembers {
				if s.onlineMembers[groupID][memberID] {
					onlineMembers = append(onlineMembers, memberID)
				}
			}
		}
		s.mu.RUnlock()
		return onlineMembers
	}

	allMembers, err := groupDao.Db.GetGroupMemberIDs(ctx, groupID)
	if err != nil {
		logger.Log.Errorf("GetOnlineMembersWithCache:GetGroupMemberIDs: %v", err)
		return []int64{}
	}

	_ = groupDao.Cache.SetGroupMembers(ctx, groupID, allMembers, constants.GroupMembersCacheTTL)

	s := gm.getShard(groupID)
	s.mu.RLock()
	onlineMembers = make([]int64, 0)
	if s.onlineMembers[groupID] != nil {
		for _, memberID := range allMembers {
			if s.onlineMembers[groupID][memberID] {
				onlineMembers = append(onlineMembers, memberID)
			}
		}
	}
	s.mu.RUnlock()

	return onlineMembers
}

func (gm *GroupManager) IsOnline(groupID int64, userID int64) bool {
	s := gm.getShard(groupID)
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.onlineMembers[groupID] == nil {
		return false
	}
	return s.onlineMembers[groupID][userID]
}

func (gm *GroupManager) GetGroupCount() int {
	count := 0
	for _, shard := range gm.shards {
		shard.mu.RLock()
		count += len(shard.onlineMembers)
		shard.mu.RUnlock()
	}
	return count
}
