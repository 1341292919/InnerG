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

func AcceptFriendRequest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.FriendTargetReq
		if err := ctx.ShouldBind(&req); err != nil {
			pack.RespError(ctx, errno.ParamMissing.WithMessage(err.Error()))
			return
		}

		userID := ctl.GetUserInfo(ctx.Request.Context()).Id
		if err := service.GetFriendSrv().AcceptFriendRequest(ctx.Request.Context(), userID, req.FriendID); err != nil {
			pack.RespError(ctx, err)
			return
		}

		pack.RespSuccess(ctx)
	}
}

func RejectFriendRequest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.FriendTargetReq
		if err := ctx.ShouldBind(&req); err != nil {
			pack.RespError(ctx, errno.ParamMissing.WithMessage(err.Error()))
			return
		}

		userID := ctl.GetUserInfo(ctx.Request.Context()).Id
		if err := service.GetFriendSrv().RejectFriendRequest(ctx.Request.Context(), userID, req.FriendID); err != nil {
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
		friends, err := service.GetFriendSrv().ListInboundRequests(ctx.Request.Context(), userID)
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
