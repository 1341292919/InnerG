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

func (db *userDB) IsUserExistById(ctx context.Context, id int64) (*model.User, bool, error) {
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

func (db *userDB) GetUserBasicByIds(ctx context.Context, ids []int64) (map[int64]*model.User, error) {
	users := make(map[int64]*model.User, len(ids))
	if len(ids) == 0 {
		return users, nil
	}

	var list []*model.User
	err := db.client.WithContext(ctx).
		Table(constants.UserTableName).
		Select("id", "username", "avatar").
		Where("id IN ?", ids).
		Find(&list).Error
	if err != nil {
		return nil, errno.NewErr(errno.MySQLDBErrorCode, "GetUserBasicByIds: "+err.Error())
	}
	for _, user := range list {
		users[int64(user.ID)] = user
	}
	return users, nil
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

func (db *userDB) UpdateUserAccount(ctx context.Context, account string, id int64) error {
	if err := db.client.WithContext(ctx).
		Table(constants.UserTableName).
		Where("id = ?", id).
		Update("account", account).
		Error; err != nil {
		return errno.NewErr(errno.MySQLDBErrorCode, "UpdateUserAccount: "+err.Error())
	}
	return nil
}

func (db *userDB) UpdateUserName(ctx context.Context, userName string, id int64) error {
	if err := db.client.WithContext(ctx).
		Table(constants.UserTableName).
		Where("id = ?", id).
		Update("username", userName).
		Error; err != nil {
		return errno.NewErr(errno.MySQLDBErrorCode, "UpdateUserName: "+err.Error())
	}
	return nil
}

func (db *userDB) UpdateUserGender(ctx context.Context, gender string, id int64) error {
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
