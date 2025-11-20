package api

import (
	"Ts3Panel/core"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jkesh/ts3-go/ts3"
)

// BanList 结构体
type BanEntry struct {
	ID       int    `ts3:"banid" json:"id"`
	IP       string `ts3:"ip" json:"ip"`
	Name     string `ts3:"name" json:"name"`
	UID      string `ts3:"uid" json:"uid"`
	Created  int64  `ts3:"created" json:"created"`
	Duration int64  `ts3:"duration" json:"duration"`
	Invoker  string `ts3:"invokername" json:"invoker"`
	Reason   string `ts3:"reason" json:"reason"`
}

// GetBanList 获取封禁列表
func GetBanList(c *gin.Context) {
	core.Mutex.Lock()
	defer core.Mutex.Unlock()

	resp, err := core.Client.Exec(c.Request.Context(), "banlist")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var bans []BanEntry
	if err := ts3.NewDecoder().Decode(resp, &bans); err != nil {
		c.JSON(200, gin.H{"data": []BanEntry{}})
		return
	}
	c.JSON(200, gin.H{"data": bans})
}

// AddBan 添加封禁
type BanReq struct {
	IP     string `json:"ip"`
	Name   string `json:"name"`
	UID    string `json:"uid"`
	Time   int64  `json:"time"` // 秒，0=永久
	Reason string `json:"reason"`
}

func AddBan(c *gin.Context) {
	var req BanReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	cmd := "banadd"
	if req.IP != "" {
		cmd += fmt.Sprintf(" ip=%s", req.IP)
	}
	if req.Name != "" {
		cmd += fmt.Sprintf(" name=%s", ts3.Escape(req.Name))
	}
	if req.UID != "" {
		cmd += fmt.Sprintf(" uid=%s", req.UID)
	}
	cmd += fmt.Sprintf(" time=%d banreason=%s", req.Time, ts3.Escape(req.Reason))

	core.Mutex.Lock()
	defer core.Mutex.Unlock()

	if _, err := core.Client.Exec(c.Request.Context(), cmd); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"msg": "Banned"})
}

// DeleteBan 删除封禁
func DeleteBan(c *gin.Context) {
	banID := c.Param("banid")

	core.Mutex.Lock()
	defer core.Mutex.Unlock()

	cmd := fmt.Sprintf("bandel banid=%s", banID)
	if _, err := core.Client.Exec(c.Request.Context(), cmd); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"msg": "Unbanned"})
}

// DeleteAllBans 清空封禁
func DeleteAllBans(c *gin.Context) {
	core.Mutex.Lock()
	defer core.Mutex.Unlock()

	if _, err := core.Client.Exec(c.Request.Context(), "bandelall"); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"msg": "All bans deleted"})
}
