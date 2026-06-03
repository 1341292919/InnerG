package v1

import (
	"InnerG/dao"
	"InnerG/dao/db/model"
	"InnerG/pack"
	"InnerG/pkg/ctl"
	"InnerG/pkg/errno"
	"InnerG/service"
	"InnerG/types"

	"github.com/gin-gonic/gin"
)

func SendFriendRequest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.FriendTargetReq
		if err := ctx.ShouldBind(&req); err != nil {
			pack.RespError(ctx, errno.ParamMissing.WithMessage(err.Error()))
			return
		}

		userID := ctl.GetUserInfo(ctx.Request.Context()).Id
		if err := service.GetFriendSrv().CreateFriendRequest(ctx.Request.Context(), userID, req.FriendID, req.Message); err != nil {
			pack.RespError(ctx, err)
			return
		}

		pack.RespSuccess(ctx)
	}
}

func HandleFriendRequest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.HandleFriendRequestReq
		if err := ctx.ShouldBind(&req); err != nil {
			pack.RespError(ctx, errno.ParamMissing.WithMessage(err.Error()))
			return
		}

		userID := ctl.GetUserInfo(ctx.Request.Context()).Id
		if err := service.GetFriendSrv().HandleFriendRequest(ctx.Request.Context(), userID, req.FriendID, req.ActionType); err != nil {
			pack.RespError(ctx, err)
			return
		}

		pack.RespSuccess(ctx)
	}
}

func DeleteFriend() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.FriendTargetReq
		if err := ctx.ShouldBind(&req); err != nil {
			pack.RespError(ctx, errno.ParamMissing.WithMessage(err.Error()))
			return
		}

		userID := ctl.GetUserInfo(ctx.Request.Context()).Id
		if err := service.GetFriendSrv().DeleteFriend(ctx.Request.Context(), userID, req.FriendID); err != nil {
			pack.RespError(ctx, err)
			return
		}

		pack.RespSuccess(ctx)
	}
}

func ListFriends() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.FriendPageReq
		if err := ctx.ShouldBind(&req); err != nil {
			pack.RespError(ctx, errno.ParamMissing.WithMessage(err.Error()))
			return
		}
		if req.PageSize <= 0 {
			req.PageSize = 20
		}
		if req.PageNum <= 0 {
			req.PageNum = 1
		}

		userID := ctl.GetUserInfo(ctx.Request.Context()).Id
		friends, total, err := service.GetFriendSrv().ListFriends(ctx.Request.Context(), userID, req.PageSize, req.PageNum)
		if err != nil {
			pack.RespError(ctx, err)
			return
		}

		users, err := getUserBasicMap(ctx, collectFriendIDs(friends))
		if err != nil {
			pack.RespError(ctx, err)
			return
		}

		pack.RespData(ctx, types.FriendListResp{
			Friends:    pack.BuildFriendList(friends, users),
			TotalCount: total,
			PageIndex:  req.PageNum,
			PageSize:   req.PageSize,
		})
	}
}

func ListFriendRequests() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.FriendPageReq
		if err := ctx.ShouldBind(&req); err != nil {
			pack.RespError(ctx, errno.ParamMissing.WithMessage(err.Error()))
			return
		}
		if req.PageSize <= 0 {
			req.PageSize = 20
		}
		if req.PageNum <= 0 {
			req.PageNum = 1
		}

		userID := ctl.GetUserInfo(ctx.Request.Context()).Id
		requests, total, err := service.GetFriendSrv().ListInboundRequests(ctx.Request.Context(), userID, req.PageSize, req.PageNum)
		if err != nil {
			pack.RespError(ctx, err)
			return
		}

		users, err := getUserBasicMap(ctx, collectRequestFromUserIDs(requests))
		if err != nil {
			pack.RespError(ctx, err)
			return
		}

		pack.RespData(ctx, types.FriendRequestListResp{
			Requests:   pack.BuildFriendRequestList(requests, users),
			TotalCount: total,
			PageIndex:  req.PageNum,
			PageSize:   req.PageSize,
		})
	}
}

func getUserBasicMap(ctx *gin.Context, ids []int64) (map[int64]*model.User, error) {
	return dao.NewUserDao(ctx.Request.Context()).Db.GetUserBasicByIds(ctx.Request.Context(), uniqueInt64(ids))
}

func collectFriendIDs(friends []*model.Friend) []int64 {
	ids := make([]int64, 0, len(friends))
	for _, friend := range friends {
		ids = append(ids, friend.FriendID)
	}
	return ids
}

func collectRequestFromUserIDs(requests []*model.FriendRequest) []int64 {
	ids := make([]int64, 0, len(requests))
	for _, request := range requests {
		ids = append(ids, request.FromUser)
	}
	return ids
}

func uniqueInt64(ids []int64) []int64 {
	seen := make(map[int64]struct{}, len(ids))
	unique := make([]int64, 0, len(ids))
	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		unique = append(unique, id)
	}
	return unique
}
