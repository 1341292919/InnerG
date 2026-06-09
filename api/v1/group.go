package v1

import (
	"InnerG/dao"
	"InnerG/dao/db/model"
	"InnerG/pack"
	"InnerG/pkg/ctl"
	"InnerG/pkg/errno"
	"InnerG/service"
	"InnerG/types"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateGroup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.CreateGroupReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			pack.RespError(ctx, errno.ParamMissing.WithMessage(err.Error()))
			return
		}

		userInfo := ctl.GetUserInfo(ctx)
		groupSrv := service.GetGroupService()

		group, err := groupSrv.CreateGroup(ctx.Request.Context(), userInfo.Id, req.Name, req.Description, req.MemberIDs)
		if err != nil {
			pack.RespError(ctx, err)
			return
		}

		pack.RespData(ctx, types.CreateGroupResp{
			GroupID:     group.GroupID,
			Name:        group.Name,
			OwnerID:     group.OwnerID,
			MemberCount: group.MemberCount,
			CreatedAt:   group.CreatedAt,
		})
	}
}

func GetGroups() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userInfo := ctl.GetUserInfo(ctx)
		groupSrv := service.GetGroupService()

		groups, err := groupSrv.GetUserGroups(ctx.Request.Context(), userInfo.Id)
		if err != nil {
			pack.RespError(ctx, err)
			return
		}

		groupInfos := make([]*types.GroupInfo, len(groups))
		for i, group := range groups {
			groupInfos[i] = &types.GroupInfo{
				GroupID:     group.GroupID,
				Name:        group.Name,
				Avatar:      group.Avatar,
				OwnerID:     group.OwnerID,
				Description: group.Description,
				MemberCount: group.MemberCount,
				MaxMembers:  group.MaxMembers,
				CreatedAt:   group.CreatedAt,
			}
		}

		pack.RespData(ctx, types.GetGroupsResp{
			Groups: groupInfos,
			Total:  len(groupInfos),
		})
	}
}

func GetGroupDetail() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		groupIDStr := ctx.Param("group_id")
		groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
		if err != nil {
			pack.RespError(ctx, errno.ParamMissing.WithMessage("群组ID格式错误"))
			return
		}

		groupSrv := service.GetGroupService()
		group, err := groupSrv.GetGroupByID(ctx.Request.Context(), groupID)
		if err != nil {
			pack.RespError(ctx, err)
			return
		}

		pack.RespData(ctx, types.GroupInfo{
			GroupID:     group.GroupID,
			Name:        group.Name,
			Avatar:      group.Avatar,
			OwnerID:     group.OwnerID,
			Description: group.Description,
			MemberCount: group.MemberCount,
			MaxMembers:  group.MaxMembers,
			CreatedAt:   group.CreatedAt,
		})
	}
}

func UpdateGroup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		groupIDStr := ctx.Param("group_id")
		groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
		if err != nil {
			pack.RespError(ctx, errno.ParamMissing.WithMessage("群组ID格式错误"))
			return
		}

		var req types.UpdateGroupReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			pack.RespError(ctx, errno.ParamMissing.WithMessage(err.Error()))
			return
		}

		userInfo := ctl.GetUserInfo(ctx)
		groupSrv := service.GetGroupService()

		err = groupSrv.UpdateGroup(ctx.Request.Context(), groupID, userInfo.Id, req.Name, req.Avatar, req.Description)
		if err != nil {
			pack.RespError(ctx, err)
			return
		}

		pack.RespSuccess(ctx)
	}
}

func DeleteGroup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		groupIDStr := ctx.Param("group_id")
		groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
		if err != nil {
			pack.RespError(ctx, errno.ParamMissing.WithMessage("群组ID格式错误"))
			return
		}

		userInfo := ctl.GetUserInfo(ctx)
		groupSrv := service.GetGroupService()

		err = groupSrv.DeleteGroup(ctx.Request.Context(), groupID, userInfo.Id)
		if err != nil {
			pack.RespError(ctx, err)
			return
		}

		pack.RespSuccess(ctx)
	}
}

func AddGroupMembers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		groupIDStr := ctx.Param("group_id")
		groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
		if err != nil {
			pack.RespError(ctx, errno.ParamMissing.WithMessage("群组ID格式错误"))
			return
		}

		var req types.AddGroupMembersReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			pack.RespError(ctx, errno.ParamMissing.WithMessage(err.Error()))
			return
		}

		userInfo := ctl.GetUserInfo(ctx)
		groupSrv := service.GetGroupService()

		err = groupSrv.AddGroupMembers(ctx.Request.Context(), groupID, userInfo.Id, req.UserIDs)
		if err != nil {
			pack.RespError(ctx, err)
			return
		}

		pack.RespSuccess(ctx)
	}
}

