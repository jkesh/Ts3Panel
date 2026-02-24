package api

import (
	"Ts3Panel/core"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jkesh/ts3-go/v2/ts3"
	ts3models "github.com/jkesh/ts3-go/v2/ts3/models"
)

type CreateChannelReq struct {
	Name     string `json:"channel_name" binding:"required"`
	Password string `json:"channel_password"`
	Topic    string `json:"channel_topic"`
}

type CreateTokenReq struct {
	Type        int    `json:"type"`
	GroupID     int    `json:"groupId" binding:"required"`
	ChannelID   int    `json:"channelId"`
	Description string `json:"description"`
}

type ServerGroupPerm struct {
	PermID  int    `json:"permid"`
	Name    string `json:"name"`
	Value   int    `json:"value"`
	Negated int    `json:"negated"`
	Skip    int    `json:"skip"`
}

// ListServerGroupPerms 获取服务器组的当前权限列表
func ListServerGroupPerms(c *gin.Context) {
	sgid, ok := parsePathInt(c, "sgid")
	if !ok {
		return
	}

	rawPerms, err := core.WithTS3Value(func(ts3Client *ts3.Client) ([]ts3models.PermissionEntry, error) {
		return ts3Client.ServerGroupPermList(c.Request.Context(), sgid, true)
	})
	if err != nil {
		var ts3Err *ts3.Error
		if errors.As(err, &ts3Err) && ts3Err.Is(ts3.ErrDatabaseEmptyResult) {
			c.JSON(http.StatusOK, gin.H{"data": []ServerGroupPerm{}})
			return
		}
		jsonError(c, http.StatusInternalServerError, err.Error())
		return
	}

	perms := make([]ServerGroupPerm, 0, len(rawPerms))
	for _, p := range rawPerms {
		perms = append(perms, ServerGroupPerm{
			PermID:  p.PermID,
			Name:    p.PermSID,
			Value:   p.PermValue,
			Negated: p.PermNegated,
			Skip:    p.PermSkip,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": perms})
}

// CreateChannel 创建频道
func CreateChannel(c *gin.Context) {
	var req CreateChannelReq
	if err := c.ShouldBindJSON(&req); err != nil {
		jsonError(c, http.StatusBadRequest, err.Error())
		return
	}

	cid, err := core.WithTS3Value(func(ts3Client *ts3.Client) (int, error) {
		return ts3Client.ChannelCreate(c.Request.Context(), ts3.ChannelCreateOptions{
			Name:        req.Name,
			Password:    req.Password,
			Topic:       req.Topic,
			IsPermanent: true,
		})
	})
	if err != nil {
		jsonError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "Channel created successfully",
		"cid": cid,
	})
}

// CreateToken 生成权限密钥
func CreateToken(c *gin.Context) {
	var req CreateTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		jsonError(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := core.WithTS3Value(func(ts3Client *ts3.Client) (string, error) {
		return ts3Client.TokenAdd(c.Request.Context(), req.Type, req.GroupID, req.ChannelID, req.Description)
	})
	if err != nil {
		jsonError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
