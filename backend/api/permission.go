package api

import (
	"Ts3Panel/core"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jkesh/ts3-go/ts3"
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
	PermID  int    `ts3:"permid" json:"permid"`
	Name    string `ts3:"permsid" json:"name"`
	Value   int    `ts3:"permvalue" json:"value"`
	Negated int    `ts3:"permnegated" json:"negated"`
	Skip    int    `ts3:"permskip" json:"skip"`
}

// ListServerGroupPerms 获取服务器组的当前权限列表
func ListServerGroupPerms(c *gin.Context) {
	sgid, ok := parsePathInt(c, "sgid")
	if !ok {
		return
	}

	cmd := fmt.Sprintf("servergrouppermlist sgid=%d -names", sgid)
	resp, err := core.WithTS3Value(func(ts3Client *ts3.Client) (string, error) {
		return ts3Client.Exec(c.Request.Context(), cmd)
	})
	if err != nil {
		jsonError(c, http.StatusInternalServerError, err.Error())
		return
	}
	if strings.TrimSpace(resp) == "" {
		c.JSON(http.StatusOK, gin.H{"data": []ServerGroupPerm{}})
		return
	}

	var perms []ServerGroupPerm
	if err := ts3.NewDecoder().Decode(resp, &perms); err != nil {
		c.JSON(http.StatusOK, gin.H{"data": []ServerGroupPerm{}})
		return
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

	var b strings.Builder
	b.WriteString("channelcreate")
	b.WriteString(" channel_name=")
	b.WriteString(ts3.Escape(req.Name))
	if req.Password != "" {
		b.WriteString(" channel_password=")
		b.WriteString(ts3.Escape(req.Password))
	}
	if req.Topic != "" {
		b.WriteString(" channel_topic=")
		b.WriteString(ts3.Escape(req.Topic))
	}
	b.WriteString(" channel_flag_permanent=1")

	err := core.WithTS3(func(ts3Client *ts3.Client) error {
		_, err := ts3Client.Exec(c.Request.Context(), b.String())
		return err
	})
	if err != nil {
		jsonError(c, http.StatusInternalServerError, err.Error())
		return
	}

	jsonMessage(c, http.StatusOK, "Channel created successfully")
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
