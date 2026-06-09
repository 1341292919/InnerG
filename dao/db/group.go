package db

import (
	"InnerG/dao/db/model"
	"InnerG/pkg/constants"
	"context"

	"gorm.io/gorm"
)

type GroupDB struct {
	db *gorm.DB
}

func NewGroupDB(db *gorm.DB) *GroupDB {
	return &GroupDB{db: db}
}

func (g *GroupDB) CreateGroup(ctx context.Context, group *model.Group) error {
	return g.db.WithContext(ctx).Create(group).Error
}

func (g *GroupDB) GetGroupByID(ctx context.Context, groupID int64) (*model.Group, error) {
	var group model.Group
	err := g.db.WithContext(ctx).Where("group_id = ?", groupID).First(&group).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (g *GroupDB) GetGroupByPK(ctx context.Context, id int64) (*model.Group, error) {
	var group model.Group
	err := g.db.WithContext(ctx).First(&group, id).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (g *GroupDB) UpdateGroup(ctx context.Context, group *model.Group) error {
	return g.db.WithContext(ctx).Model(group).Updates(group).Error
}

func (g *GroupDB) DeleteGroup(ctx context.Context, groupID int64) error {
	return g.db.WithContext(ctx).Where("group_id = ?", groupID).Delete(&model.Group{}).Error
}

func (g *GroupDB) GetUserGroups(ctx context.Context, userID int64) ([]*model.Group, error) {
	var groups []*model.Group
	err := g.db.WithContext(ctx).
		Table(constants.GroupTableName+" g").
		Joins("INNER JOIN "+constants.GroupMemberTableName+" gm ON g.group_id = gm.group_id").
		Where("gm.user_id = ? AND gm.deleted_at IS NULL", userID).
		Select("g.*").
		Order("g.created_at DESC").
		Find(&groups).Error
	return groups, err
}

func (g *GroupDB) IncrementMemberCount(ctx context.Context, groupID int64, delta int) error {
	return g.db.WithContext(ctx).
		Model(&model.Group{}).
		Where("group_id = ?", groupID).
		Update("member_count", gorm.Expr("member_count + ?", delta)).Error
}

func (g *GroupDB) AddGroupMember(ctx context.Context, member *model.GroupMember) error {
	return g.db.WithContext(ctx).Create(member).Error
}

func (g *GroupDB) RemoveGroupMember(ctx context.Context, groupID, userID int64) error {
	return g.db.WithContext(ctx).
		Where("group_id = ? AND user_id = ?", groupID, userID).
		Delete(&model.GroupMember{}).Error
}

func (g *GroupDB) GetGroupMembers(ctx context.Context, groupID int64, page, pageSize int) ([]*model.GroupMember, int64, error) {
	var members []*model.GroupMember
	var total int64

	db := g.db.WithContext(ctx).Model(&model.GroupMember{}).Where("group_id = ?", groupID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := db.Order("role DESC, joined_at ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&members).Error

	return members, total, err
}

func (g *GroupDB) GetGroupMemberIDs(ctx context.Context, groupID int64) ([]int64, error) {
	var userIDs []int64
	err := g.db.WithContext(ctx).
		Model(&model.GroupMember{}).
		Where("group_id = ?", groupID).
		Pluck("user_id", &userIDs).Error
	return userIDs, err
}

func (g *GroupDB) IsMember(ctx context.Context, groupID, userID int64) (bool, error) {
	var count int64
	err := g.db.WithContext(ctx).
		Model(&model.GroupMember{}).
		Where("group_id = ? AND user_id = ?", groupID, userID).
		Count(&count).Error
	return count > 0, err
}

func (g *GroupDB) GetMemberCount(ctx context.Context, groupID int64) (int, error) {
	var count int64
	err := g.db.WithContext(ctx).
		Model(&model.GroupMember{}).
		Where("group_id = ?", groupID).
		Count(&count).Error
	return int(count), err
}

func (g *GroupDB) GetUserGroupIDs(ctx context.Context, userID int64) ([]int64, error) {
	var groupIDs []int64
	err := g.db.WithContext(ctx).
		Model(&model.GroupMember{}).
		Where("user_id = ?", userID).
		Pluck("group_id", &groupIDs).Error
	return groupIDs, err
}

func (g *GroupDB) GetGroupMember(ctx context.Context, groupID, userID int64) (*model.GroupMember, error) {
	var member model.GroupMember
	err := g.db.WithContext(ctx).
		Where("group_id = ? AND user_id = ?", groupID, userID).
		First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (g *GroupDB) InsertGroupMessage(ctx context.Context, msg *model.GroupMessage) error {
	return g.db.WithContext(ctx).Create(msg).Error
}

func (g *GroupDB) GetGroupMessages(ctx context.Context, groupID int64, before int64, page, pageSize int) ([]*model.GroupMessage, int64, error) {
	var messages []*model.GroupMessage
	var total int64

	db := g.db.WithContext(ctx).Model(&model.GroupMessage{}).
		Where("group_id = ? AND status != ?", groupID, constants.GroupMessageStatusDeleted)

	if before > 0 {
		db = db.Where("created_at < ?", before)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := db.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&messages).Error

	return messages, total, err
}

func (g *GroupDB) GetGroupMessageByMsgID(ctx context.Context, msgID string) (*model.GroupMessage, error) {
	var msg model.GroupMessage
	err := g.db.WithContext(ctx).Where("msg_id = ?", msgID).First(&msg).Error
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

// GetGroupMessagesAfterID 获取用户所在群的增量消息（基于雪花ID游标）
func (g *GroupDB) GetGroupMessagesAfterID(ctx context.Context, userID, lastID int64, limit int) ([]*model.GroupMessage, error) {
	var messages []*model.GroupMessage

	err := g.db.WithContext(ctx).
		Table(constants.GroupMessageTableName+" gm").
		Joins("INNER JOIN "+constants.GroupMemberTableName+" gm2 ON gm.group_id = gm2.group_id").
		Where("gm2.user_id = ? AND gm.id > ? AND gm.status = ?",
			userID, lastID, constants.GroupMessageStatusNormal).
		Select("gm.*").
		Order("gm.id ASC").
		Limit(limit + 1). // 多查一条用于判断 has_more
		Find(&messages).Error

	if err != nil {
		return nil, err
	}
	return messages, nil
}
