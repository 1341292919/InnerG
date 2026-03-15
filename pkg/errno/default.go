package errno

var (
	ParamFileMiss        = NewErr(20000, "file miss")
	Success              = NewErr(10000, "success")
	InternalServiceError = NewErr(50000, "internal server error")
	AuthError            = NewErr(AuthErrorCode, "鉴权失败")            // 鉴权失败，通常是内部错误，如解析失败
	AuthInvalid          = NewErr(AuthInvalidCode, "鉴权无效")          // 鉴权无效，如令牌颁发者不是 west2-online
	AuthAccessExpired    = NewErr(AuthAccessExpiredCode, "访问令牌过期")  // 访问令牌过期
	AuthRefreshExpired   = NewErr(AuthRefreshExpiredCode, "刷新令牌过期") // 刷新令牌过期
	AuthMissing          = NewErr(AuthInvalidCode, "缺失合法鉴权数据")      // 鉴权缺失，如访问令牌缺失
)
