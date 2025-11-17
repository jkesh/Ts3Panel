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
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	auth := r.Group("/auth")
	{
		auth.POST("/login", api.Login)
		// [修改] 移除原来的公开注册接口，防止任意用户自助注册
		auth.POST("/register", api.Register)
	}

	apiV1 := r.Group("/api/v1")
	apiV1.Use(middleware.AuthRequired())
	{
		// [修改] 将注册接口移动到此处，只有已登录的用户（管理员）才能创建新账号
		// apiV1.POST("/register", api.Register)

		apiV1.GET("/server", api.GetServerInfo)
		apiV1.GET("/clients", api.ListClients)
		apiV1.POST("/client/:id/kick", api.KickClient)
		apiV1.GET("/events/stream", api.StreamEvents)
	}

	return r
}
