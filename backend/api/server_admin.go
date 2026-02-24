package api

import (
	"Ts3Panel/core"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jkesh/ts3-go/v2/ts3"
)

type UpdateServerSettingsReq struct {
	Name            string `json:"name"`
	WelcomeMessage  string `json:"welcome_message"`
	Password        string `json:"password"`
	MaxClients      int    `json:"max_clients"`
	HostMessage     string `json:"host_message"`
	HostMessageMode int    `json:"host_message_mode"`
}

func UpdateServerSettings(c *gin.Context) {
	if !ensureAdmin(c) {
		return
	}

	var req UpdateServerSettingsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		jsonError(c, http.StatusBadRequest, err.Error())
		return
	}

	err := core.WithTS3(func(ts3Client *ts3.Client) error {
		return ts3Client.ServerEdit(c.Request.Context(), ts3.ServerEditOptions{
			Name:            strings.TrimSpace(req.Name),
			WelcomeMessage:  strings.TrimSpace(req.WelcomeMessage),
			Password:        req.Password,
			MaxClients:      req.MaxClients,
			HostMessage:     strings.TrimSpace(req.HostMessage),
			HostMessageMode: req.HostMessageMode,
		})
	})
	if err != nil {
		jsonError(c, http.StatusInternalServerError, err.Error())
		return
	}

	jsonMessage(c, http.StatusOK, "Server settings updated")
}

type TempPasswordReq struct {
	Password              string `json:"password" binding:"required"`
	Description           string `json:"description"`
	DurationSeconds       int64  `json:"duration_seconds" binding:"required"`
	TargetChannelID       int    `json:"target_channel_id"`
	TargetChannelPassword string `json:"target_channel_password"`
}

func ListTempPasswords(c *gin.Context) {
	passwords, err := core.WithTS3Value(func(ts3Client *ts3.Client) (any, error) {
		return ts3Client.ServerTempPasswordList(c.Request.Context())
	})
	if err != nil {
		jsonError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": passwords})
}

func CreateTempPassword(c *gin.Context) {
	if !ensureAdmin(c) {
		return
	}

	var req TempPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		jsonError(c, http.StatusBadRequest, err.Error())
		return
	}

	err := core.WithTS3(func(ts3Client *ts3.Client) error {
		return ts3Client.ServerTempPasswordAdd(c.Request.Context(), ts3.ServerTempPasswordOptions{
			Password:              req.Password,
			Description:           req.Description,
			DurationSeconds:       req.DurationSeconds,
			TargetChannelID:       req.TargetChannelID,
			TargetChannelPassword: req.TargetChannelPassword,
		})
	})
	if err != nil {
		jsonError(c, http.StatusInternalServerError, err.Error())
		return
	}

	jsonMessage(c, http.StatusOK, "Temporary password created")
}

func DeleteTempPassword(c *gin.Context) {
	if !ensureAdmin(c) {
		return
	}

	password := strings.TrimSpace(c.Param("password"))
	if password == "" {
		jsonError(c, http.StatusBadRequest, "Invalid password")
		return
	}

	err := core.WithTS3(func(ts3Client *ts3.Client) error {
		return ts3Client.ServerTempPasswordDelete(c.Request.Context(), password)
	})
	if err != nil {
		jsonError(c, http.StatusInternalServerError, err.Error())
		return
	}

	jsonMessage(c, http.StatusOK, "Temporary password deleted")
}

type QueryLoginReq struct {
	ClientDBID int `json:"client_db_id" binding:"required"`
	ServerID   int `json:"server_id"`
}

func ListQueryLogins(c *gin.Context) {
	list, err := core.WithTS3Value(func(ts3Client *ts3.Client) (any, error) {
		return ts3Client.QueryLoginList(c.Request.Context())
	})
	if err != nil {
		jsonError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": list})
}

func CreateQueryLogin(c *gin.Context) {
	if !ensureAdmin(c) {
		return
	}

	var req QueryLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		jsonError(c, http.StatusBadRequest, err.Error())
		return
	}

	credentials, err := core.WithTS3Value(func(ts3Client *ts3.Client) (any, error) {
		return ts3Client.QueryLoginAdd(c.Request.Context(), req.ClientDBID, req.ServerID)
	})
	if err != nil {
		jsonError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": credentials})
}

func DeleteQueryLogin(c *gin.Context) {
	if !ensureAdmin(c) {
		return
	}

	cldbid, ok := parsePathInt(c, "cldbid")
	if !ok {
		return
	}

	err := core.WithTS3(func(ts3Client *ts3.Client) error {
		return ts3Client.QueryLoginDelete(c.Request.Context(), cldbid)
	})
	if err != nil {
		jsonError(c, http.StatusInternalServerError, err.Error())
		return
	}

	jsonMessage(c, http.StatusOK, "Query login deleted")
}

type ServerGroupClientReq struct {
	ServerGroupID int `json:"server_group_id" binding:"required"`
	ClientDBID    int `json:"client_db_id" binding:"required"`
}

func AddServerGroupClient(c *gin.Context) {
	if !ensureAdmin(c) {
		return
	}

	var req ServerGroupClientReq
	if err := c.ShouldBindJSON(&req); err != nil {
		jsonError(c, http.StatusBadRequest, err.Error())
		return
	}

	err := core.WithTS3(func(ts3Client *ts3.Client) error {
		return ts3Client.ServerGroupAddClient(c.Request.Context(), req.ServerGroupID, req.ClientDBID)
	})
	if err != nil {
		jsonError(c, http.StatusInternalServerError, err.Error())
		return
	}

	jsonMessage(c, http.StatusOK, "Client added to server group")
}

func RemoveServerGroupClient(c *gin.Context) {
	if !ensureAdmin(c) {
		return
	}

	var req ServerGroupClientReq
	if err := c.ShouldBindJSON(&req); err != nil {
		jsonError(c, http.StatusBadRequest, err.Error())
		return
	}

	err := core.WithTS3(func(ts3Client *ts3.Client) error {
		return ts3Client.ServerGroupDelClient(c.Request.Context(), req.ServerGroupID, req.ClientDBID)
	})
	if err != nil {
		jsonError(c, http.StatusInternalServerError, err.Error())
		return
	}

	jsonMessage(c, http.StatusOK, "Client removed from server group")
}
