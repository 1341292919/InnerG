package dao

import (
	"InnerG/dao/cache"
	"InnerG/dao/db"
	_interface "InnerG/dao/interface"
	"context"
)

type MusicDao struct {
	Db _interface.MusicDB
	Ca _interface.MusicCache
}

func NewMusicDao(ctx context.Context) *MusicDao {
	return &MusicDao{
		Db: db.NewMusicDBClient(),
		Ca: cache.NewMusicClient(),
	}
}
