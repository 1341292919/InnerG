package _interface

import (
	"InnerG/dao/db/model"
	"context"
)

type MusicDB interface {
	GetPlaylistList(ctx context.Context, pageNum, pageSize int) ([]*model.Playlist, int, error)
	GetPlaylistById(ctx context.Context, id string) (*model.Playlist, bool, error)
	GetPlaylistSongListByPlaylistId(ctx context.Context, playlistId string) ([]*model.PlaylistSong, error)
	GetSongList(ctx context.Context, pageNum, pageSize int) ([]*model.Song, int, error)
	GetSongById(ctx context.Context, id string) (*model.Song, bool, error)
}
