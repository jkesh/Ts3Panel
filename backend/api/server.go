package api

import (
	"Ts3Panel/core"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jkesh/ts3-go/ts3"
)

type serverInfoResp struct {
	Name          string `ts3:"virtualserver_name"`
	Version       string `ts3:"virtualserver_version"`
	Platform      string `ts3:"virtualserver_platform"`
	ClientsOnline int    `ts3:"virtualserver_clientsonline"`
	MaxClients    int    `ts3:"virtualserver_maxclients"`
	Uptime        int64  `ts3:"virtualserver_uptime"`
	Status        string `ts3:"virtualserver_status"`
}

// GetServerInfo 获取服务器信息
func GetServerInfo(c *gin.Context) {
	var info serverInfoResp
	err := core.WithTS3(func(ts3Client *ts3.Client) error {
		resp, err := ts3Client.Exec(c.Request.Context(), "serverinfo")
		if err != nil {
			return err
		}
		if strings.TrimSpace(resp) == "" {
			return errors.New("empty response from server")
		}
		return ts3.NewDecoder().Decode(resp, &info)
	})
	if err != nil {
		jsonError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"name":           info.Name,
		"version":        info.Version,
		"platform":       info.Platform,
		"clients_online": info.ClientsOnline,
		"max_clients":    info.MaxClients,
		"uptime":         info.Uptime,
		"status":         info.Status,
	})
}

type BroadcastReq struct {
	Message string `json:"message" binding:"required"`
}

// Broadcast 发送全服公告
func Broadcast(c *gin.Context) {
	if !ensureAdmin(c) {
		return
	}

	var req BroadcastReq
	if err := c.ShouldBindJSON(&req); err != nil {
		jsonError(c, http.StatusBadRequest, err.Error())
		return
	}

	err := core.WithTS3(func(ts3Client *ts3.Client) error {
		return ts3Client.Broadcast(c.Request.Context(), req.Message)
	})
	if err != nil {
		jsonError(c, http.StatusInternalServerError, err.Error())
		return
	}

	jsonMessage(c, http.StatusOK, "Broadcast sent")
}
