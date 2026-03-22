package types

type GetPlaylistListReq struct {
	PageSize int `form:"pageSize" binding:"required"`
	PageNum  int `form:"pageNum" binding:"required"`
}

type GetPlaylistListResp struct {
	Total        int
	PlaylistList []*Playlist
}

type GetPlaylistDetailReq struct {
	PlaylistId string `form:"playlistId" binding:"required"`
}

type GetPlaylistDetailResp struct {
	PlaylistDetail *PlaylistDetail
}

type GetSongDetailReq struct {
	SongId string `form:"songId" binding:"required"`
}

type GetSongDetailResp struct {
	SongDetail *SongDetail
}

type GetSongDetailListReq struct {
	PageSize int `form:"pageSize" binding:"required"`
	PageNum  int `form:"pageNum" binding:"required"`
}

type GetSongDetailListResp struct {
	Total    int
	SongList []*SongDetail
}

type Playlist struct {
	Id          string
	Name        string
	Description string
	CoverUrl    string
	Status      int
	PlayCount   int64
	SongCount   int
	Tags        string
	CreatedAt   int64
	UpdatedAt   int64
}

type PlaylistDetail struct {
	Id          string
	Name        string
	Description string
	CoverUrl    string
	Status      int
	PlayCount   int64
	SongCount   int
	Tags        string
	Songs       []*Song
	CreatedAt   int64
	UpdatedAt   int64
}

type Song struct {
	Id         string
	Name       string
	SingerName string
	CreatedAt  int64
}

type SongDetail struct {
	Id          string
	Name        string
	Description string
	CoverUrl    string
	Status      int
	SingerName  string
	Album       string
	SourceUrl   string
	Duration    int
	PlayCount   int64
	Tags        string
	CreatedAt   int64
	UpdatedAt   int64
}
