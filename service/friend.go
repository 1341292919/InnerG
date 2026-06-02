package service

import (
	"InnerG/config"
	"InnerG/dao"
	"InnerG/dao/db/model"
	_interface "InnerG/dao/interface"
	"InnerG/pkg/constants"
	"InnerG/pkg/utils"
	"context"
	"fmt"
	"sync"
	"time"
)

var FriendSrvIns *FriendSrv
var FriendSrvOnce sync.Once

type FriendSrv struct {
	db _interface.FriendDB
	sf *utils.Snowflake
}

func GetFriendSrv() *FriendSrv {
	FriendSrvOnce.Do(func() {
		sf, _ := utils.NewSnowflake(config.Snowflake.DatancenterID, config.Snowflake.WorkerID)
		FriendSrvIns = &FriendSrv{
			db: dao.NewFriendDao(context.Background()).Db,
			sf: sf,
		}
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
	if exist && existing.Status == constants.FriendActiveStatus {
		return fmt.Errorf("已经是好友")
	}

	request, exist, err := s.db.GetFriendRequest(ctx, userID, friendID)
	if err != nil {
		return err
	}
	if exist && request.Status == constants.FriendRequestPendingStatus {
		return fmt.Errorf("好友请求已存在")
	}

	id, err := s.sf.NextVal()
	if err != nil {
		return err
	}

	return s.db.CreateFriendRequest(ctx, &model.FriendRequest{
		ID:        id,
		FromUser:  userID,
		ToUser:    friendID,
		Status:    constants.FriendRequestPendingStatus,
		CreatedAt: time.Now().Unix(),
	})
}

func (s *FriendSrv) AcceptFriendRequest(ctx context.Context, currentUserID, requesterID int64) error {
	request, exist, err := s.db.GetFriendRequest(ctx, requesterID, currentUserID)
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("好友请求不存在")
	}
	if request.Status != constants.FriendRequestPendingStatus {
		return fmt.Errorf("好友请求状态不是待处理")
	}

	forwardID, err := s.sf.NextVal()
	if err != nil {
		return err
	}
	reverseID, err := s.sf.NextVal()
	if err != nil {
		return err
	}

	return s.db.AcceptFriendRequest(ctx, requesterID, currentUserID, forwardID, reverseID, time.Now().Unix())
}

func (s *FriendSrv) HandleFriendRequest(ctx context.Context, currentUserID, requesterID int64, actionType string) error {
	switch actionType {
	case "accept":
		return s.AcceptFriendRequest(ctx, currentUserID, requesterID)
	case "reject":
		return s.RejectFriendRequest(ctx, currentUserID, requesterID)
	default:
		return fmt.Errorf("好友请求操作类型错误")
	}
}

func (s *FriendSrv) RejectFriendRequest(ctx context.Context, currentUserID, requesterID int64) error {
	request, exist, err := s.db.GetFriendRequest(ctx, requesterID, currentUserID)
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("好友请求不存在")
	}
	if request.Status != constants.FriendRequestPendingStatus {
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

func (s *FriendSrv) ListInboundRequests(ctx context.Context, userID int64) ([]*model.FriendRequest, error) {
	return s.db.ListInboundRequests(ctx, userID)
}

func (s *FriendSrv) IsFriend(ctx context.Context, userID, friendID int64) (bool, error) {
	return s.db.IsFriend(ctx, userID, friendID)
}
