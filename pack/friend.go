package pack

import (
	"InnerG/dao/db/model"
	"InnerG/types"
)

func BuildFriend(friend *model.Friend, users map[int64]*model.User) *types.FriendResp {
	if friend == nil {
		return nil
	}
	resp := &types.FriendResp{
		UserID:    friend.UserID,
		FriendID:  friend.FriendID,
		Status:    friend.Status,
		CreatedAt: friend.CreatedAt,
	}
	if user, ok := users[friend.FriendID]; ok && user != nil {
		resp.Avatar = user.Avatar
		resp.Username = user.Username
	}
	return resp
}

func BuildFriendList(friends []*model.Friend, users map[int64]*model.User) []*types.FriendResp {
	res := make([]*types.FriendResp, 0, len(friends))
	for _, friend := range friends {
		res = append(res, BuildFriend(friend, users))
	}
	return res
}

func BuildFriendRequest(request *model.FriendRequest, users map[int64]*model.User) *types.FriendRequestResp {
	if request == nil {
		return nil
	}
	resp := &types.FriendRequestResp{
		ID:        request.ID,
		FromUser:  request.FromUser,
		ToUser:    request.ToUser,
		Status:    request.Status,
		CreatedAt: request.CreatedAt,
	}
	if user, ok := users[request.FromUser]; ok && user != nil {
		resp.FromUserAvatar = user.Avatar
		resp.FromUserName = user.Username
	}
	return resp
}

func BuildFriendRequestList(requests []*model.FriendRequest, users map[int64]*model.User) []*types.FriendRequestResp {
	res := make([]*types.FriendRequestResp, 0, len(requests))
	for _, request := range requests {
		res = append(res, BuildFriendRequest(request, users))
	}
	return res
}
