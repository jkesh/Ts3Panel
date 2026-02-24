package api

import (
	"Ts3Panel/core"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jkesh/ts3-go/v2/ts3"
	ts3models "github.com/jkesh/ts3-go/v2/ts3/models"
)

type PermReq struct {
	PermName  string `json:"perm_name" binding:"required"`
	PermValue int    `json:"perm_value"`
}

// AddChannelPerm 修改频道权限
func AddChannelPerm(c *gin.Context) {
	cid, ok := parsePathInt(c, "cid")
	if !ok {
		return
	}

	var req PermReq
	if err := c.ShouldBindJSON(&req); err != nil {
		jsonError(c, http.StatusBadRequest, err.Error())
		return
	}

	err := core.WithTS3(func(ts3Client *ts3.Client) error {
		return ts3Client.ChannelAddPerm(c.Request.Context(), cid, req.PermName, req.PermValue)
	})
	if err != nil {
		jsonError(c, http.StatusInternalServerError, err.Error())
		return
	}

	jsonMessage(c, http.StatusOK, "Updated")
}

// AddServerGroupPerm 修改服务器组权限
func AddServerGroupPerm(c *gin.Context) {
	sgid, ok := parsePathInt(c, "sgid")
	if !ok {
		return
	}

	var req PermReq
	if err := c.ShouldBindJSON(&req); err != nil {
		jsonError(c, http.StatusBadRequest, err.Error())
		return
	}

	err := core.WithTS3(func(ts3Client *ts3.Client) error {
		return ts3Client.ServerGroupAddPerm(c.Request.Context(), sgid, req.PermName, req.PermValue, false, false)
	})
	if err != nil {
		jsonError(c, http.StatusInternalServerError, err.Error())
		return
	}

	jsonMessage(c, http.StatusOK, "Updated")
}

// AddClientDbPerm 修改客户端(数据库ID)权限
func AddClientDbPerm(c *gin.Context) {
	cldbid, ok := parsePathInt(c, "cldbid")
	if !ok {
		return
	}

	var req PermReq
	if err := c.ShouldBindJSON(&req); err != nil {
		jsonError(c, http.StatusBadRequest, err.Error())
		return
	}

	err := core.WithTS3(func(ts3Client *ts3.Client) error {
		return ts3Client.ClientAddPerm(c.Request.Context(), cldbid, req.PermName, req.PermValue, false)
	})
	if err != nil {
		jsonError(c, http.StatusInternalServerError, err.Error())
		return
	}

	jsonMessage(c, http.StatusOK, "Updated")
}

func GetChannels(c *gin.Context) {
	channels, err := core.WithTS3Value(func(ts3Client *ts3.Client) ([]ts3models.Channel, error) {
		return ts3Client.ChannelList(c.Request.Context(), "-topic", "-flags")
	})
	if err != nil {
		jsonError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": channels})
}

// DeleteChannel 删除频道
func DeleteChannel(c *gin.Context) {
	cid, ok := parsePathInt(c, "cid")
	if !ok {
		return
	}
	force, ok := parseQueryIntWithDefault(c, "force", 0)
	if !ok {
		return
	}

	err := core.WithTS3(func(ts3Client *ts3.Client) error {
		return ts3Client.ChannelDelete(c.Request.Context(), cid, force == 1)
	})
	if err != nil {
		jsonError(c, http.StatusInternalServerError, err.Error())
		return
	}

	jsonMessage(c, http.StatusOK, "Deleted")
}

// GetServerGroups 获取服务器组列表
func GetServerGroups(c *gin.Context) {
	groups, err := core.WithTS3Value(func(ts3Client *ts3.Client) ([]ts3models.ServerGroup, error) {
		return ts3Client.ServerGroupList(c.Request.Context())
	})
	if err != nil {
		jsonError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": groups})
}

// DeleteServerGroup 删除服务器组
func DeleteServerGroup(c *gin.Context) {
	sgid, ok := parsePathInt(c, "sgid")
	if !ok {
		return
	}
	force, ok := parseQueryIntWithDefault(c, "force", 1)
	if !ok {
		return
	}

	err := core.WithTS3(func(ts3Client *ts3.Client) error {
		return ts3Client.ServerGroupDelete(c.Request.Context(), sgid, force == 1)
	})
	if err != nil {
		jsonError(c, http.StatusInternalServerError, err.Error())
		return
	}

	jsonMessage(c, http.StatusOK, "Server group deleted")
}
