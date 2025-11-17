package main

import (
	"Ts3Panel/config"
	"Ts3Panel/core"
	"Ts3Panel/database"
	"Ts3Panel/router"
	"log"
)

func main() {
	// 1. 加载配置
	config.LoadConfig()

	// 2. 初始化数据库
	database.InitDB()

	// 3. 初始化 TS3 连接
	core.InitTS3()

	// 4. 启动 Web 服务
	r := router.Setup()
	port := config.GlobalConfig.App.Port
	log.Printf("Starting server on %s", port)
	r.Run(port)
}
