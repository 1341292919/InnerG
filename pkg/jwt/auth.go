package jwt

import (
	"InnerG/dao"
	"InnerG/pack"
	"InnerG/pkg/constants"
	"InnerG/pkg/ctl"
	"InnerG/pkg/errno"
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

		access, refresh, err := CreateAllToken(userId)
		if err != nil {
			pack.RespError(ctx, err)
			ctx.Abort()
			return
		}
		ctx.Header(constants.AccessTokenHeader, access)
		ctx.Header(constants.RefreshTokenHeader, refresh)
		ctx.Request = ctx.Request.WithContext(ctl.NewContext(ctx.Request.Context(), &ctl.UserInfo{Id: userId, Token: token}))
		ctx.Next()
	}
}
