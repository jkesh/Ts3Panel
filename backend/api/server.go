package api

import (
	"Ts3Panel/core"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetServerInfo 获取服务器信息
func GetServerInfo(c *gin.Context) {
	core.Mutex.Lock()
	defer core.Mutex.Unlock()

	// [关键] 使用 Exec 传入 Context 和命令字符串
	resp, err := core.Client.Exec(c.Request.Context(), "serverinfo")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 库返回的是 string，直接使用
	fullResp := resp
	if fullResp == "" {
		c.JSON(500, gin.H{"error": "Empty response from server"})
		return
	}

	lines := strings.Split(fullResp, "\n")
	rawLine := strings.TrimSpace(lines[0])

	data := make(map[string]interface{})
	parts := strings.Split(rawLine, " ")
	for _, part := range parts {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) == 2 {
			key := kv[0]
			val := unescape(kv[1])
			data[key] = val
		}
	}

	result := map[string]interface{}{
		"name":           data["virtualserver_name"],
		"version":        data["virtualserver_version"],
		"platform":       data["virtualserver_platform"],
		"clients_online": data["virtualserver_clientsonline"],
		"max_clients":    data["virtualserver_maxclients"],
		"uptime":         data["virtualserver_uptime"],
		"status":         data["virtualserver_status"],
	}

	c.JSON(200, result)
}

type BroadcastReq struct {
	Message string `json:"message" binding:"required"`
}

// Broadcast 发送全服公告
func Broadcast(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied: Admins only"})
		return
	}

	var req BroadcastReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	core.Mutex.Lock()
	defer core.Mutex.Unlock()

	escapedMsg := escape(req.Message)

	// [修复] 将命令和参数拼接成一个单独的字符串
	// 错误写法: Exec(ctx, "cmd", "arg1", "arg2")
	// 正确写法: Exec(ctx, "cmd arg1 arg2")
	cmd := "sendtextmessage targetmode=3 msg=" + escapedMsg

	_, err := core.Client.Exec(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"msg": "Broadcast sent"})
}

// TS3 协议辅助函数
func unescape(str string) string {
	str = strings.ReplaceAll(str, "\\s", " ")
	str = strings.ReplaceAll(str, "\\/", "/")
	str = strings.ReplaceAll(str, "\\p", "|")
	str = strings.ReplaceAll(str, "\\b", "\b")
	str = strings.ReplaceAll(str, "\\f", "\f")
	str = strings.ReplaceAll(str, "\\n", "\n")
	str = strings.ReplaceAll(str, "\\r", "\r")
	str = strings.ReplaceAll(str, "\\t", "\t")
	str = strings.ReplaceAll(str, "\\v", "\v")
	str = strings.ReplaceAll(str, "\\\\", "\\")
	return str
}

func escape(str string) string {
	str = strings.ReplaceAll(str, "\\", "\\\\")
	str = strings.ReplaceAll(str, "/", "\\/")
	str = strings.ReplaceAll(str, " ", "\\s")
	str = strings.ReplaceAll(str, "|", "\\p")
	str = strings.ReplaceAll(str, "\a", "\\a")
	str = strings.ReplaceAll(str, "\b", "\\b")
	str = strings.ReplaceAll(str, "\f", "\\f")
	str = strings.ReplaceAll(str, "\n", "\\n")
	str = strings.ReplaceAll(str, "\r", "\\r")
	str = strings.ReplaceAll(str, "\t", "\\t")
	str = strings.ReplaceAll(str, "\v", "\\v")
	return str
}
