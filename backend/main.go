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
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("config load failed: %v", err)
	}

	// 2. 初始化数据库
	if err := database.InitDB(); err != nil {
		log.Fatalf("database init failed: %v", err)
	}

	// 3. 初始化 TS3 连接
	if err := core.InitTS3(); err != nil {
		log.Fatalf("ts3 init failed: %v", err)
	}

	// 4. 启动 Web 服务
	r := router.Setup()
	port := config.GlobalConfig.App.Port
	log.Printf("Starting server on %s", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("server start failed: %v", err)
	}
}
