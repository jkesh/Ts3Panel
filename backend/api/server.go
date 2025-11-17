package api

import (
	"Ts3Panel/core"

	"github.com/gin-gonic/gin"
)

func GetServerInfo(c *gin.Context) {
	core.Mutex.Lock()
	defer core.Mutex.Unlock()

	info, err := core.Client.ServerInfo(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, info)
}
