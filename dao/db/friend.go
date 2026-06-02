package db

import (
	"InnerG/dao/db/model"
	_interface "InnerG/dao/interface"
	"InnerG/pkg/constants"
	"InnerG/pkg/errno"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type friendDB struct {
	client *gorm.DB
}

type friendBusinessError struct {
	error
}

func newFriendBusinessError(format string, args ...interface{}) error {
	return friendBusinessError{error: fmt.Errorf(format, args...)}
}

func wrapFriendDBError(op string, err error) error {
	var businessErr friendBusinessError
	if errors.As(err, &businessErr) {
		return err
	}
	return errno.NewErr(errno.MySQLDBErrorCode, op+": "+err.Error())
}

func NewFriendDB(db *gorm.DB) _interface.FriendDB {
	return &friendDB{client: db}
}

func (db *friendDB) GetFriend(ctx context.Context, userID, friendID int64) (*model.Friend, bool, error) {
	var friend model.Friend
	err := db.client.WithContext(ctx).
		Table(constants.FriendTableName).
		Where("user_id = ? AND friend_id = ?", userID, friendID).
		First(&friend).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, errno.NewErr(errno.MySQLDBErrorCode, "GetFriend: "+err.Error())
	}
	return &friend, true, nil
}

func (db *friendDB) GetFriendRequest(ctx context.Context, fromUser, toUser int64) (*model.FriendRequest, bool, error) {
	var request model.FriendRequest
	err := db.client.WithContext(ctx).
		Table(constants.FriendRequestTableName).
		Where("from_user = ? AND to_user = ?", fromUser, toUser).
		First(&request).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, errno.NewErr(errno.MySQLDBErrorCode, "GetFriendRequest: "+err.Error())
	}
	return &request, true, nil
}

func (db *friendDB) CreateFriendRequest(ctx context.Context, request *model.FriendRequest) error {
	err := db.client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing model.FriendRequest
		err := tx.Table(constants.FriendRequestTableName).
			Where("from_user = ? AND to_user = ?", request.FromUser, request.ToUser).
			First(&existing).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return tx.Table(constants.FriendRequestTableName).Create(request).Error
			}
			return err
		}

		switch existing.Status {
		case constants.FriendRequestPendingStatus:
			return newFriendBusinessError("好友请求已存在")
		case constants.FriendRequestAcceptedStatus:
			return newFriendBusinessError("已经是好友")
		case constants.FriendRequestRejectedStatus, constants.FriendRequestCancelledStatus:
			result := tx.Table(constants.FriendRequestTableName).
				Where("from_user = ? AND to_user = ? AND status IN ?", request.FromUser, request.ToUser, []int8{
					constants.FriendRequestRejectedStatus,
					constants.FriendRequestCancelledStatus,
				}).
				Updates(map[string]interface{}{
					"status":     constants.FriendRequestPendingStatus,
					"created_at": request.CreatedAt,
				})
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return newFriendBusinessError("好友请求不存在")
			}
			return nil
		default:
			return newFriendBusinessError("好友请求状态异常")
		}
	})
	if err != nil {
		return wrapFriendDBError("CreateFriendRequest", err)
	}
	return nil
}

func (db *friendDB) AcceptFriendRequest(ctx context.Context, requesterID, addresseeID int64, forwardID, reverseID int64, now int64) error {
	err := db.client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Table(constants.FriendRequestTableName).
			Where("from_user = ? AND to_user = ? AND status = ?", requesterID, addresseeID, constants.FriendRequestPendingStatus).
			Update("status", constants.FriendRequestAcceptedStatus)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return newFriendBusinessError("好友请求不存在")
		}

		friends := []*model.Friend{
			{ID: forwardID, UserID: requesterID, FriendID: addresseeID, Status: constants.FriendActiveStatus, CreatedAt: now},
			{ID: reverseID, UserID: addresseeID, FriendID: requesterID, Status: constants.FriendActiveStatus, CreatedAt: now},
		}
		for _, friend := range friends {
			if err := tx.Table(constants.FriendTableName).
				Clauses(clause.OnConflict{
					Columns: []clause.Column{{Name: "user_id"}, {Name: "friend_id"}},
					DoUpdates: clause.Assignments(map[string]interface{}{
						"status": constants.FriendActiveStatus,
					}),
				}).
				Create(friend).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return wrapFriendDBError("AcceptFriendRequest", err)
	}
	return nil
}

func (db *friendDB) RejectFriendRequest(ctx context.Context, requesterID, addresseeID int64) error {
	result := db.client.WithContext(ctx).
		Table(constants.FriendRequestTableName).
		Where("from_user = ? AND to_user = ? AND status = ?", requesterID, addresseeID, constants.FriendRequestPendingStatus).
		Update("status", constants.FriendRequestRejectedStatus)
	if result.Error != nil {
		return errno.NewErr(errno.MySQLDBErrorCode, "RejectFriendRequest: "+result.Error.Error())
	}
	if result.RowsAffected != 1 {
		return fmt.Errorf("好友请求不存在")
	}
	return nil
}

func (db *friendDB) DeleteFriend(ctx context.Context, userID, friendID int64) error {
	err := db.client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Table(constants.FriendTableName).
			Where("user_id = ? AND friend_id = ? AND status = ?", userID, friendID, constants.FriendActiveStatus).
			Update("status", constants.FriendDeletedStatus)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return newFriendBusinessError("双方不是好友")
		}

		result = tx.Table(constants.FriendTableName).
			Where("user_id = ? AND friend_id = ? AND status = ?", friendID, userID, constants.FriendActiveStatus).
			Update("status", constants.FriendDeletedStatus)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return newFriendBusinessError("双方不是好友")
		}
		return nil
	})
	if err != nil {
		return wrapFriendDBError("DeleteFriend", err)
	}
	return nil
}

func (db *friendDB) ListFriends(ctx context.Context, userID int64) ([]*model.Friend, error) {
	var friends []*model.Friend
	err := db.client.WithContext(ctx).
		Table(constants.FriendTableName).
		Where("user_id = ? AND status = ?", userID, constants.FriendActiveStatus).
		Order("updated_at DESC").
		Find(&friends).Error
	if err != nil {
		return nil, errno.NewErr(errno.MySQLDBErrorCode, "ListFriends: "+err.Error())
	}
	return friends, nil
}

func (db *friendDB) ListInboundRequests(ctx context.Context, userID int64) ([]*model.FriendRequest, error) {
	var requests []*model.FriendRequest
	err := db.client.WithContext(ctx).
		Table(constants.FriendRequestTableName).
		Where("to_user = ? AND status = ?", userID, constants.FriendRequestPendingStatus).
		Order("created_at DESC").
		Find(&requests).Error
	if err != nil {
		return nil, errno.NewErr(errno.MySQLDBErrorCode, "ListInboundRequests: "+err.Error())
	}
	return requests, nil
}

func (db *friendDB) IsFriend(ctx context.Context, userID, friendID int64) (bool, error) {
	var count int64
	err := db.client.WithContext(ctx).
		Table(constants.FriendTableName).
		Where("user_id = ? AND friend_id = ? AND status = ?", userID, friendID, constants.FriendActiveStatus).
		Count(&count).Error
	if err != nil {
		return false, errno.NewErr(errno.MySQLDBErrorCode, "IsFriend: "+err.Error())
	}
	return count > 0, nil
}
