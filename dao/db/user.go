package db

import (
	"InnerG/dao/db/model"
	_interface "InnerG/dao/interface"
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
	return db.client.WithContext(ctx).Table("").Create(user).Error
}

func (db *userDB) IsUserExistByEmail(ctx context.Context, email string) (*model.User, bool, error) {
	var user *model.User
	err := db.client.WithContext(ctx).Table("").Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return user, true, nil
}
