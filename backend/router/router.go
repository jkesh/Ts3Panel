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
		auth.POST("/register", api.Register)
	}

	apiV1 := r.Group("/api/v1")
	apiV1.Use(middleware.AuthRequired())
	{
		apiV1.GET("/server", api.GetServerInfo)
		apiV1.GET("/clients", api.ListClients)
		apiV1.POST("/client/:id/kick", api.KickClient)
		apiV1.GET("/events/stream", api.StreamEvents) // SSE 需要处理 Token 传递问题，通常用 Query Param
	}

	return r
}
