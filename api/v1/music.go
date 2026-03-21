package v1

import (
	"InnerG/pack"
	"InnerG/pkg/errno"
	"InnerG/service"
	"InnerG/types"
	"github.com/gin-gonic/gin"
)

func GetPlaylistList() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.GetPlaylistListReq
		if err := ctx.ShouldBind(&req); err != nil {
			pack.RespError(ctx, errno.ParamMissing.WithMessage(err.Error()))
			return
		}
		l := service.GetMusicSrv()
		list, total, err := l.GetPlaylistList(ctx.Request.Context(), &req)
		if err != nil {
			pack.RespError(ctx, err)
			return
		}
		resp := types.GetPlaylistListResp{
			Total:        total,
			PlaylistList: list,
		}
		pack.RespData(ctx, resp)
	}
}

func GetPlaylistDetail() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.GetPlaylistDetailReq
		if err := ctx.ShouldBind(&req); err != nil {
			pack.RespError(ctx, errno.ParamMissing.WithMessage(err.Error()))
			return
		}
		l := service.GetMusicSrv()
		data, err := l.GetPlaylistDetail(ctx.Request.Context(), &req)
		if err != nil {
			pack.RespError(ctx, err)
			return
		}
		pack.RespData(ctx, types.GetPlaylistDetailResp{PlaylistDetail: data})
	}
}

func GetSongDetail() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.GetSongDetailReq
		if err := ctx.ShouldBind(&req); err != nil {
			pack.RespError(ctx, errno.ParamMissing.WithMessage(err.Error()))
			return
		}
		l := service.GetMusicSrv()
		data, err := l.GetSongDetail(ctx.Request.Context(), &req)
		if err != nil {
			pack.RespError(ctx, err)
			return
		}
		pack.RespData(ctx, types.GetSongDetailResp{SongDetail: data})
	}
}

func GetSongDetailList() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.GetSongDetailListReq
		if err := ctx.ShouldBind(&req); err != nil {
			pack.RespError(ctx, errno.ParamMissing.WithMessage(err.Error()))
			return
		}
		l := service.GetMusicSrv()
		list, total, err := l.GetSongDetailList(ctx.Request.Context(), &req)
		if err != nil {
			pack.RespError(ctx, err)
			return
		}
		pack.RespData(ctx, types.GetSongDetailListResp{Total: total, SongList: list})
	}
}
