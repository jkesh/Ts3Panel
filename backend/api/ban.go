package api

import (
	"Ts3Panel/core"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jkesh/ts3-go/ts3"
	ts3models "github.com/jkesh/ts3-go/ts3/models"
)

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

type BanReq struct {
	IP     string `json:"ip"`
	Name   string `json:"name"`
	UID    string `json:"uid"`
	Time   int64  `json:"time"` // 秒，0=永久
	Reason string `json:"reason"`
}

// GetBanList 获取封禁列表
func GetBanList(c *gin.Context) {
	rawBans, err := core.WithTS3Value(func(ts3Client *ts3.Client) ([]ts3models.BanEntry, error) {
		return ts3Client.BanList(c.Request.Context())
	})
	if err != nil {
		var ts3Err *ts3.Error
		if errors.As(err, &ts3Err) && ts3Err.Is(ts3.ErrDatabaseEmptyResult) {
			c.JSON(http.StatusOK, gin.H{"data": []BanEntry{}})
			return
		}
		jsonError(c, http.StatusInternalServerError, err.Error())
		return
	}

	bans := make([]BanEntry, 0, len(rawBans))
	for _, b := range rawBans {
		bans = append(bans, BanEntry{
			ID:       b.BanID,
			IP:       b.IP,
			Name:     b.Name,
			UID:      b.UID,
			Created:  b.Created,
			Duration: b.Duration,
			Invoker:  b.InvokerName,
			Reason:   b.Reason,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": bans})
}

// AddBan 添加封禁
func AddBan(c *gin.Context) {
	var req BanReq
	if err := c.ShouldBindJSON(&req); err != nil {
		jsonError(c, http.StatusBadRequest, err.Error())
		return
	}
	if req.IP == "" && req.Name == "" && req.UID == "" {
		jsonError(c, http.StatusBadRequest, "ip/name/uid at least one is required")
		return
	}

	var b strings.Builder
	b.WriteString("banadd")
	if req.IP != "" {
		b.WriteString(" ip=")
		b.WriteString(ts3.Escape(req.IP))
	}
	if req.Name != "" {
		b.WriteString(" name=")
		b.WriteString(ts3.Escape(req.Name))
	}
	if req.UID != "" {
		b.WriteString(" uid=")
		b.WriteString(ts3.Escape(req.UID))
	}
	b.WriteString(fmt.Sprintf(" time=%d banreason=%s", req.Time, ts3.Escape(req.Reason)))

	err := core.WithTS3(func(ts3Client *ts3.Client) error {
		_, err := ts3Client.Exec(c.Request.Context(), b.String())
		return err
	})
	if err != nil {
		jsonError(c, http.StatusInternalServerError, err.Error())
		return
	}

	jsonMessage(c, http.StatusOK, "Banned")
}

// DeleteBan 删除封禁
func DeleteBan(c *gin.Context) {
	banID, ok := parsePathInt(c, "banid")
	if !ok {
		return
	}

	err := core.WithTS3(func(ts3Client *ts3.Client) error {
		return ts3Client.BanDelete(c.Request.Context(), banID)
	})
	if err != nil {
		jsonError(c, http.StatusInternalServerError, err.Error())
		return
	}

	jsonMessage(c, http.StatusOK, "Unbanned")
}

// DeleteAllBans 清空封禁
func DeleteAllBans(c *gin.Context) {
	err := core.WithTS3(func(ts3Client *ts3.Client) error {
		return ts3Client.BanDeleteAll(c.Request.Context())
	})
	if err != nil {
		jsonError(c, http.StatusInternalServerError, err.Error())
		return
	}

	jsonMessage(c, http.StatusOK, "All bans deleted")
}
