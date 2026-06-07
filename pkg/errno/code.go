package errno

const (
	SuccessCode              = 10000
	AuthErrorCode            = 30001 // 鉴权错误
	AuthInvalidCode          = 30002 // 鉴权无效
	AuthAccessExpiredCode    = 30003 // 访问令牌过期
	AuthRefreshExpiredCode   = 30004 // 刷新令牌过期
	RateLimitErrorCode       = 40001 // 请求频率限制
	TrafficGuardErrorCode    = 40002 // 服务流量保护
	FileTooLargeCode         = 40003 // 上传文件过大
	InternalServiceErrorCode = 50000

	// dao
	MySQLDBErrorCode = 50001
	MongoDBErrorCode = 50002
	RedisDBErrorCode = 50003
)
