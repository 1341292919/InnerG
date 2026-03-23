package cache

import (
	"InnerG/dao/db/model"
	_interface "InnerG/dao/interface"
	"InnerG/pkg/constants"
	"InnerG/types"
	"context"
	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
)

type musicCache struct {
	client *redis.Client
}

func NewMusicCache(cache *redis.Client) _interface.MusicCache {
	return &musicCache{
		client: cache,
	}
}
func (ca *musicCache) IsKeyExist(ctx context.Context, key string) bool {
	return ca.client.Exists(ctx, key).Val() == 1
}
func (ca *musicCache) SetSongsCache(ctx context.Context, key string, song *model.Song) error {
	songsJson, err := sonic.Marshal(song)
	if err != nil {
		return err
	}
	return ca.client.Set(ctx, key, songsJson, constants.SongsKeyExpire).Err()
}

func (ca *musicCache) GetSongsCache(ctx context.Context, key string) (*model.Song, error) {
	songsJson, err := ca.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var songs *model.Song
	err = sonic.Unmarshal([]byte(songsJson), &songs)
	if err != nil {
		return nil, err
	}
	return songs, nil
}

func (ca *musicCache) SetPlaylistDetailCache(ctx context.Context, key string, playlist *types.PlaylistDetail) error {
	playlistJson, err := sonic.Marshal(playlist)
	if err != nil {
		return err
	}
	return ca.client.Set(ctx, key, playlistJson, constants.PlaylistKeyExpire).Err()
}
func (ca *musicCache) GetPlaylistDetailCache(ctx context.Context, key string) (*types.PlaylistDetail, error) {
	playlistJson, err := ca.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var playlist *types.PlaylistDetail
	err = sonic.Unmarshal([]byte(playlistJson), &playlist)
	if err != nil {
		return nil, err
	}
	return playlist, nil
}
