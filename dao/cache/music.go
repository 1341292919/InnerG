package cache

import (
	"InnerG/dao/db/model"
	_interface "InnerG/dao/interface"
	"InnerG/pkg/constants"
	"InnerG/pkg/errno"
	"InnerG/types"
	"context"
	"fmt"
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
		return errno.NewErr(errno.RedisDBErrorCode, fmt.Sprintf("SetSongsCache Marshal: %v", err))
	}
	if err := ca.client.Set(ctx, key, songsJson, constants.SongsKeyExpire).Err(); err != nil {
		return errno.NewErr(errno.RedisDBErrorCode, fmt.Sprintf("SetSongsCache Set: %v", err))
	}
	return nil
}

func (ca *musicCache) GetSongsCache(ctx context.Context, key string) (*model.Song, error) {
	songsJson, err := ca.client.Get(ctx, key).Result()
	if err != nil {
		return nil, errno.NewErr(errno.RedisDBErrorCode, fmt.Sprintf("GetSongsCache Get: %v", err))
	}
	var songs *model.Song
	err = sonic.Unmarshal([]byte(songsJson), &songs)
	if err != nil {
		return nil, errno.NewErr(errno.RedisDBErrorCode, fmt.Sprintf("GetSongsCache Unmarshal: %v", err))
	}
	return songs, nil
}

func (ca *musicCache) SetPlaylistDetailCache(ctx context.Context, key string, playlist *types.PlaylistDetail) error {
	playlistJson, err := sonic.Marshal(playlist)
	if err != nil {
		return errno.NewErr(errno.RedisDBErrorCode, fmt.Sprintf("SetPlaylistDetailCache Marshal: %v", err))
	}
	if err := ca.client.Set(ctx, key, playlistJson, constants.PlaylistKeyExpire).Err(); err != nil {
		return errno.NewErr(errno.RedisDBErrorCode, fmt.Sprintf("SetPlaylistDetailCache Set: %v", err))
	}
	return nil
}

func (ca *musicCache) GetPlaylistDetailCache(ctx context.Context, key string) (*types.PlaylistDetail, error) {
	playlistJson, err := ca.client.Get(ctx, key).Result()
	if err != nil {
		return nil, errno.NewErr(errno.RedisDBErrorCode, fmt.Sprintf("GetPlaylistDetailCache Get: %v", err))
	}
	var playlist *types.PlaylistDetail
	err = sonic.Unmarshal([]byte(playlistJson), &playlist)
	if err != nil {
		return nil, errno.NewErr(errno.RedisDBErrorCode, fmt.Sprintf("GetPlaylistDetailCache Unmarshal: %v", err))
	}
	return playlist, nil
}
