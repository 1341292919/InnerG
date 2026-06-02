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

func BuildFriendRequest(request *model.FriendRequest) *types.FriendRequestResp {
	if request == nil {
		return nil
	}

	return &types.FriendRequestResp{
		ID:        request.ID,
		FromUser:  request.FromUser,
		ToUser:    request.ToUser,
		Status:    request.Status,
		CreatedAt: request.CreatedAt,
	}
}

func BuildFriendRequestList(requests []*model.FriendRequest) []*types.FriendRequestResp {
	res := make([]*types.FriendRequestResp, 0, len(requests))
	for _, request := range requests {
		res = append(res, BuildFriendRequest(request))
	}
	return res
}
