package service

import (
	"InnerG/dao"
	"InnerG/dao/db/model"
	"InnerG/pack"
	"InnerG/types"
	"context"
	"fmt"
	"log"
	"sync"
)

var MusicSrvIns *MusicSrv
var MusicSrvOnce sync.Once

type MusicSrv struct{}

func GetMusicSrv() *MusicSrv {
	MusicSrvOnce.Do(func() {
		MusicSrvIns = &MusicSrv{}
	})
	return MusicSrvIns
}

func (s *MusicSrv) GetPlaylistList(ctx context.Context, req *types.GetPlaylistListReq) ([]*types.Playlist, int, error) {
	musicDao := dao.NewMusicDao(ctx)
	list, total, err := musicDao.Db.GetPlaylistList(ctx, req.PageNum, req.PageSize)
	if err != nil {
		return nil, -1, err
	}
	return pack.BuildPlaylistList(list), total, nil
}

func (s *MusicSrv) GetPlaylistDetail(ctx context.Context, req *types.GetPlaylistDetailReq) (*types.PlaylistDetail, error) {
	musicDao := dao.NewMusicDao(ctx)
	key := fmt.Sprintf("playListDetail:%s", req.PlaylistId)
	var res *types.PlaylistDetail
	var err error
	if musicDao.Ca.IsKeyExist(ctx, key) {
		res, err = musicDao.Ca.GetPlaylistDetailCache(ctx, key)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
	data, exist, err := musicDao.Db.GetPlaylistById(ctx, req.PlaylistId)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("playlist not exist")
	}
	songs, err := musicDao.Db.GetPlaylistSongListByPlaylistId(ctx, req.PlaylistId)
	if err != nil {
		return nil, err
	}
	res = pack.BuildPlaylistDetail(data, songs)
	go func() {
		err = musicDao.Ca.SetPlaylistDetailCache(context.Background(), key, res)
		if err != nil {
			log.Println(err)
		}
	}()
	return res, nil
}

func (s *MusicSrv) GetSongDetail(ctx context.Context, req *types.GetSongDetailReq) (*types.SongDetail, error) {
	musicDao := dao.NewMusicDao(ctx)
	key := fmt.Sprintf("song:%s", req.SongId)
	var res *model.Song
	var err error
	if musicDao.Ca.IsKeyExist(ctx, key) {
		res, err = musicDao.Ca.GetSongsCache(ctx, key)
		if err != nil {
			return nil, err
		}
		return pack.BuildSongDetail(res), nil
	}
	data, exist, err := musicDao.Db.GetSongById(ctx, req.SongId)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("song not exist")
	}
	go func() {
		err = musicDao.Ca.SetSongsCache(context.Background(), key, data)
		if err != nil {
			log.Println(err)
		}
	}()
	return pack.BuildSongDetail(data), nil
}

func (s *MusicSrv) GetSongDetailList(ctx context.Context, req *types.GetSongDetailListReq) ([]*types.SongDetail, int, error) {
	musicDao := dao.NewMusicDao(ctx)
	list, total, err := musicDao.Db.GetSongList(ctx, req.PageNum, req.PageSize)
	if err != nil {
		return nil, -1, err
	}
	return pack.BuildSongDetailList(list), total, nil
}
