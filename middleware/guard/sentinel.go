package guard

import (
	"InnerG/pack"
	"InnerG/pkg/errno"
	"InnerG/pkg/logger"
	"fmt"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/gin-gonic/gin"
)

const globalResource = "api"

func InitSentinel() {
	if err := sentinel.InitDefault(); err != nil {
		panic(err)
	}

	_, err := flow.LoadRules([]*flow.Rule{
		{
			Resource:               globalResource,
			Threshold:              1000,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject, // 超过阈值后直接拒绝，不排队等待
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "POST:/api/v1/user/email/code",
			Threshold:              100,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject, // 超过阈值后直接拒绝，不排队等待
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "POST:/api/v1/user/login",
			Threshold:              100,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject, // 超过阈值后直接拒绝，不排队等待
			StatIntervalInMs:       1000,
		},
	})
	if err != nil {
		panic(err)
	}
}

func SentinelMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/metrics" { // Prometheus 监控接口不受 Sentinel 限流
			c.Next()
			return
		}

		globalEntry, globalBlock := sentinel.Entry(globalResource)
		if globalBlock != nil {
			logSentinelBlock(c, globalResource)
			pack.RespError(c, errno.TrafficGuardExceeded)
			c.Abort()
			return
		}
		defer globalEntry.Exit()

		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		pathResource := fmt.Sprintf("%s:%s", c.Request.Method, path)
		pathEntry, pathBlock := sentinel.Entry(pathResource)
		if pathBlock != nil {
			logSentinelBlock(c, pathResource)
			pack.RespError(c, errno.TrafficGuardExceeded)
			c.Abort()
			return
		}
		defer pathEntry.Exit()

		c.Next()
	}
}

func logSentinelBlock(c *gin.Context, resource string) {
	if logger.Log == nil {
		return
	}
	logger.Log.Errorf("sentinel rate limit exceeded: resource=%s ip=%s path=%s", resource, c.ClientIP(), c.Request.URL.Path)
}
