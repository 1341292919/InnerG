package v1

import (
	"InnerG/pack"
	"InnerG/pkg/errno"
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
			pack.RespError(ctx, errno.ParamMissing.WithMessage(err.Error()))
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
			pack.RespError(ctx, errno.ParamMissing.WithMessage(err.Error()))
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

func UserLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.UserLoginReq
		if err := ctx.ShouldBind(&req); err != nil {
			pack.RespError(ctx, errno.ParamMissing.WithMessage(err.Error()))
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
			pack.RespError(ctx, errno.ParamMissing.WithMessage(err.Error()))
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

func GetUserInfo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		l := service.GetUserSrv()
		u, err := l.GetUserInfo(ctx.Request.Context())
		if err != nil {
			pack.RespError(ctx, err)
			return
		}
		pack.RespData(ctx, pack.BuildUser(u))
	}
}

func UserUpdateAccount() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.UpdateUserAccountReq
		if err := ctx.ShouldBind(&req); err != nil {
			pack.RespError(ctx, errno.ParamMissing.WithMessage(err.Error()))
			return
		}
		l := service.GetUserSrv()
		err := l.UpdateUserAccount(ctx.Request.Context(), req.Account)
		if err != nil {
			pack.RespError(ctx, err)
			return
		}
		pack.RespSuccess(ctx)
	}
}

func UserUpdateUserName() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.UpdateUserNameReq
		if err := ctx.ShouldBind(&req); err != nil {
			pack.RespError(ctx, errno.ParamMissing.WithMessage(err.Error()))
			return
		}
		l := service.GetUserSrv()
		err := l.UpdateUserName(ctx.Request.Context(), req.UserName)
		if err != nil {
			pack.RespError(ctx, err)
			return
		}
		pack.RespSuccess(ctx)
	}
}

func UserUpdateGender() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.UpdateUserGenderReq
		if err := ctx.ShouldBind(&req); err != nil {
			pack.RespError(ctx, errno.ParamMissing.WithMessage(err.Error()))
			return
		}
		l := service.GetUserSrv()
		err := l.UpdateUserGender(ctx.Request.Context(), req.Gender)
		if err != nil {
			pack.RespError(ctx, err)
			return
		}
		pack.RespSuccess(ctx)
	}
}

func UserLogOut() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		l := service.GetUserSrv()
		err := l.LogOut(ctx.Request.Context())
		if err != nil {
			pack.RespError(ctx, err)
			return
		}
		pack.RespSuccess(ctx)
	}
}

func UserUploadAvatar() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, fileHeader, err := ctx.Request.FormFile("file")
		if err != nil {
			pack.RespError(ctx, errno.ParamMissing.WithMessage(err.Error()))
			return
		}
		l := service.GetUserSrv()
		url, err := l.UpdateUserAvatar(ctx.Request.Context(), fileHeader)
		if err != nil {
			pack.RespError(ctx, err)
			return
		}
		pack.RespData(ctx, types.UpdateUserAvatarResp{AvatarUrl: url})
	}
}
