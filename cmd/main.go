package main

import (
	"InnerG/config"
	"InnerG/dao"
	"InnerG/pkg/logger"
	"InnerG/routes"
	"fmt"
)

func loading() {
	config.Init()
	dao.Init()
	logger.InitLogger(config.Log.LogPath, config.Log.LogPrefix)
	logger.InitGinLogger(config.Log.LogPath, config.Log.GinLogPrefix)
}

func main() {
	loading()
	defer logger.CloseAll()
	r := routes.NewRouter()
	_ = r.Run(config.Service.Address)
	fmt.Println("启动配成功...")
}
