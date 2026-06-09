package model

import (
	"InnerG/pkg/constants"
	"time"

	"gorm.io/gorm"
)

type Group struct {
	ID           int64          `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	GroupID      int64          `gorm:"type:bigint;not null;uniqueIndex;column:group_id" json:"group_id"`
	Name         string         `gorm:"type:varchar(100);not null;column:name" json:"name"`
	Avatar       string         `gorm:"type:varchar(512);column:avatar" json:"avatar"`
	OwnerID      int64          `gorm:"type:bigint;not null;column:owner_id;index:idx_owner" json:"owner_id"`
	Description  string         `gorm:"type:varchar(500);default:'';column:description" json:"description"`
	MemberCount  int            `gorm:"type:int;default:0;not null;column:member_count" json:"member_count"`
	MaxMembers   int            `gorm:"type:int;default:200;not null;column:max_members" json:"max_members"`
	Announcement string         `gorm:"type:text;column:announcement" json:"announcement"`
	JoinMode     int8           `gorm:"type:tinyint;default:0;not null;column:join_mode" json:"join_mode"`
	MuteAll      int8           `gorm:"type:tinyint;default:0;not null;column:mute_all" json:"mute_all"`
	Settings     string         `gorm:"type:json;column:settings" json:"settings"`
	CreatedAt    int64          `gorm:"type:bigint;not null;column:created_at;index:idx_created_at" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP;autoUpdateTime;column:updated_at" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (Group) TableName() string {
	return constants.GroupTableName
}

type GroupMember struct {
	ID        int64          `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	GroupID   int64          `gorm:"type:bigint;not null;column:group_id;uniqueIndex:uk_group_user,priority:1;index:idx_group_members" json:"group_id"`
	UserID    int64          `gorm:"type:bigint;not null;column:user_id;uniqueIndex:uk_group_user,priority:2;index:idx_user_groups" json:"user_id"`
	JoinedAt  int64          `gorm:"type:bigint;not null;column:joined_at" json:"joined_at"`
	Role      int8           `gorm:"type:tinyint;default:0;not null;column:role" json:"role"`
	Nickname  string         `gorm:"type:varchar(50);column:nickname" json:"nickname"`
	Muted     int8           `gorm:"type:tinyint;default:0;not null;column:muted" json:"muted"`
	MuteUntil int64          `gorm:"type:bigint;column:mute_until" json:"mute_until"`
	CreatedAt time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP;column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP;autoUpdateTime;column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (GroupMember) TableName() string {
	return constants.GroupMemberTableName
}

type GroupMessage struct {
	ID             int64          `gorm:"primaryKey;column:id" json:"id"`
	MsgID          string         `gorm:"type:varchar(64);not null;uniqueIndex;column:msg_id;index:idx_msg_id" json:"msg_id"`
	GroupID        int64          `gorm:"type:bigint;not null;column:group_id;index:idx_group_id,priority:1;index:idx_group_time,priority:1" json:"group_id"`
	FromUser       int64          `gorm:"type:bigint;not null;column:from_user" json:"from_user"`
	Content        string         `gorm:"type:text;not null;column:content" json:"content"`
	Type           int8           `gorm:"type:tinyint;not null;default:1;column:type" json:"type"`
	Status         int8           `gorm:"type:tinyint;not null;default:0;column:status" json:"status"`
	RecalledAt     int64          `gorm:"type:bigint;column:recalled_at" json:"recalled_at"`
	MentionedUsers string         `gorm:"type:json;column:mentioned_users" json:"mentioned_users"`
	CreatedAt      int64          `gorm:"type:bigint;not null;column:created_at;index:idx_group_time,priority:2;index:idx_created_at" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP;autoUpdateTime;column:updated_at" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (GroupMessage) TableName() string {
	return constants.GroupMessageTableName
}
