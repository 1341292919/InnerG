package v1

import (
	"InnerG/pack"
	"InnerG/pkg/jwt"
	"InnerG/service"
	"InnerG/types"
	"github.com/gin-gonic/gin"
	"strconv"
)

func UserGetEmailCodeHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.UserGetEmailCodeReq
		if err := ctx.ShouldBind(&req); err != nil {
			// log
			pack.RespError(ctx, err)
			return
		}

		// 参数检验-邮箱地址检验

		l := service.GetUserSrv()
		err := l.GetEmailCode(ctx.Request.Context(), &req)
		if err != nil {
			pack.RespError(ctx, err)
			return
		}
		pack.RespSuccess(ctx)
	}
}

func UserVerifyEmailAndRegister() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.UserVerifyEmailAndRegisterReq
		if err := ctx.ShouldBind(&req); err != nil {
			// log
			pack.RespError(ctx, err)
			return
		}
		l := service.GetUserSrv()
		err := l.VerifyEmailAndRegister(ctx.Request.Context(), &req)
		if err != nil {
			pack.RespError(ctx, err)
			return
		}
		pack.RespSuccess(ctx)
	}
}

func UserLoginRegister() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.UserLoginReq
		if err := ctx.ShouldBind(&req); err != nil {
			// log
			pack.RespError(ctx, err)
			return
		}
		l := service.GetUserSrv()
		u, err := l.Login(ctx, &req)
		if err != nil {
			pack.RespError(ctx, err)
			return
		}
		access, refresh, err := jwt.CreateAllToken(strconv.FormatInt(int64(u.ID), 10))
		if err != nil {
			pack.RespError(ctx, err)
			return
		}
		pack.WithToken(ctx, access, refresh)
		pack.RespData(ctx, pack.BuildUser(u))
	}
}

func UserVerifyEmailAndLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.UserVerifyEmailAndLoginReq
		if err := ctx.ShouldBind(&req); err != nil {
			pack.RespError(ctx, err)
			return
		}
		l := service.GetUserSrv()
		u, err := l.VerifyEmailAndLogin(ctx, &req)
		if err != nil {
			pack.RespError(ctx, err)
			return
		}
		access, refresh, err := jwt.CreateAllToken(strconv.FormatInt(int64(u.ID), 10))
		if err != nil {
			pack.RespError(ctx, err)
			return
		}
		pack.WithToken(ctx, access, refresh)
		pack.RespData(ctx, pack.BuildUser(u))
	}
}
