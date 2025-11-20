package api

import (
	"Ts3Panel/core"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PermReq struct {
	PermName  string `json:"perm_name" binding:"required"`
	PermValue int    `json:"perm_value"` // 如果是删除权限，这个值可能不重要，但在 Add 时必需
}

// AddChannelPerm 修改频道权限
func AddChannelPerm(c *gin.Context) {
	cid, _ := strconv.Atoi(c.Param("cid"))
	var req PermReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	core.Mutex.Lock()
	defer core.Mutex.Unlock()

	if err := core.Client.ChannelAddPerm(c.Request.Context(), cid, req.PermName, req.PermValue); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "Updated"})
}

// AddServerGroupPerm 修改服务器组权限
func AddServerGroupPerm(c *gin.Context) {
	sgid, _ := strconv.Atoi(c.Param("sgid"))
	var req PermReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	core.Mutex.Lock()
	defer core.Mutex.Unlock()

	// negated=false, skip=false 为默认
	if err := core.Client.ServerGroupAddPerm(c.Request.Context(), sgid, req.PermName, req.PermValue, false, false); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "Updated"})
}

// AddClientDbPerm 修改客户端(数据库ID)权限
func AddClientDbPerm(c *gin.Context) {
	cldbid, _ := strconv.Atoi(c.Param("cldbid"))
	var req PermReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	core.Mutex.Lock()
	defer core.Mutex.Unlock()

	// skip=false 为默认
	if err := core.Client.ClientAddPerm(c.Request.Context(), cldbid, req.PermName, req.PermValue, false); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "Updated"})
}
func GetChannels(c *gin.Context) {
	core.Mutex.Lock()
	defer core.Mutex.Unlock()

	// 调用 ts3-go 的 ChannelList 方法
	channels, err := core.Client.ChannelList(c.Request.Context(), "-topic", "-flags")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": channels})
}

// DeleteChannel 删除频道
func DeleteChannel(c *gin.Context) {
	cidStr := c.Param("cid")
	forceStr := c.DefaultQuery("force", "0") // 0=只删除空频道, 1=强制删除(包括里面的人)

	core.Mutex.Lock()
	defer core.Mutex.Unlock()

	// [修正] 将 HfStr 改为 forceStr
	cmd := fmt.Sprintf("channeldelete cid=%s force=%s", cidStr, forceStr)

	if _, err := core.Client.Exec(c.Request.Context(), cmd); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "Deleted"})
}

// GetServerGroups 获取服务器组列表
func GetServerGroups(c *gin.Context) {
	core.Mutex.Lock()
	defer core.Mutex.Unlock()

	// 调用 ts3-go 的 ServerGroupList 方法
	groups, err := core.Client.ServerGroupList(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": groups})
}

// DeleteServerGroup 删除服务器组
func DeleteServerGroup(c *gin.Context) {
	sgid := c.Param("sgid")
	force := c.DefaultQuery("force", "1") // 1=强制删除

	core.Mutex.Lock()
	defer core.Mutex.Unlock()

	// 构造命令: servergroupdel sgid=xxx force=1
	cmd := fmt.Sprintf("servergroupdel sgid=%s force=%s", sgid, force)
	if _, err := core.Client.Exec(c.Request.Context(), cmd); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "Server group deleted"})
}
