package db

import (
	"InnerG/dao/db/model"
	_interface "InnerG/dao/interface"
	"InnerG/pkg/constants"
	"InnerG/pkg/errno"
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type friendDB struct {
	client *gorm.DB
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
	var existing model.Friend
	err := db.client.WithContext(ctx).
		Table(constants.FriendTableName).
		Where("user_id = ? AND friend_id = ?", request.UserID, request.FriendID).
		First(&existing).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := db.client.WithContext(ctx).Table(constants.FriendTableName).Create(request).Error; err != nil {
				return errno.NewErr(errno.MySQLDBErrorCode, "CreateFriendRequest: "+err.Error())
			}
			return nil
		}
		return errno.NewErr(errno.MySQLDBErrorCode, "CreateFriendRequest: "+err.Error())
	}

	if existing.Status != constants.FriendRejectedStatus && existing.Status != constants.FriendDeletedStatus {
		return errno.NewErr(errno.MySQLDBErrorCode, "CreateFriendRequest: active friend row exists")
	}

	err = db.client.WithContext(ctx).
		Table(constants.FriendTableName).
		Where("user_id = ? AND friend_id = ?", request.UserID, request.FriendID).
		Updates(map[string]interface{}{
			"status":     constants.FriendPendingStatus,
			"created_at": request.CreatedAt,
		}).Error
	if err != nil {
		return errno.NewErr(errno.MySQLDBErrorCode, "CreateFriendRequest: "+err.Error())
	}
	return nil
}

func (db *friendDB) AcceptFriendRequest(ctx context.Context, requesterID, addresseeID int64, now int64) error {
	err := db.client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(constants.FriendTableName).
			Where("user_id = ? AND friend_id = ?", requesterID, addresseeID).
			Update("status", constants.FriendAcceptedStatus).Error; err != nil {
			return err
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
		return errno.NewErr(errno.MySQLDBErrorCode, "AcceptFriendRequest: "+err.Error())
	}
	return nil
}

func (db *friendDB) RejectFriendRequest(ctx context.Context, requesterID, addresseeID int64) error {
	err := db.client.WithContext(ctx).
		Table(constants.FriendTableName).
		Where("user_id = ? AND friend_id = ?", requesterID, addresseeID).
		Update("status", constants.FriendRejectedStatus).Error
	if err != nil {
		return errno.NewErr(errno.MySQLDBErrorCode, "RejectFriendRequest: "+err.Error())
	}
	return nil
}

func (db *friendDB) DeleteFriend(ctx context.Context, userID, friendID int64) error {
	err := db.client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(constants.FriendTableName).
			Where("user_id = ? AND friend_id = ?", userID, friendID).
			Update("status", constants.FriendDeletedStatus).Error; err != nil {
			return err
		}
		return tx.Table(constants.FriendTableName).
			Where("user_id = ? AND friend_id = ?", friendID, userID).
			Update("status", constants.FriendDeletedStatus).Error
	})
	if err != nil {
		return errno.NewErr(errno.MySQLDBErrorCode, "DeleteFriend: "+err.Error())
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
