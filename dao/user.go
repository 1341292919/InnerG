package dao

import (
	"InnerG/dao/cache"
	"InnerG/dao/db"
	_interface "InnerG/dao/interface"
	"context"
)

type UserDao struct {
	Db    _interface.UserDB
	Cache _interface.UserCache
}

func NewUserDao(ctx context.Context) *UserDao {
	return &UserDao{
		Db:    db.NewDBClient(),
		Cache: cache.NewRedisClient(),
	}
}
