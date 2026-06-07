package main

import (
	"InnerG/config"
	"InnerG/dao"
	"InnerG/middleware/guard"
	"InnerG/middleware/metrics"
	"InnerG/pkg/logger"
	"InnerG/routes"
	"InnerG/service/websocket"
	"fmt"
)

func loading() {
	config.Init()
	dao.Init()
	metrics.Init()
	guard.InitSentinel()
	logger.InitLogger(config.Log.LogPath, config.Log.LogPrefix, config.Log.LogMaxDays)
	logger.InitGinLogger(config.Log.LogPath, config.Log.GinLogPrefix, config.Log.LogMaxDays)
}

func main() {
	loading()
	defer logger.CloseAll()
	r := routes.NewRouter()
	metrics.StartMetricsUpdater()
	// 消息队列启动
	websocket.NewWebsocketConsume().Run()
	_ = r.Run(config.Service.Address)
	fmt.Println("启动配成功...")
}
