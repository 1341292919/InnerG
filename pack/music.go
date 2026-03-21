package pack

import (
	"InnerG/dao/db/model"
	"InnerG/types"
	"database/sql"
	"strconv"
)

func BuildPlaylistList(list []*model.Playlist) []*types.Playlist {
	res := make([]*types.Playlist, 0)
	for _, item := range list {
		res = append(res, &types.Playlist{
			Id:          strconv.FormatInt(int64(item.ID), 10),
			Name:        item.Name,
			Description: convertNullString(item.Description),
			CoverUrl:    convertNullString(item.CoverURL),
			Status:      int(item.Status),
			PlayCount:   item.PlayCount,
			SongCount:   item.SongCount,
			Tags:        convertNullString(item.Tags),
			CreatedAt:   item.CreatedAt.Unix(),
			UpdatedAt:   item.UpdatedAt.Unix(),
		})
	}
	return res
}

func BuildPlaylistDetail(data *model.Playlist, songs []*model.PlaylistSong) *types.PlaylistDetail {
	return &types.PlaylistDetail{
		Id:          strconv.FormatInt(int64(data.ID), 10),
		Name:        data.Name,
		Description: convertNullString(data.Description),
		CoverUrl:    convertNullString(data.CoverURL),
		Status:      int(data.Status),
		PlayCount:   data.PlayCount,
		SongCount:   data.SongCount,
		Tags:        convertNullString(data.Tags),
		Songs:       BuildPlaylistSongList(songs),
		CreatedAt:   data.CreatedAt.Unix(),
		UpdatedAt:   data.UpdatedAt.Unix(),
	}
}

func BuildPlaylistSongList(list []*model.PlaylistSong) []*types.Song {
	res := make([]*types.Song, 0)
	for _, item := range list {
		res = append(res, &types.Song{
			Id:         strconv.FormatUint(item.ID, 10),
			Name:       item.Name,
			SingerName: item.SingerName,
			CreatedAt:  item.CreatedAt.Unix(),
		})
	}
	return res
}

func BuildSongDetail(data *model.Song) *types.SongDetail {
	return &types.SongDetail{
		Id:          strconv.FormatInt(int64(data.ID), 10),
		Name:        data.Name,
		Description: convertNullString(data.Description),
		CoverUrl:    convertNullString(data.CoverURL),
		Status:      int(data.Status),
		SingerId:    strconv.FormatUint(data.SingerID, 10),
		SourceUrl:   data.SourceURL,
		Duration:    data.Duration,
		PlayCount:   data.PlayCount,
		Tags:        convertNullString(data.Tags),
		CreatedAt:   data.CreatedAt.Unix(),
		UpdatedAt:   data.UpdatedAt.Unix(),
	}
}

func BuildSongDetailList(list []*model.Song) []*types.SongDetail {
	res := make([]*types.SongDetail, 0)
	for _, item := range list {
		res = append(res, BuildSongDetail(item))
	}
	return res
}

func convertNullString(field sql.NullString) string {
	if !field.Valid {
		return ""
	}
	return field.String
}
