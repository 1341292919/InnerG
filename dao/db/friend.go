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
	return &friendDB{
		client: db,
	}
}

func (db *friendDB) GetFriend(ctx context.Context, userID, friendID int64) (*model.Friend, bool, error) {
	var friend *model.Friend
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
	return friend, true, nil
}

func (db *friendDB) CreateFriendRequest(ctx context.Context, request *model.Friend) error {
	err := db.client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing model.Friend
		err := tx.Table(constants.FriendTableName).
			Where("user_id = ? AND friend_id = ?", request.UserID, request.FriendID).
			First(&existing).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return tx.Table(constants.FriendTableName).Create(request).Error
			}
			return err
		}

		switch existing.Status {
		case constants.FriendDeletedStatus, constants.FriendRejectedStatus:
			result := tx.Table(constants.FriendTableName).
				Where("user_id = ? AND friend_id = ? AND status IN ?", request.UserID, request.FriendID, []int8{
					constants.FriendDeletedStatus,
					constants.FriendRejectedStatus,
				}).
				Updates(map[string]interface{}{
					"status":     constants.FriendPendingStatus,
					"created_at": request.CreatedAt,
				})
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return newFriendBusinessError("好友请求不存在")
			}
			return nil
		case constants.FriendPendingStatus, constants.FriendAcceptedStatus:
			return newFriendBusinessError("好友请求已存在")
		default:
			return newFriendBusinessError("好友请求状态异常")
		}
	})
	if err != nil {
		return wrapFriendDBError("CreateFriendRequest", err)
	}
	return nil
}

func (db *friendDB) AcceptFriendRequest(ctx context.Context, requesterID, addresseeID int64, now int64) error {
	err := db.client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Table(constants.FriendTableName).
			Where("user_id = ? AND friend_id = ? AND status = ?", requesterID, addresseeID, constants.FriendPendingStatus).
			Update("status", constants.FriendAcceptedStatus)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return newFriendBusinessError("好友请求不存在")
		}

		friend := &model.Friend{
			UserID:    addresseeID,
			FriendID:  requesterID,
			Status:    constants.FriendAcceptedStatus,
			CreatedAt: now,
		}
		return tx.Table(constants.FriendTableName).
			Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "user_id"}, {Name: "friend_id"}},
				DoUpdates: clause.Assignments(map[string]interface{}{
					"status": constants.FriendAcceptedStatus,
				}),
			}).
			Create(friend).Error
	})
	if err != nil {
		return wrapFriendDBError("AcceptFriendRequest", err)
	}
	return nil
}

func (db *friendDB) RejectFriendRequest(ctx context.Context, requesterID, addresseeID int64) error {
	result := db.client.WithContext(ctx).
		Table(constants.FriendTableName).
		Where("user_id = ? AND friend_id = ? AND status = ?", requesterID, addresseeID, constants.FriendPendingStatus).
		Update("status", constants.FriendRejectedStatus)
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
			Where("user_id = ? AND friend_id = ? AND status = ?", userID, friendID, constants.FriendAcceptedStatus).
			Update("status", constants.FriendDeletedStatus)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return newFriendBusinessError("双方不是好友")
		}

		result = tx.Table(constants.FriendTableName).
			Where("user_id = ? AND friend_id = ? AND status = ?", friendID, userID, constants.FriendAcceptedStatus).
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
		Where("user_id = ? AND status = ?", userID, constants.FriendAcceptedStatus).
		Order("updated_at DESC").
		Find(&friends).Error
	if err != nil {
		return nil, errno.NewErr(errno.MySQLDBErrorCode, "ListFriends: "+err.Error())
	}
	return friends, nil
}

func (db *friendDB) ListInboundRequests(ctx context.Context, userID int64) ([]*model.Friend, error) {
	var requests []*model.Friend
	err := db.client.WithContext(ctx).
		Table(constants.FriendTableName).
		Where("friend_id = ? AND status = ?", userID, constants.FriendPendingStatus).
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
		Where("user_id = ? AND friend_id = ? AND status = ?", userID, friendID, constants.FriendAcceptedStatus).
		Count(&count).Error
	if err != nil {
		return false, errno.NewErr(errno.MySQLDBErrorCode, "IsFriend: "+err.Error())
	}
	return count > 0, nil
}
