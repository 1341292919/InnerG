package v1

import (
	"InnerG/pack"
	"InnerG/pkg/errno"
	"InnerG/service"
	"InnerG/types"
	"github.com/gin-gonic/gin"
)

func NewChatSession() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.NewChatSessionReq
		if err := ctx.ShouldBind(&req); err != nil {
			pack.RespError(ctx, errno.ParamMissing.WithMessage(err.Error()))
			return
		}
		l := service.GetContactSrv()
		id, err := l.NewChatSession(ctx.Request.Context(), &req)
		if err != nil {
			pack.RespError(ctx, err)
			return
		}
		pack.RespData(ctx, types.NewChatSessionResp{SessionId: id})
	}
}

func StreamChat() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Content-Type", "text/event-stream")
		ctx.Header("Cache-Control", "no-cache")
		ctx.Header("Connection", "keep-alive")
		var req types.StreamChatReq
		if err := ctx.ShouldBind(&req); err != nil {
			pack.RespError(ctx, errno.ParamMissing.WithMessage(err.Error()))
			return
		}
		l := service.GetContactSrv()
		err := l.StreamChat(ctx, &req)
		if err != nil {
			pack.RespError(ctx, err)
			return
		}
		pack.RespSuccess(ctx)
	}
}

func GetUserSession() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.GetUserSessionListReq
		if err := ctx.ShouldBind(&req); err != nil {
			pack.RespError(ctx, errno.ParamMissing.WithMessage(err.Error()))
			return
		}
		l := service.GetContactSrv()
		list, total, err := l.GetUserSessionHistory(ctx.Request.Context(), &req)
		if err != nil {
			pack.RespError(ctx, err)
			return
		}
		resp := types.GetUserSessionListResp{
			SessionList: list,
			Total:       total,
		}
		pack.RespData(ctx, resp)
	}
}

func GetUserSessionDetail() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.GetUserSessionDetailReq
		if err := ctx.ShouldBind(&req); err != nil {
			pack.RespError(ctx, errno.ParamMissing.WithMessage(err.Error()))
			return
		}
		l := service.GetContactSrv()
		data, err := l.GetUserSessionDetail(ctx.Request.Context(), &req)
		if err != nil {
			pack.RespError(ctx, err)
			return
		}
		resp := types.GetUserSessionDetailResp{SessionDetail: data}
		pack.RespData(ctx, resp)
	}
}

func DeleteUserSession() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.DeleteUserSessionReq
		if err := ctx.ShouldBind(&req); err != nil {
			pack.RespError(ctx, errno.ParamMissing.WithMessage(err.Error()))
			return
		}
		l := service.GetContactSrv()
		err := l.DeleteUserSession(ctx.Request.Context(), req.SessionId)
		if err != nil {
			pack.RespError(ctx, err)
			return
		}
		pack.RespSuccess(ctx)
	}
}
