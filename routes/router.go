package routes

import (
	api "InnerG/api/v1"
	"InnerG/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	v1 := r.Group("api/v1")
	{
		v1.POST("user/email/code", api.UserGetEmailCodeHandler())
		v1.POST("user/register", api.UserVerifyEmailAndRegister())
		v1.POST("user/login", api.UserLoginRegister())
		v1.POST("user/email/login", api.UserVerifyEmailAndLogin())
		authed := v1.Group("/") // 需要登陆保护
		authed.Use(jwt.Auth())
		{
			// 咨询聊天部分
			authed.GET("contact", api.NewContactHandler())
			authed.POST("contact/session/start", api.NewChatSession())
			authed.POST("contact/session/stream", api.StreamChat())
			authed.GET("contact/session/list", api.GetUserSession())
			authed.GET("contact/session/detail", api.GetUserSessionDetail())
		}
	}
	return r
}
