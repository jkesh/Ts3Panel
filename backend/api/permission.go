package api

import (
	"Ts3Panel/core"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jkesh/ts3-go/ts3"
)

// --- 请求结构体 ---

type CreateChannelReq struct {
	Name     string `json:"channel_name" binding:"required"`
	Password string `json:"channel_password"`
	Topic    string `json:"channel_topic"`
}

type CreateTokenReq struct {
	Type        int    `json:"type"` // 0=ServerGroup, 1=ChannelGroup
	GroupID     int    `json:"groupId" binding:"required"`
	ChannelID   int    `json:"channelId"`
	Description string `json:"description"`
}

// 定义权限列表响应结构
type ServerGroupPerm struct {
	PermID  int    `ts3:"permid" json:"permid"` // [新增] 数字ID
	Name    string `ts3:"permsid" json:"name"`  // 字符串ID
	Value   int    `ts3:"permvalue" json:"value"`
	Negated int    `ts3:"permnegated" json:"negated"`
	Skip    int    `ts3:"permskip" json:"skip"`
}

// ListServerGroupPerms 获取服务器组的当前权限列表
func ListServerGroupPerms(c *gin.Context) {
	sgid, _ := strconv.Atoi(c.Param("sgid"))

	core.Mutex.Lock()
	defer core.Mutex.Unlock()

	// 发送命令
	cmd := fmt.Sprintf("servergrouppermlist sgid=%d -names", sgid)
	resp, err := core.Client.Exec(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var perms []ServerGroupPerm
	if err := ts3.NewDecoder().Decode(resp, &perms); err != nil {
		// 忽略空结果错误
		c.JSON(200, gin.H{"data": []ServerGroupPerm{}})
		return
	}

	c.JSON(200, gin.H{"data": perms})
}

// --- 接口实现 ---

// CreateChannel 创建频道
func CreateChannel(c *gin.Context) {
	var req CreateChannelReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 构造 TS3 命令: channelcreate channel_name=...
	cmd := fmt.Sprintf("channelcreate channel_name=%s", ts3.Escape(req.Name))
	if req.Password != "" {
		cmd += fmt.Sprintf(" channel_password=%s", ts3.Escape(req.Password))
	}
	if req.Topic != "" {
		cmd += fmt.Sprintf(" channel_topic=%s", ts3.Escape(req.Topic))
	}
	// 设置为永久频道 (channel_flag_permanent=1)
	cmd += " channel_flag_permanent=1"

	core.Mutex.Lock()
	defer core.Mutex.Unlock()

	if _, err := core.Client.Exec(c.Request.Context(), cmd); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "Channel created successfully"})
}

// CreateToken 生成权限密钥
func CreateToken(c *gin.Context) {
	var req CreateTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	core.Mutex.Lock()
	defer core.Mutex.Unlock()

	// 调用 ts3-go 库现有的 TokenAdd 方法
	token, err := core.Client.TokenAdd(c.Request.Context(), req.Type, req.GroupID, req.ChannelID, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
