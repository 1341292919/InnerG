package db

import (
	"InnerG/dao/db/model"
	_interface "InnerG/dao/interface"
	"InnerG/pkg/constants"
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
	return db.client.WithContext(ctx).Table(constants.UserTableName).Create(user).Error
}

func (db *userDB) IsUserExistById(ctx context.Context, id string) (*model.User, bool, error) {
	var user *model.User
	err := db.client.WithContext(ctx).Table(constants.UserTableName).Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
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
		return nil, false, err
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
		return nil, false, err
	}
	return user, true, nil
}

func (db *userDB) UpdateUserAccount(ctx context.Context, account string, id string) error {
	return db.client.WithContext(ctx).Table(constants.UserTableName).Where("id = ?", id).Update("account", account).Error
}

func (db *userDB) UpdateUserAvatar(ctx context.Context, id, avatarUrl string) error {
	return db.client.WithContext(ctx).Table(constants.UserTableName).Where("id = ?", id).Update("avatar", avatarUrl).Error
}
