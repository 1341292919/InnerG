package _interface

import (
	"InnerG/dao/db/model"
	"context"
	"time"
)

type GroupDB interface {
	CreateGroup(ctx context.Context, group *model.Group) error
	GetGroupByID(ctx context.Context, groupID int64) (*model.Group, error)
	GetGroupByPK(ctx context.Context, id int64) (*model.Group, error)
	UpdateGroup(ctx context.Context, group *model.Group) error
	DeleteGroup(ctx context.Context, groupID int64) error
	GetUserGroups(ctx context.Context, userID int64) ([]*model.Group, error)
	IncrementMemberCount(ctx context.Context, groupID int64, delta int) error

	AddGroupMember(ctx context.Context, member *model.GroupMember) error
	RemoveGroupMember(ctx context.Context, groupID, userID int64) error
	GetGroupMembers(ctx context.Context, groupID int64, page, pageSize int) ([]*model.GroupMember, int64, error)
	GetGroupMemberIDs(ctx context.Context, groupID int64) ([]int64, error)
	IsMember(ctx context.Context, groupID, userID int64) (bool, error)
	GetMemberCount(ctx context.Context, groupID int64) (int, error)
	GetUserGroupIDs(ctx context.Context, userID int64) ([]int64, error)
	GetGroupMember(ctx context.Context, groupID, userID int64) (*model.GroupMember, error)

	InsertGroupMessage(ctx context.Context, msg *model.GroupMessage) error
	GetGroupMessages(ctx context.Context, groupID int64, before int64, page, pageSize int) ([]*model.GroupMessage, int64, error)
	GetGroupMessageByMsgID(ctx context.Context, msgID string) (*model.GroupMessage, error)
	GetGroupMessagesAfterID(ctx context.Context, userID, lastID int64, limit int) ([]*model.GroupMessage, error)
}

type GroupCache interface {
	SetGroupMembers(ctx context.Context, groupID int64, userIDs []int64, ttl time.Duration) error
	GetGroupMembers(ctx context.Context, groupID int64) ([]int64, error)
	DeleteGroupMembers(ctx context.Context, groupID int64) error
}
