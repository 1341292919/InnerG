package model

import (
	"InnerG/pkg/constants"
	"time"

	"gorm.io/gorm"
)

type Friend struct {
	ID        int64          `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	UserID    int64          `gorm:"type:bigint;not null;column:user_id;uniqueIndex:idx_user_friend,priority:1;index:idx_user_status,priority:1" json:"user_id"`
	FriendID  int64          `gorm:"type:bigint;not null;column:friend_id;uniqueIndex:idx_user_friend,priority:2;index:idx_friend_status,priority:1" json:"friend_id"`
	Status    int8           `gorm:"type:tinyint;not null;default:1;column:status;index:idx_user_status,priority:2;index:idx_friend_status,priority:2" json:"status"`
	CreatedAt int64          `gorm:"type:bigint;not null;column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP;autoUpdateTime;column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (Friend) TableName() string {
	return constants.FriendTableName
}
