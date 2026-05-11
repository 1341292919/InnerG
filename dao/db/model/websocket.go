package model

import (
	"InnerG/pkg/constants"
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID        int64          `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	MsgID     string         `gorm:"type:varchar(64);not null;uniqueIndex;column:msg_id" json:"msg_id"`
	FromUser  int64          `gorm:"type:bigint;not null;column:from_user;index:idx_from_to,priority:1" json:"from_user"`
	ToUser    int64          `gorm:"type:bigint;not null;column:to_user;index:idx_from_to,priority:2;index:idx_to_from,priority:1" json:"to_user"`
	Content   string         `gorm:"type:text;not null;column:content" json:"content"`
	Type      int8           `gorm:"type:tinyint;not null;default:1;column:type" json:"type"`
	Status    int8           `gorm:"type:tinyint;not null;default:0;column:status" json:"status"`
	CreatedAt int64          `gorm:"type:bigint;not null;column:created_at;index:idx_created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP;autoUpdateTime;column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (Message) TableName() string {
	return constants.MessageTableName
}
