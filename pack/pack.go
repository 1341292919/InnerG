package pack

import (
	"InnerG/pkg/constants"
	"InnerG/pkg/errno"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Base struct {
	Code string `json:"code"`
	Msg  string `json:"message"`
}

type RespWithData struct {
	Code string `json:"code"`
	Msg  string `json:"message"`
	Data any    `json:"data"`
}

func RespError(ctx *gin.Context, err error) {
	Errno := errno.ConvertErr(err)
	ctx.JSON(http.StatusInternalServerError, Base{
		Code: strconv.FormatInt(Errno.ErrorCode, 10),
		Msg:  Errno.ErrorMsg,
	})
}

func RespSuccess(ctx *gin.Context) {
	Errno := errno.Success
	ctx.JSON(http.StatusOK, Base{
		Code: strconv.FormatInt(Errno.ErrorCode, 10),
		Msg:  Errno.ErrorMsg,
	})
}

func RespData(ctx *gin.Context, data any) {
	Errno := errno.Success
	ctx.JSON(http.StatusOK, RespWithData{
		Code: strconv.FormatInt(Errno.ErrorCode, 10),
		Msg:  Errno.ErrorMsg,
		Data: data,
	})
}

func WithToken(ctx *gin.Context, access, refresh string) {
	ctx.Header(constants.AccessTokenHeader, access)
	ctx.Header(constants.RefreshTokenHeader, refresh)
}
