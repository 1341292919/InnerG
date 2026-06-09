package service

import (
	"InnerG/config"
	"InnerG/dao"
	"InnerG/dao/db/model"
	"InnerG/pkg/constants"
	"InnerG/pkg/errno"
	"InnerG/pkg/utils"
	"context"
	"fmt"
	"sync"
	"time"
)

type GroupService struct {
	sf *utils.Snowflake
}

var groupServiceIns *GroupService
var groupServiceOnce sync.Once

func GetGroupService() *GroupService {
	groupServiceOnce.Do(func() {
		sf, _ := utils.NewSnowflake(config.Snowflake.DatancenterID, config.Snowflake.WorkerID)
		groupServiceIns = &GroupService{sf: sf}
	})
	return groupServiceIns
}

func (s *GroupService) CreateGroup(ctx context.Context, ownerID int64, name, description string, memberIDs []int64) (*model.Group, error) {
	groupDao := dao.NewGroupDao()
	members := uniqueGroupMemberIDs(memberIDs, ownerID)

	groupID, err := s.sf.NextVal()
	if err != nil {
		return nil, err
	}
	now := time.Now().Unix()

	group := &model.Group{
		GroupID:     groupID,
		Name:        name,
		OwnerID:     ownerID,
		Description: description,
		MemberCount: len(members) + 1,
		MaxMembers:  constants.DefaultMaxGroupMembers,
		CreatedAt:   now,
	}

	if err = groupDao.Db.CreateGroup(ctx, group); err != nil {
		return nil, errno.InternalServiceError.WithMessage("创建群组失败")
	}

	if err = groupDao.Db.AddGroupMember(ctx, &model.GroupMember{
		GroupID:  groupID,
		UserID:   ownerID,
		Role:     constants.GroupRoleOwner,
		JoinedAt: now,
	}); err != nil {
		return nil, errno.InternalServiceError.WithMessage("添加群主失败")
	}

	for _, memberID := range members {
		if err = groupDao.Db.AddGroupMember(ctx, &model.GroupMember{
			GroupID:  groupID,
			UserID:   memberID,
			Role:     constants.GroupRoleMember,
			JoinedAt: now,
		}); err != nil {
			return nil, errno.InternalServiceError.WithMessage("添加群成员失败")
		}
	}

	_ = groupDao.Cache.DeleteGroupMembers(ctx, groupID)

	return group, nil
}

func (s *GroupService) GetUserGroups(ctx context.Context, userID int64) ([]*model.Group, error) {
	groupDao := dao.NewGroupDao()
	return groupDao.Db.GetUserGroups(ctx, userID)
}

func (s *GroupService) GetGroupByID(ctx context.Context, groupID int64) (*model.Group, error) {
	groupDao := dao.NewGroupDao()
	return groupDao.Db.GetGroupByID(ctx, groupID)
}

func (s *GroupService) UpdateGroup(ctx context.Context, groupID, userID int64, name, avatar, description *string) error {
	groupDao := dao.NewGroupDao()

	group, err := groupDao.Db.GetGroupByID(ctx, groupID)
	if err != nil {
		return errno.ParamMissing.WithMessage("群组不存在")
	}
	if group.OwnerID != userID {
		return errno.AuthError.WithMessage("只有群主可以修改群信息")
	}

	if name != nil {
		group.Name = *name
	}
	if avatar != nil {
		group.Avatar = *avatar
	}
	if description != nil {
		group.Description = *description
	}

	return groupDao.Db.UpdateGroup(ctx, group)
}

func (s *GroupService) DeleteGroup(ctx context.Context, groupID, userID int64) error {
	groupDao := dao.NewGroupDao()

	group, err := groupDao.Db.GetGroupByID(ctx, groupID)
	if err != nil {
		return errno.ParamMissing.WithMessage("群组不存在")
	}
	if group.OwnerID != userID {
		return errno.AuthError.WithMessage("只有群主可以解散群组")
	}

	if err = groupDao.Db.DeleteGroup(ctx, groupID); err != nil {
		return errno.InternalServiceError.WithMessage("解散群组失败")
	}

	_ = groupDao.Cache.DeleteGroupMembers(ctx, groupID)

	return nil
}

