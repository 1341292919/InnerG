package dao

import (
	"InnerG/dao/db"
	_interface "InnerG/dao/interface"
	"context"
)

type MusicDao struct {
	Db _interface.MusicDB
}

func NewMusicDao(ctx context.Context) *MusicDao {
	return &MusicDao{
		Db: db.NewMusicDBClient(),
	}
}
