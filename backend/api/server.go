package api

import (
	"Ts3Panel/core"
	"github.com/gin-gonic/gin"
)

func GetServerInfo(c *gin.Context) {
	// [关键] 传入 c.Request.Context() 以便 Client.Exec 使用
	info, err := core.Client.ServerInfo(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, info)
}
