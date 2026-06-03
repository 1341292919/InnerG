package service

import (
	"InnerG/config"
	"InnerG/dao"
	"InnerG/dao/db/model"
	_interface "InnerG/dao/interface"
	"InnerG/dao/rabbitmq"
	"InnerG/pkg/constants"
	"InnerG/pkg/logger"
	"InnerG/pkg/utils"
	"InnerG/types"
	"context"
	"fmt"
	"sync"
	"time"
)

var FriendSrvIns *FriendSrv
var FriendSrvOnce sync.Once

type FriendSrv struct {
	db             _interface.FriendDB
	sf             *utils.Snowflake
	eventPublisher friendEventPublisher
}

type friendEventPublisher interface {
	SendMessage(topic string, data interface{}) error
}

func GetFriendSrv() *FriendSrv {
	FriendSrvOnce.Do(func() {
		sf, _ := utils.NewSnowflake(config.Snowflake.DatancenterID, config.Snowflake.WorkerID)
		producer, err := rabbitmq.NewProducer(constants.WebsocketService)
		if err != nil {
			logFriendEventError("init friend event publisher: ", err)
		}
		FriendSrvIns = &FriendSrv{
			db:             dao.NewFriendDao(context.Background()).Db,
			sf:             sf,
			eventPublisher: producer,
		}
	})
	return FriendSrvIns
}

func (s *FriendSrv) publishFriendRequestEvent(event types.FriendRequestEvent) {
	if s.eventPublisher == nil {
		return
	}
	if err := s.eventPublisher.SendMessage(constants.FriendRequestMessageTopic, event); err != nil {
		logFriendEventError("publish friend request event: ", err)
	}
}

func logFriendEventError(v ...interface{}) {
	if logger.Log == nil {
		return
	}
	logger.Log.Error(v...)
}

func (s *FriendSrv) CreateFriendRequest(ctx context.Context, userID, friendID int64, message string) error {
	if userID == friendID {
		return fmt.Errorf("不能添加自己为好友")
	}
	if len([]rune(message)) > 100 {
		return fmt.Errorf("好友申请打招呼内容不能超过100个字符")
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
	eventRequestID := id
	if exist && request != nil {
		eventRequestID = request.ID
	}
	now := time.Now().Unix()

	err = s.db.CreateFriendRequest(ctx, &model.FriendRequest{
		ID:        id,
		FromUser:  userID,
		ToUser:    friendID,
		Status:    constants.FriendRequestPendingStatus,
		Message:   message,
		CreatedAt: now,
	})
	if err != nil {
		return err
	}
	s.publishFriendRequestEvent(types.FriendRequestEvent{
		Type:      constants.FriendRequestEventType,
		RequestID: eventRequestID,
		FromUser:  userID,
		ToUser:    friendID,
		Message:   message,
		CreatedAt: now,
	})
	return nil
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

	now := time.Now().Unix()
	if err := s.db.AcceptFriendRequest(ctx, requesterID, currentUserID, forwardID, reverseID, now); err != nil {
		return err
	}
	s.publishFriendRequestEvent(types.FriendRequestEvent{
		Type:      constants.FriendRequestAcceptedEventType,
		RequestID: request.ID,
		FromUser:  currentUserID,
		ToUser:    requesterID,
		Message:   constants.FriendRequestAcceptedMessage,
		CreatedAt: now,
	})
	return nil
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

func (s *FriendSrv) ListFriends(ctx context.Context, userID int64, pageSize, pageNum int) ([]*model.Friend, int64, error) {
	return s.db.ListFriends(ctx, userID, pageSize, pageNum)
}

func (s *FriendSrv) ListInboundRequests(ctx context.Context, userID int64, pageSize, pageNum int) ([]*model.FriendRequest, int64, error) {
	return s.db.ListInboundRequests(ctx, userID, pageSize, pageNum)
}

func (s *FriendSrv) IsFriend(ctx context.Context, userID, friendID int64) (bool, error) {
	return s.db.IsFriend(ctx, userID, friendID)
}
