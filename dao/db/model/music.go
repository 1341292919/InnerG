package model

import (
	"database/sql"
	"gorm.io/gorm"
	"time"
)

type Song struct {
	gorm.Model
	Name        string
	Description sql.NullString
	CoverURL    sql.NullString `gorm:"column:cover_url"`
	Status      int8
	SingerID    uint64 `gorm:"column:singer_id"`
	SourceURL   string `gorm:"column:source_url"`
	Duration    int
	PlayCount   int64 `gorm:"column:play_count"`
	Tags        sql.NullString
}

type Playlist struct {
	gorm.Model
	Name        string
	Description sql.NullString
	CoverURL    sql.NullString `gorm:"column:cover_url"`
	Status      int8
	PlayCount   int64 `gorm:"column:play_count"`
	SongCount   int   `gorm:"column:song_count"`
	Tags        sql.NullString
}

type PlaylistSong struct {
	ID         uint64    `gorm:"column:id"`
	Name       string    `gorm:"column:name"`
	SingerName string    `gorm:"column:singer_name"`
	CreatedAt  time.Time `gorm:"column:created_at"`
}