func (s *GroupService) AddGroupMembers(ctx context.Context, groupID, inviterID int64, userIDs []int64) error {
	groupDao := dao.NewGroupDao()
	userIDs = uniqueGroupMemberIDs(userIDs, 0)

	isMember, err := groupDao.Db.IsMember(ctx, groupID, inviterID)
	if err != nil || !isMember {
		return errno.AuthError.WithMessage("您不是群成员")
	}

	group, err := groupDao.Db.GetGroupByID(ctx, groupID)
	if err != nil {
		return errno.ParamMissing.WithMessage("群组不存在")
	}

	if group.MemberCount+len(userIDs) > group.MaxMembers {
		return errno.ParamMissing.WithMessage(fmt.Sprintf("群组成员已达上限（%d人）", group.MaxMembers))
	}

	now := time.Now().Unix()
	addedCount := 0

	for _, userID := range userIDs {
		isExist, err := groupDao.Db.IsMember(ctx, groupID, userID)
		if err != nil {
			return err
		}
		if isExist {
			continue
		}

		if err = groupDao.Db.AddGroupMember(ctx, &model.GroupMember{
			GroupID:  groupID,
			UserID:   userID,
			Role:     constants.GroupRoleMember,
			JoinedAt: now,
		}); err != nil {
			return errno.InternalServiceError.WithMessage("添加群成员失败")
		}
		addedCount++
	}

	if addedCount > 0 {
		if err = groupDao.Db.IncrementMemberCount(ctx, groupID, addedCount); err != nil {
			return errno.InternalServiceError.WithMessage("更新群成员数量失败")
		}
		_ = groupDao.Cache.DeleteGroupMembers(ctx, groupID)
	}

	return nil
}

func (s *GroupService) RemoveGroupMember(ctx context.Context, groupID, operatorID, targetUserID int64) error {
	groupDao := dao.NewGroupDao()

	group, err := groupDao.Db.GetGroupByID(ctx, groupID)
	if err != nil {
		return errno.ParamMissing.WithMessage("群组不存在")
	}

	if group.OwnerID != operatorID {
		return errno.AuthError.WithMessage("只有群主可以移除成员")
	}

	if targetUserID == group.OwnerID {
		return errno.ParamMissing.WithMessage("不能移除群主")
	}

	if err = groupDao.Db.RemoveGroupMember(ctx, groupID, targetUserID); err != nil {
		return errno.InternalServiceError.WithMessage("移除成员失败")
	}

	if err = groupDao.Db.IncrementMemberCount(ctx, groupID, -1); err != nil {
		return errno.InternalServiceError.WithMessage("更新群成员数量失败")
	}
	_ = groupDao.Cache.DeleteGroupMembers(ctx, groupID)

	return nil
}

func (s *GroupService) QuitGroup(ctx context.Context, groupID, userID int64) error {
	groupDao := dao.NewGroupDao()

	group, err := groupDao.Db.GetGroupByID(ctx, groupID)
	if err != nil {
		return errno.ParamMissing.WithMessage("群组不存在")
	}

	if group.OwnerID == userID {
		return errno.ParamMissing.WithMessage("群主不能退出群组，请先转让群主或解散群组")
	}

	if err = groupDao.Db.RemoveGroupMember(ctx, groupID, userID); err != nil {
		return errno.InternalServiceError.WithMessage("退出群组失败")
	}

	if err = groupDao.Db.IncrementMemberCount(ctx, groupID, -1); err != nil {
		return errno.InternalServiceError.WithMessage("更新群成员数量失败")
	}
	_ = groupDao.Cache.DeleteGroupMembers(ctx, groupID)

	return nil
}

func (s *GroupService) GetGroupMembers(ctx context.Context, groupID int64, page, pageSize int) ([]*model.GroupMember, int64, error) {
	groupDao := dao.NewGroupDao()
	return groupDao.Db.GetGroupMembers(ctx, groupID, page, pageSize)
}

func (s *GroupService) GetGroupMessages(ctx context.Context, groupID, userID int64, before int64, page, pageSize int) ([]*model.GroupMessage, int64, error) {
	groupDao := dao.NewGroupDao()

	isMember, err := groupDao.Db.IsMember(ctx, groupID, userID)
	if err != nil || !isMember {
		return nil, 0, errno.AuthError.WithMessage("您不是群成员")
	}

	return groupDao.Db.GetGroupMessages(ctx, groupID, before, page, pageSize)
}

func uniqueGroupMemberIDs(ids []int64, excludedID int64) []int64 {
	seen := make(map[int64]struct{}, len(ids))
	unique := make([]int64, 0, len(ids))
	for _, id := range ids {
		if id == 0 || id == excludedID {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		unique = append(unique, id)
	}
	return unique
}
