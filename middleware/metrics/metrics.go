package metrics

import (
	"github.com/gin-gonic/gin"
	"time"
)

var qpsCounter *QpsCounter

// QPSMiddleware 统计每个请求
func QPSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path != "/metrics" {
			// 请求进来时增加计数
			qpsCounter.Increment()
		}
		c.Next()
	}
}

// GetQPSData 给Prometheus用的获取方法
func GetCurrentQPSData() (currentQPS int64) {
	// 3秒缓冲，避免出现读到未完成的数据
	return qpsCounter.GetTargetSecondQPS(time.Now().Add(-3 * time.Second).Unix())
}
func GetCurrentMinuteQPS() (currentQPS int64) {
	// 统计近5分钟的QPS
	return qpsCounter.GetMinuteQps(time.Now().Add(-300 * time.Second).Unix())
}
