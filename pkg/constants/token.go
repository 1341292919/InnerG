package constants

import (
	"fmt"
	"time"
)

const (
	TypeAccessToken  = 0
	TypeRefreshToken = 1

	AccessTokenTTL  = time.Hour * 24 * 7  // Access Token 有效期7天
	RefreshTokenTTL = time.Hour * 24 * 30 // Refresh Token 有效期30天\

	AuthHeader         = "Authorization" // 获取 Token 时的请求头
	AccessTokenHeader  = "Access-Token"  // 响应时的访问令牌头
	RefreshTokenHeader = "Refresh-Token" // 响应时的刷新令牌头

	Issuer = "yang" // token 颁发者

	ContextIdKey = "user_info"
)

var PublicKey = fmt.Sprintf("%v\n%v\n%v", "-----BEGIN PUBLIC KEY-----",
	"MCowBQYDK2VwAyEAmoXLgJUzBOvMIpAPv4stYZVg3iWT1BLoKKgXyWZUjb4=",
	"-----END PUBLIC KEY-----")
