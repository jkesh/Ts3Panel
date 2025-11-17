package api

import (
	"Ts3Panel/core"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ListClients(c *gin.Context) {
	// [关键] 传入 Context
	clients, err := core.Client.ClientList(c.Request.Context(), "-uid", "-away", "-groups")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, clients)
}

func KickClient(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	// [关键] 传入 Context
	err := core.Client.KickFromServer(c.Request.Context(), id, "Kicked by WebAdmin")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"msg": "Kicked"})
}
