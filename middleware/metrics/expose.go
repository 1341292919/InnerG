package metrics

import (
	"InnerG/dao/cache"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"time"
)

var (
	// 定义Prometheus指标
	qpsGauge prometheus.Gauge

	qpsMinuteGauge prometheus.Gauge
)

// 定时更新QPS指标到Prometheus
func StartMetricsUpdater() {
	go func() {
		for {
			for i := 0; i < 300; i++ {
				currentQPS := GetCurrentQPSData()
				qpsGauge.Set(float64(currentQPS))
				time.Sleep(1 * time.Second)
			}
			qpsMinuteGauge.Set(float64(GetCurrentMinuteQPS()))
		}
	}()

}

func Init() {
	qpsCounter = NewQpsCounter(cache.GetRedisClient())
	qpsGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "gin_qps_current",
		Help: "Current QPS (requests per second)",
	})

	qpsMinuteGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "gin_qps_sum_5min",
		Help: "ALL QPS over last 5 Minutes",
	})
}
