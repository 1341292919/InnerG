package pack

import (
	"InnerG/dao/db/model"
	"InnerG/types"
)

func BuildFriend(friend *model.Friend) *types.FriendResp {
	if friend == nil {
		return nil
	}

	return &types.FriendResp{
		UserID:    friend.UserID,
		FriendID:  friend.FriendID,
		Status:    friend.Status,
		CreatedAt: friend.CreatedAt,
	}
}

func BuildFriendList(friends []*model.Friend) []*types.FriendResp {
	res := make([]*types.FriendResp, 0, len(friends))
	for _, friend := range friends {
		res = append(res, BuildFriend(friend))
	}
	return res
}
