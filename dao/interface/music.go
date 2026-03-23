package _interface

import (
	"InnerG/dao/db/model"
	"InnerG/types"
	"context"
)

type MusicDB interface {
	GetPlaylistList(ctx context.Context, pageNum, pageSize int) ([]*model.Playlist, int, error)
	GetPlaylistById(ctx context.Context, id string) (*model.Playlist, bool, error)
	GetPlaylistSongListByPlaylistId(ctx context.Context, playlistId string) ([]*model.PlaylistSong, error)
	GetSongList(ctx context.Context, pageNum, pageSize int) ([]*model.Song, int, error)
	GetSongById(ctx context.Context, id string) (*model.Song, bool, error)
}

type MusicCache interface {
	IsKeyExist(ctx context.Context, key string) bool
	SetSongsCache(ctx context.Context, key string, song *model.Song) error
	GetSongsCache(ctx context.Context, key string) (*model.Song, error)
	SetPlaylistDetailCache(ctx context.Context, key string, playlist *types.PlaylistDetail) error
	GetPlaylistDetailCache(ctx context.Context, key string) (*types.PlaylistDetail, error)
}
