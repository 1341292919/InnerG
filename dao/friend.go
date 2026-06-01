package dao

import (
	"InnerG/dao/db"
	_interface "InnerG/dao/interface"
	"context"
)

type FriendDao struct {
	Db _interface.FriendDB
}

func NewFriendDao(ctx context.Context) *FriendDao {
	return &FriendDao{Db: db.NewFriendDBClient()}
}
