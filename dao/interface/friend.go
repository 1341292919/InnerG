package _interface

import (
	"InnerG/dao/db/model"
	"context"
)

type FriendDB interface {
	GetFriend(ctx context.Context, userID, friendID int64) (*model.Friend, bool, error)
	GetFriendRequest(ctx context.Context, fromUser, toUser int64) (*model.FriendRequest, bool, error)
	CreateFriendRequest(ctx context.Context, request *model.FriendRequest) error
	AcceptFriendRequest(ctx context.Context, requesterID, addresseeID int64, forwardID, reverseID int64, now int64) error
	RejectFriendRequest(ctx context.Context, requesterID, addresseeID int64) error
	DeleteFriend(ctx context.Context, userID, friendID int64) error
	ListFriends(ctx context.Context, userID int64) ([]*model.Friend, error)
	ListInboundRequests(ctx context.Context, userID int64) ([]*model.FriendRequest, error)
	IsFriend(ctx context.Context, userID, friendID int64) (bool, error)
}
