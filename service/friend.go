package service

import (
	"InnerG/dao"
	"InnerG/dao/db/model"
	_interface "InnerG/dao/interface"
	"InnerG/pkg/constants"
	"context"
	"fmt"
	"sync"
	"time"
)

var FriendSrvIns *FriendSrv
var FriendSrvOnce sync.Once

type FriendSrv struct {
	db _interface.FriendDB
}

func GetFriendSrv() *FriendSrv {
	FriendSrvOnce.Do(func() {
		FriendSrvIns = &FriendSrv{db: dao.NewFriendDao(context.Background()).Db}
	})
	return FriendSrvIns
}

func (s *FriendSrv) CreateFriendRequest(ctx context.Context, userID, friendID int64) error {
	if userID == friendID {
		return fmt.Errorf("不能添加自己为好友")
	}

	existing, exist, err := s.db.GetFriend(ctx, userID, friendID)
	if err != nil {
		return err
	}
	if exist {
		switch existing.Status {
		case constants.FriendPendingStatus:
			return fmt.Errorf("好友请求已存在")
		case constants.FriendAcceptedStatus:
			return fmt.Errorf("已经是好友")
		}
	}

	return s.db.CreateFriendRequest(ctx, &model.Friend{
		UserID:    userID,
		FriendID:  friendID,
		Status:    constants.FriendPendingStatus,
		CreatedAt: time.Now().Unix(),
	})
}

func (s *FriendSrv) AcceptFriendRequest(ctx context.Context, currentUserID, requesterID int64) error {
	request, exist, err := s.db.GetFriend(ctx, requesterID, currentUserID)
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("好友请求不存在")
	}
	if request.Status != constants.FriendPendingStatus {
		return fmt.Errorf("好友请求状态不是待处理")
	}

	return s.db.AcceptFriendRequest(ctx, requesterID, currentUserID, time.Now().Unix())
}

func (s *FriendSrv) RejectFriendRequest(ctx context.Context, currentUserID, requesterID int64) error {
	request, exist, err := s.db.GetFriend(ctx, requesterID, currentUserID)
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("好友请求不存在")
	}
	if request.Status != constants.FriendPendingStatus {
		return fmt.Errorf("好友请求状态不是待处理")
	}

	return s.db.RejectFriendRequest(ctx, requesterID, currentUserID)
}

func (s *FriendSrv) DeleteFriend(ctx context.Context, userID, friendID int64) error {
	ok, err := s.db.IsFriend(ctx, userID, friendID)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("双方不是好友")
	}

	return s.db.DeleteFriend(ctx, userID, friendID)
}

func (s *FriendSrv) ListFriends(ctx context.Context, userID int64) ([]*model.Friend, error) {
	return s.db.ListFriends(ctx, userID)
}

func (s *FriendSrv) ListInboundRequests(ctx context.Context, userID int64) ([]*model.Friend, error) {
	return s.db.ListInboundRequests(ctx, userID)
}

func (s *FriendSrv) IsFriend(ctx context.Context, userID, friendID int64) (bool, error) {
	return s.db.IsFriend(ctx, userID, friendID)
}
