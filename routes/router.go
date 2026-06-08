package routes

import (
	api "InnerG/api/v1"
	"InnerG/middleware/confine"
	"InnerG/middleware/guard"
	"InnerG/middleware/metrics"
	"InnerG/pkg/constants"
	"InnerG/pkg/jwt"
	"InnerG/pkg/logger"
	"InnerG/types"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter() *gin.Engine {
	r := gin.New()

	// 设置 Gin 的日志输出到自定义 Writer
	gin.DefaultWriter = logger.GinWriter{}
	gin.DefaultErrorWriter = logger.GinWriter{}

	// 使用自定义的恢复中间件（可选，也可以使用 gin.Recovery()）
	r.Use(gin.Recovery())
	r.Use(logger.GinLoggerMiddleware())
	r.Use(metrics.QPSMiddleware())
	r.Use(guard.SentinelMiddleware())
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	v1 := r.Group("api/v1")
	{
		emailCodeKey := confine.ByBind[types.UserGetEmailCodeReq](func(req *types.UserGetEmailCodeReq) string { return req.Email })
		registerEmailKey := confine.ByBind[types.UserVerifyEmailAndRegisterReq](func(req *types.UserVerifyEmailAndRegisterReq) string { return req.Email })
		loginAccountKey := confine.ByBind[types.UserLoginReq](func(req *types.UserLoginReq) string { return req.Account })
		emailLoginKey := confine.ByBind[types.UserVerifyEmailAndLoginReq](func(req *types.UserVerifyEmailAndLoginReq) string { return req.Email })

		v1.POST("user/email/code", confine.Limit(
			confine.Rule{Name: "email_code_ip", Window: time.Minute, Max: 5, Key: confine.ByIP()},
			confine.Rule{Name: "email_code_email", Window: time.Hour, Max: 5, Key: emailCodeKey},
			confine.Rule{Name: "email_code_ip_email", Window: time.Minute, Max: 2, Key: confine.Compose(confine.ByIP(), emailCodeKey)},
		), api.UserGetEmailCodeHandler())
		v1.POST("user/register", confine.Limit(
			confine.Rule{Name: "register_ip", Window: time.Minute, Max: 20, Key: confine.ByIP()},
			confine.Rule{Name: "register_email", Window: 5 * time.Minute, Max: 5, Key: registerEmailKey},
			confine.Rule{Name: "register_ip_email", Window: 5 * time.Minute, Max: 5, Key: confine.Compose(confine.ByIP(), registerEmailKey)},
		), api.UserVerifyEmailAndRegister())
		v1.POST("user/login", confine.Limit(
			confine.Rule{Name: "login_ip", Window: time.Minute, Max: 10, Key: confine.ByIP()},
			confine.Rule{Name: "login_account", Window: 5 * time.Minute, Max: 5, Key: loginAccountKey},
			confine.Rule{Name: "login_ip_account", Window: 5 * time.Minute, Max: 5, Key: confine.Compose(confine.ByIP(), loginAccountKey)},
		), api.UserLogin())
		v1.POST("user/email/login", confine.Limit(
			confine.Rule{Name: "email_login_ip", Window: time.Minute, Max: 10, Key: confine.ByIP()},
			confine.Rule{Name: "email_login_email", Window: 5 * time.Minute, Max: 5, Key: emailLoginKey},
			confine.Rule{Name: "email_login_ip_email", Window: 5 * time.Minute, Max: 5, Key: confine.Compose(confine.ByIP(), emailLoginKey)},
		), api.UserVerifyEmailAndLogin())
		v1.GET("user/refresh-token", api.UserRefreshToken())
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

			// 社交聊天
			authed.POST("friend/request", api.SendFriendRequest())
			authed.POST("friend/request/handle", api.HandleFriendRequest())
			authed.POST("friend/delete", api.DeleteFriend())
			authed.GET("friend/list", api.ListFriends())
			authed.GET("friend/requests", api.ListFriendRequests())

			authed.POST("ws/ticket", api.CreateWebSocketTicket())
			authed.GET("ws/messages", api.GetWebSocketMessages())
			authed.GET("ws/unread", api.GetWebSocketUnread())
			authed.POST("ws/messages/ack", api.AckWebSocketMessages())
			authed.POST("ws/upload/image", api.UploadWebsocketImage())
			authed.POST("ws/upload/video", api.UploadWebsocketVideo())
		}

		ws := v1.Group("ws")
		ws.Use(jwt.WSTicketAuth())
		{
			ws.GET("connection", confine.Limit(
				confine.Rule{Name: "ws_connect_user", Window: constants.WebsocketConnectWindow, Max: constants.WebsocketConnectUserLimit, Key: confine.ByUserID()},
				confine.Rule{Name: "ws_connect_ip", Window: constants.WebsocketConnectWindow, Max: constants.WebsocketConnectIPLimit, Key: confine.ByIP()},
			), api.WebSocketConnectionHandler())
		}
	}
	return r
}
