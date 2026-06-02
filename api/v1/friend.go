package v1

import (
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
		if err := service.GetFriendSrv().CreateFriendRequest(ctx.Request.Context(), userID, req.FriendID); err != nil {
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
		userID := ctl.GetUserInfo(ctx.Request.Context()).Id
		friends, err := service.GetFriendSrv().ListFriends(ctx.Request.Context(), userID)
		if err != nil {
			pack.RespError(ctx, err)
			return
		}

		pack.RespData(ctx, types.FriendListResp{
			Friends: pack.BuildFriendList(friends),
			Total:   len(friends),
		})
	}
}

func ListFriendRequests() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctl.GetUserInfo(ctx.Request.Context()).Id
		requests, err := service.GetFriendSrv().ListInboundRequests(ctx.Request.Context(), userID)
		if err != nil {
			pack.RespError(ctx, err)
			return
		}

		pack.RespData(ctx, types.FriendRequestListResp{
			Requests: pack.BuildFriendRequestList(requests),
			Total:    len(requests),
		})
	}
}
