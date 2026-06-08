package jwt

import (
	"InnerG/dao"
	"InnerG/pack"
	"InnerG/pkg/constants"
	"InnerG/pkg/ctl"
	"InnerG/pkg/errno"
	service "InnerG/service/websocket"
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := string(ctx.GetHeader(constants.AuthHeader))
		_, userId, err := CheckToken(token)
		if err != nil {
			pack.RespError(ctx, err)
			ctx.Abort()
			return
		}
		userDao := dao.NewUserDao(context.Background())
		if userDao.Cache.IsKeyExist(ctx.Request.Context(), fmt.Sprintf("token:%s", token)) {
			pack.RespError(ctx, errno.AuthInvalid.WithMessage("token have been confined"))
			ctx.Abort()
			return
		}

		ctx.Request = ctx.Request.WithContext(ctl.NewContext(ctx.Request.Context(), &ctl.UserInfo{Id: userId, Token: token}))
		ctx.Next()
	}
}

func WSTicketAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ticket := ctx.Query("ticket")
		if ticket == "" {
			pack.RespError(ctx, errno.AuthMissing)
			ctx.Abort()
			return
		}

		userID, err := service.NewWebSocketSrv().ConsumeTicket(ctx.Request.Context(), ticket)
		if err != nil {
			pack.RespError(ctx, err)
			ctx.Abort()
			return
		}

		ctx.Request = ctx.Request.WithContext(ctl.NewContext(ctx.Request.Context(), &ctl.UserInfo{Id: userID}))
		ctx.Next()
	}
}
