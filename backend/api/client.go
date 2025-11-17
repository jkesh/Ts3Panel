package api

import (
	"Ts3Panel/core"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ListClients(c *gin.Context) {
	core.Mutex.Lock()
	defer core.Mutex.Unlock()
	clients, err := core.Client.ClientList(c.Request.Context(), "-uid", "-away", "-groups")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, clients)
}

func KickClient(c *gin.Context) {
	// 1. 权限检查
	role := c.GetString("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied: Admins only"})
		return
	}

	idStr := c.Param("id")

	uID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid Client ID format"})
		return
	}

	// 转换为 int 传递给 TS3 库 (如果是 2^64-1，转为 int 会变成 -1，这是符合底层逻辑的)
	id := int(uID)

	core.Mutex.Lock()
	defer core.Mutex.Unlock()
	// 2. 执行踢出
	err = core.Client.KickFromServer(c.Request.Context(), id, "Kicked by WebAdmin")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"msg": "Kicked"})
}
