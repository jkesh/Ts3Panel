package router

import (
	"Ts3Panel/api"
	"Ts3Panel/middleware"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.Default()

	// CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	auth := r.Group("/auth")
	{
		auth.POST("/login", api.Login)
		auth.POST("/register", api.Register)
	}

	apiV1 := r.Group("/api/v1")
	apiV1.Use(middleware.AuthRequired())
	{
		apiV1.GET("/server", api.GetServerInfo)
		apiV1.GET("/clients", api.ListClients)
		apiV1.POST("/client/:id/kick", api.KickClient)
		apiV1.GET("/events/stream", api.StreamEvents)
		apiV1.POST("/broadcast", api.Broadcast)
		apiV1.POST("/channel/create", api.CreateChannel) // 新建频道
		apiV1.POST("/token/create", api.CreateToken)     // 生成 Token

		apiV1.POST("/channel/:cid/permission", api.AddChannelPerm)
		apiV1.POST("/clientdb/:cldbid/permission", api.AddClientDbPerm)

		apiV1.GET("/channels", api.GetChannels)          // 获取频道列表
		apiV1.DELETE("/channel/:cid", api.DeleteChannel) // 删除频道

		// 服务器组管理
		apiV1.GET("/servergroups", api.GetServerGroups)
		apiV1.DELETE("/servergroup/:sgid", api.DeleteServerGroup)

		// 权限管理
		apiV1.POST("/servergroup/:sgid/permission", api.AddServerGroupPerm)
		apiV1.GET("/servergroup/:sgid/permissions", api.ListServerGroupPerms)

		// 封禁处理
		apiV1.GET("/bans", api.GetBanList)
		apiV1.POST("/ban", api.AddBan)
		apiV1.DELETE("/ban/:banid", api.DeleteBan)
		apiV1.DELETE("/bans/all", api.DeleteAllBans)

		// 机器人管理路由
		apiV1.GET("/bots", api.ListBots)
		apiV1.POST("/bot", api.AddBot)
		apiV1.DELETE("/bot/:id", api.DeleteBot)
		apiV1.POST("/bot/:id/command", api.SendCommand) // 核心控制接口
	}

	return r
}
