package v1

import (
	"InnerG/pack"
	"InnerG/pkg/ctl"
	"InnerG/pkg/errno"
	service "InnerG/service/websocket"
	"InnerG/types"

	"github.com/gin-gonic/gin"
)

// SyncMessages 游标同步接口（合并私聊+群聊的增量消息）
func SyncMessages() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.SyncMessagesReq
		if err := ctx.ShouldBindQuery(&req); err != nil {
			pack.RespError(ctx, errno.ParamMissing.WithMessage(err.Error()))
			return
		}

		if req.Limit <= 0 {
			req.Limit = 100
		}
		if req.Limit > 200 {
			req.Limit = 200
		}

		uid := ctl.GetUserInfo(ctx.Request.Context()).Id
		srv := service.NewWebSocketSrv()

		resp, err := srv.SyncMessagesHTTP(ctx.Request.Context(), uid, req.LastID, req.Limit)
		if err != nil {
			pack.RespError(ctx, err)
			return
		}

		pack.RespData(ctx, types.SyncMessagesResp{
			Messages: resp.Messages,
			HasMore:  resp.HasMore,
			NextID:   resp.NextID,
		})
	}
}
