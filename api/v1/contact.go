package v1

import (
	"InnerG/pack"
	"InnerG/pkg/errno"
	"InnerG/service"
	"InnerG/types"
	"github.com/gin-gonic/gin"
)

type reqBody struct {
	Model    string    `json:"model"`    // 添加标签，指定JSON中的键名为小写
	Messages []message `json:"messages"` // 添加标签
	Stream   bool      `json:"stream"`   // 添加标签
}

type message struct {
	Role    string `json:"role"`    // 添加标签
	Content string `json:"content"` // 添加标签
}

func NewContactHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Content-Type", "text/event-stream")
		ctx.Header("Cache-Control", "no-cache")
		ctx.Header("Connection", "keep-alive")

		l := service.GetContactSrv()
		err := l.NewContact(ctx)
		if err != nil {
			pack.RespError(ctx, err)
			return
		}
		pack.RespSuccess(ctx)
	}
}

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

func NewStreamChat() gin.HandlerFunc {
	return func(ctx *gin.Context) {
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
		list, total, err := l.GetUserSessionHistory(ctx.Request.Context())
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