func RemoveGroupMember() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		groupIDStr := ctx.Param("group_id")
		groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
		if err != nil {
			pack.RespError(ctx, errno.ParamMissing.WithMessage("群组ID格式错误"))
			return
		}

		userIDStr := ctx.Param("user_id")
		targetUserID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			pack.RespError(ctx, errno.ParamMissing.WithMessage("用户ID格式错误"))
			return
		}

		userInfo := ctl.GetUserInfo(ctx)
		groupSrv := service.GetGroupService()

		err = groupSrv.RemoveGroupMember(ctx.Request.Context(), groupID, userInfo.Id, targetUserID)
		if err != nil {
			pack.RespError(ctx, err)
			return
		}

		pack.RespSuccess(ctx)
	}
}

func QuitGroup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		groupIDStr := ctx.Param("group_id")
		groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
		if err != nil {
			pack.RespError(ctx, errno.ParamMissing.WithMessage("群组ID格式错误"))
			return
		}

		userInfo := ctl.GetUserInfo(ctx)
		groupSrv := service.GetGroupService()

		err = groupSrv.QuitGroup(ctx.Request.Context(), groupID, userInfo.Id)
		if err != nil {
			pack.RespError(ctx, err)
			return
		}

		pack.RespSuccess(ctx)
	}
}

func GetGroupMembers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		groupIDStr := ctx.Param("group_id")
		groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
		if err != nil {
			pack.RespError(ctx, errno.ParamMissing.WithMessage("群组ID格式错误"))
			return
		}

		page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "50"))
		if page < 1 {
			page = 1
		}
		if pageSize < 1 || pageSize > 100 {
			pageSize = 50
		}

		groupSrv := service.GetGroupService()
		members, total, err := groupSrv.GetGroupMembers(ctx.Request.Context(), groupID, page, pageSize)
		if err != nil {
			pack.RespError(ctx, err)
			return
		}

		users, err := getUserBasicMapByIDs(ctx, collectGroupMemberIDs(members))
		if err != nil {
			pack.RespError(ctx, err)
			return
		}

		memberInfos := make([]*types.GroupMemberInfo, len(members))
		for i, member := range members {
			user := users[member.UserID]
			memberInfos[i] = &types.GroupMemberInfo{
				UserID:   member.UserID,
				Username: usernameOf(user),
				Avatar:   avatarOf(user),
				Role:     member.Role,
				JoinedAt: member.JoinedAt,
			}
		}

		pack.RespData(ctx, types.GetGroupMembersResp{
			Members:   memberInfos,
			Total:     total,
			PageIndex: page,
			PageSize:  pageSize,
		})
	}
}

func GetGroupMessages() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.GetGroupMessagesReq
		if err := ctx.ShouldBindQuery(&req); err != nil {
			pack.RespError(ctx, errno.ParamMissing.WithMessage(err.Error()))
			return
		}

		if req.PageSize <= 0 || req.PageSize > 100 {
			req.PageSize = 20
		}
		if req.PageNum < 1 {
			req.PageNum = 1
		}

		userInfo := ctl.GetUserInfo(ctx)
		groupSrv := service.GetGroupService()

		messages, total, err := groupSrv.GetGroupMessages(ctx.Request.Context(), req.GroupID, userInfo.Id, req.Before, req.PageNum, req.PageSize)
		if err != nil {
			pack.RespError(ctx, err)
			return
		}

		users, err := getUserBasicMapByIDs(ctx, collectGroupMessageUserIDs(messages))
		if err != nil {
			pack.RespError(ctx, err)
			return
		}

		messageResps := make([]*types.GroupMessageResp, len(messages))
		for i, msg := range messages {
			fromName := "系统"
			fromAvatar := ""

			if msg.FromUser != 0 {
				user := users[msg.FromUser]
				if user != nil {
					fromName = user.Username
					fromAvatar = user.Avatar
				}
			}

			messageResps[i] = &types.GroupMessageResp{
				MsgID:      msg.MsgID,
				FromUser:   msg.FromUser,
				FromName:   fromName,
				FromAvatar: fromAvatar,
				Content:    msg.Content,
				Type:       msg.Type,
				CreatedAt:  msg.CreatedAt,
			}
		}

		pack.RespData(ctx, types.GetGroupMessagesResp{
			Messages:   messageResps,
			TotalCount: total,
			PageIndex:  req.PageNum,
			PageSize:   req.PageSize,
		})
	}
}

func getUserBasicMapByIDs(ctx *gin.Context, ids []int64) (map[int64]*model.User, error) {
	return dao.NewUserDao(ctx.Request.Context()).Db.GetUserBasicByIds(ctx.Request.Context(), uniqueInt64(ids))
}

func collectGroupMemberIDs(members []*model.GroupMember) []int64 {
	ids := make([]int64, 0, len(members))
	for _, member := range members {
		ids = append(ids, member.UserID)
	}
	return ids
}

func collectGroupMessageUserIDs(messages []*model.GroupMessage) []int64 {
	ids := make([]int64, 0, len(messages))
	for _, message := range messages {
		if message.FromUser != 0 {
			ids = append(ids, message.FromUser)
		}
	}
	return ids
}

func usernameOf(user *model.User) string {
	if user == nil {
		return ""
	}
	return user.Username
}

func avatarOf(user *model.User) string {
	if user == nil {
		return ""
	}
	return user.Avatar
}
