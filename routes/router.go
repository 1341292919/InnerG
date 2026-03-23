package routes

import (
	api "InnerG/api/v1"
	"InnerG/pkg/jwt"
	"InnerG/pkg/logger"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.New()

	// 设置 Gin 的日志输出到自定义 Writer
	gin.DefaultWriter = logger.GinWriter{}
	gin.DefaultErrorWriter = logger.GinWriter{}

	// 使用自定义的恢复中间件（可选，也可以使用 gin.Recovery()）
	r.Use(gin.Recovery())
	r.Use(logger.GinLoggerMiddleware())

	v1 := r.Group("api/v1")
	{
		v1.POST("user/email/code", api.UserGetEmailCodeHandler())
		v1.POST("user/register", api.UserVerifyEmailAndRegister())
		v1.POST("user/login", api.UserLogin())
		v1.POST("user/email/login", api.UserVerifyEmailAndLogin())
		authed := v1.Group("/") // 需要登陆保护
		authed.Use(jwt.Auth())
		{
			// 用户部分
			authed.GET("user/info", api.GetUserInfo())
			authed.POST("user/update/account", api.UserUpdateAccount())
			authed.POST("user/update/username", api.UserUpdateUserName())
			authed.POST("user/update/gender", api.UserUpdateGender())
			authed.POST("user/logout", api.UserLogOut())
			authed.POST("user/avatar", api.UserUploadAvatar())

			// 咨询聊天部分
			authed.POST("contact/session/start", api.NewChatSession())
			authed.POST("contact/session/stream", api.StreamChat())
			authed.GET("contact/session/list", api.GetUserSession())
			authed.GET("contact/session/detail", api.GetUserSessionDetail())
			authed.POST("contact/session/delete", api.DeleteUserSession())

			// 音乐服务部分
			authed.GET("music/playlist/list", api.GetPlaylistList())
			authed.GET("music/playlist/detail", api.GetPlaylistDetail())
			authed.GET("music/song/list", api.GetSongDetailList())
			authed.GET("music/song/detail", api.GetSongDetail())
		}
	}
	return r
}
