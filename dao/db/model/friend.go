package model

import (
	"InnerG/pkg/constants"
	"time"

	"gorm.io/gorm"
)

type Friend struct {
	ID        int64          `gorm:"primaryKey;column:id" json:"id"`
	UserID    int64          `gorm:"type:bigint;not null;column:user_id;uniqueIndex:idx_user_friend,priority:1;index:idx_user_status,priority:1" json:"user_id"`
	FriendID  int64          `gorm:"type:bigint;not null;column:friend_id;uniqueIndex:idx_user_friend,priority:2" json:"friend_id"`
	Status    int8           `gorm:"type:tinyint;not null;default:1;column:status;index:idx_user_status,priority:2" json:"status"`
	CreatedAt int64          `gorm:"type:bigint;not null;column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP;autoUpdateTime;column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (Friend) TableName() string {
	return constants.FriendTableName
}

type FriendRequest struct {
	ID        int64          `gorm:"primaryKey;column:id" json:"id"`
	FromUser  int64          `gorm:"type:bigint;not null;column:from_user;uniqueIndex:idx_from_to,priority:1;index:idx_from_status,priority:1" json:"from_user"`
	ToUser    int64          `gorm:"type:bigint;not null;column:to_user;uniqueIndex:idx_from_to,priority:2;index:idx_to_status,priority:1" json:"to_user"`
	Status    int8           `gorm:"type:tinyint;not null;default:1;column:status;index:idx_from_status,priority:2;index:idx_to_status,priority:2" json:"status"`
	CreatedAt int64          `gorm:"type:bigint;not null;column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP;autoUpdateTime;column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (FriendRequest) TableName() string {
	return constants.FriendRequestTableName
}
