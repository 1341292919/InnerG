package main

import (
	"InnerG/config"
	"InnerG/dao"
	"InnerG/routes"
	"fmt"
)

func loading() {
	config.Init()
	dao.Init()
}

func main() {
	loading()
	r := routes.NewRouter()
	_ = r.Run(config.Service.Address)
	fmt.Println("启动配成功...")
}
