package dao

import (
	"InnerG/dao/cache"
	"InnerG/dao/db"
	_interface "InnerG/dao/interface"
)

type GroupDao struct {
	Db    _interface.GroupDB
	Cache _interface.GroupCache
}

func NewGroupDao() *GroupDao {
	return &GroupDao{
		Db:    db.NewGroupDBClient(),
		Cache: cache.NewGroupCacheClient(),
	}
}
