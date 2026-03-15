package jwt

import (
	"InnerG/pack"
	"InnerG/pkg/constants"
	"InnerG/pkg/ctl"
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
		access, refresh, err := CreateAllToken(userId)
		if err != nil {
			pack.RespError(ctx, err)
			ctx.Abort()
			return
		}

		ctx.Header(constants.AccessTokenHeader, access)
		ctx.Header(constants.RefreshTokenHeader, refresh)
		ctx.Request = ctx.Request.WithContext(ctl.NewContext(ctx.Request.Context(), &ctl.UserInfo{Id: userId}))
		ctx.Next()
	}
}
