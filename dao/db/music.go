package db

import (
	"InnerG/dao/db/model"
	_interface "InnerG/dao/interface"
	"InnerG/pkg/constants"
	"context"
	"errors"
	"gorm.io/gorm"
)

type musicDB struct {
	client *gorm.DB
}

func NewMusicDB(db *gorm.DB) _interface.MusicDB {
	return &musicDB{
		client: db,
	}
}

func (db *musicDB) GetPlaylistList(ctx context.Context, pageNum, pageSize int) ([]*model.Playlist, int, error) {
	if pageNum < 1 {
		pageNum = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	var total int64
	err := db.client.WithContext(ctx).Table(constants.PlaylistTableName).Where("status <> ?", 2).Count(&total).Error
	if err != nil {
		return nil, -1, err
	}
	list := make([]*model.Playlist, 0)
	err = db.client.WithContext(ctx).
		Table(constants.PlaylistTableName).
		Where("status <> ?", 2).
		Order("updated_at DESC").
		Limit(pageSize).
		Offset((pageNum - 1) * pageSize).
		Find(&list).Error
	if err != nil {
		return nil, -1, err
	}
	return list, int(total), nil
}

func (db *musicDB) GetPlaylistById(ctx context.Context, id string) (*model.Playlist, bool, error) {
	var playlist *model.Playlist
	err := db.client.WithContext(ctx).Table(constants.PlaylistTableName).Where("id = ?", id).First(&playlist).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return playlist, true, nil
}

func (db *musicDB) GetPlaylistSongListByPlaylistId(ctx context.Context, playlistId string) ([]*model.PlaylistSong, error) {
	res := make([]*model.PlaylistSong, 0)
	err := db.client.WithContext(ctx).
		Table(constants.PlaylistSongTableName+" ps").
		Select("s.id, s.name, sg.name as singer_name, s.created_at").
		Joins("JOIN "+constants.SongTableName+" s ON s.id = ps.song_id").
		Joins("LEFT JOIN "+constants.SingerTableName+" sg ON sg.id = s.singer_id").
		Where("ps.playlist_id = ?", playlistId).
		Where("ps.deleted_at IS NULL").
		Where("s.deleted_at IS NULL").
		Order("ps.sort_order ASC, ps.id ASC").
		Scan(&res).Error
	return res, err
}

func (db *musicDB) GetSongList(ctx context.Context, pageNum, pageSize int) ([]*model.Song, int, error) {
	if pageNum < 1 {
		pageNum = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	var total int64
	err := db.client.WithContext(ctx).Table(constants.SongTableName).Where("status <> ?", 2).Count(&total).Error
	if err != nil {
		return nil, -1, err
	}
	list := make([]*model.Song, 0)
	err = db.client.WithContext(ctx).
		Table(constants.SongTableName).
		Where("status <> ?", 2).
		Order("updated_at DESC").
		Limit(pageSize).
		Offset((pageNum - 1) * pageSize).
		Find(&list).Error
	if err != nil {
		return nil, -1, err
	}
	return list, int(total), nil
}

func (db *musicDB) GetSongById(ctx context.Context, id string) (*model.Song, bool, error) {
	var song *model.Song
	err := db.client.WithContext(ctx).Table(constants.SongTableName).Where("id = ?", id).First(&song).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return song, true, nil
}
