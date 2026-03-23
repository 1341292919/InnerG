package db

import (
	"InnerG/dao/db/model"
	_interface "InnerG/dao/interface"
	"InnerG/pkg/constants"
	"InnerG/pkg/errno"
	"context"
	"errors"
	"gorm.io/gorm"
)

type userDB struct {
	client *gorm.DB
}

func NewUserDB(db *gorm.DB) _interface.UserDB {
	return &userDB{
		client: db,
	}
}
func (db *userDB) CreateNewUser(ctx context.Context, user *model.User) error {
	err := db.client.WithContext(ctx).Table(constants.UserTableName).Create(user).Error
	if err != nil {
		return errno.NewErr(errno.MySQLDBErrorCode, "CreateNewUser: "+err.Error())
	}
	return nil
}

func (db *userDB) IsUserExistById(ctx context.Context, id string) (*model.User, bool, error) {
	var user *model.User
	err := db.client.WithContext(ctx).Table(constants.UserTableName).Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, errno.NewErr(errno.MySQLDBErrorCode, "IsUserExistById: "+err.Error())
	}
	return user, true, nil
}

func (db *userDB) IsUserExistByEmail(ctx context.Context, email string) (*model.User, bool, error) {
	var user *model.User
	err := db.client.WithContext(ctx).Table(constants.UserTableName).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, errno.NewErr(errno.MySQLDBErrorCode, "IsUserExistByEmail: "+err.Error())
	}
	return user, true, nil
}

func (db *userDB) IsUserExistByAccount(ctx context.Context, account string) (*model.User, bool, error) {
	var user *model.User
	err := db.client.WithContext(ctx).Table(constants.UserTableName).Where("account = ?", account).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, errno.NewErr(errno.MySQLDBErrorCode, "IsUserExistByAccount: "+err.Error())
	}
	return user, true, nil
}

func (db *userDB) UpdateUserAccount(ctx context.Context, account string, id string) error {
	if err := db.client.WithContext(ctx).
		Table(constants.UserTableName).
		Where("id = ?", id).
		Update("account", account).
		Error; err != nil {
		return errno.NewErr(errno.MySQLDBErrorCode, "UpdateUserAccount: "+err.Error())
	}
	return nil
}

func (db *userDB) UpdateUserName(ctx context.Context, userName string, id string) error {
	if err := db.client.WithContext(ctx).
		Table(constants.UserTableName).
		Where("id = ?", id).
		Update("username", userName).
		Error; err != nil {
		return errno.NewErr(errno.MySQLDBErrorCode, "UpdateUserName: "+err.Error())
	}
	return nil
}

func (db *userDB) UpdateUserGender(ctx context.Context, gender string, id string) error {
	if err := db.client.WithContext(ctx).
		Table(constants.UserTableName).
		Where("id = ?", id).
		Update("gender", gender).
		Error; err != nil {
		return errno.NewErr(errno.MySQLDBErrorCode, "UpdateUserGender: "+err.Error())
	}
	return nil
}

func (db *userDB) UpdateUserAvatar(ctx context.Context, id, avatarUrl string) error {
	if err := db.client.WithContext(ctx).
		Table(constants.UserTableName).
		Where("id = ?", id).
		Update("avatar", avatarUrl).
		Error; err != nil {
		return errno.NewErr(errno.MySQLDBErrorCode, "UpdateUserAvatar: "+err.Error())
	}
	return nil
}
