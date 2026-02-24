package api

import (
	"Ts3Panel/core"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jkesh/ts3-go/ts3"
	ts3models "github.com/jkesh/ts3-go/ts3/models"
)

type kickClientReq struct {
	Reason string `json:"reason"`
}

func ListClients(c *gin.Context) {
	clients, err := core.WithTS3Value(func(ts3Client *ts3.Client) ([]ts3models.OnlineClient, error) {
		return ts3Client.ClientList(c.Request.Context(), "-uid", "-away", "-groups")
	})
	if err != nil {
		jsonError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, clients)
}

func KickClient(c *gin.Context) {
	if !ensureAdmin(c) {
		return
	}

	clid, ok := parsePathInt(c, "id")
	if !ok {
		return
	}

	var req kickClientReq
	_ = c.ShouldBindJSON(&req)
	reason := req.Reason
	if reason == "" {
		reason = "Kicked by WebAdmin"
	}

	err := core.WithTS3(func(ts3Client *ts3.Client) error {
		return ts3Client.KickFromServer(c.Request.Context(), clid, reason)
	})
	if err != nil {
		jsonError(c, http.StatusInternalServerError, err.Error())
		return
	}

	jsonMessage(c, http.StatusOK, "Kicked")
}
