package api

import (
	"Ts3Panel/database"
	"Ts3Panel/models"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var botHTTPClient = &http.Client{Timeout: 10 * time.Second}

type BotInfo struct {
	ID     int    `json:"Id"`
	Name   string `json:"Name"`
	Status int    `json:"Status"` // 0=Offline, 1=Booting, 2=Connecting, 3=Online, 4=Transferring, 5=Stopping
}

type BotResponse struct {
	models.MusicBot
	RealStatus  string `json:"real_status"` // Online, Offline, Playing...
	CurrentSong string `json:"current_song"`
}

type botCommandReq struct {
	Command string `json:"command" binding:"required"`
	Value   string `json:"value"`
}

func endpointKey(bot models.MusicBot) string {
	return bot.ApiUrl + "|" + bot.ApiToken
}

// isSystemCmd: true 表示是对整个程序的指令(如 connect)，false 表示是对具体 Bot 的指令(如 play)
func sendBotCommand(bot *models.MusicBot, isSystemCmd bool, cmd string, args ...string) (string, error) {
	if bot.ApiUrl == "" {
		return "", errors.New("bot api_url is empty")
	}

	base := strings.TrimRight(bot.ApiUrl, "/")
	var apiPath string
	if isSystemCmd {
		apiPath = fmt.Sprintf("%s/api/bot/%s", base, strings.TrimPrefix(cmd, "/"))
	} else {
		apiPath = fmt.Sprintf("%s/api/bot/use/%d/%s", base, bot.BotId, strings.TrimPrefix(cmd, "/"))
	}
	for _, v := range args {
		apiPath += "/" + url.PathEscape(v)
	}

	req, err := http.NewRequest(http.MethodGet, apiPath, nil)
	if err != nil {
		return "", err
	}

	q := req.URL.Query()
	q.Set("token", bot.ApiToken)
	req.URL.RawQuery = q.Encode()

	resp, err := botHTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	if resp.StatusCode >= http.StatusBadRequest {
		return "", fmt.Errorf("bot api error (%d): %s", resp.StatusCode, bodyString)
	}

	return bodyString, nil
}

func AddBot(c *gin.Context) {
	var req models.MusicBot
	if err := c.ShouldBindJSON(&req); err != nil {
		jsonError(c, http.StatusBadRequest, err.Error())
		return
	}
	if req.ApiUrl == "" {
		req.ApiUrl = "http://127.0.0.1:58913"
	}
	req.BotId = -1

	if err := database.DB.Create(&req).Error; err != nil {
		jsonError(c, http.StatusInternalServerError, "Failed to save bot")
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "Bot added", "data": req})
}

func ListBots(c *gin.Context) {
	var dbBots []models.MusicBot
	if err := database.DB.Find(&dbBots).Error; err != nil {
		jsonError(c, http.StatusInternalServerError, "Failed to query bots")
		return
	}

	activeByEndpoint := make(map[string]map[int]BotInfo, len(dbBots))
	seenEndpoint := make(map[string]struct{}, len(dbBots))

	for _, dbBot := range dbBots {
		key := endpointKey(dbBot)
		if _, ok := seenEndpoint[key]; ok {
			continue
		}
		seenEndpoint[key] = struct{}{}

		listJSON, err := sendBotCommand(&dbBot, true, "list")
		if err != nil {
			continue
		}

		var activeList []BotInfo
		if err := json.Unmarshal([]byte(listJSON), &activeList); err != nil {
			continue
		}

		activeMap := make(map[int]BotInfo, len(activeList))
		for _, b := range activeList {
			activeMap[b.ID] = b
		}
		activeByEndpoint[key] = activeMap
	}

	responseList := make([]BotResponse, 0, len(dbBots))
	for _, dbBot := range dbBots {
		res := BotResponse{MusicBot: dbBot}
		res.ApiToken = "******"
		res.RealStatus = "Offline"

		if byID, ok := activeByEndpoint[endpointKey(dbBot)]; ok {
			if activeBot, ok := byID[dbBot.BotId]; ok {
				res.RealStatus = "Online"
				if activeBot.Status == 3 {
					if _, err := sendBotCommand(&dbBot, false, "song"); err == nil {
						res.RealStatus = "Playing"
					} else {
						res.RealStatus = "Idle"
					}
				}
			}
		}

		responseList = append(responseList, res)
	}

	c.JSON(http.StatusOK, gin.H{"data": responseList})
}

func DeleteBot(c *gin.Context) {
	id, ok := parsePathInt(c, "id")
	if !ok {
		return
	}

	var bot models.MusicBot
	if err := database.DB.First(&bot, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			jsonError(c, http.StatusNotFound, "Bot not found")
			return
		}
		jsonError(c, http.StatusInternalServerError, "Failed to query bot")
		return
	}

	if bot.BotId >= 0 {
		_, _ = sendBotCommand(&bot, true, "delete", strconv.Itoa(bot.BotId))
	}

	if err := database.DB.Delete(&models.MusicBot{}, id).Error; err != nil {
		jsonError(c, http.StatusInternalServerError, "Failed to delete bot")
		return
	}

	jsonMessage(c, http.StatusOK, "Deleted")
}

func SendCommand(c *gin.Context) {
	id, ok := parsePathInt(c, "id")
	if !ok {
		return
	}

	var req botCommandReq
	if err := c.ShouldBindJSON(&req); err != nil {
		jsonError(c, http.StatusBadRequest, err.Error())
		return
	}
	req.Command = strings.TrimSpace(req.Command)

	var bot models.MusicBot
	if err := database.DB.First(&bot, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			jsonError(c, http.StatusNotFound, "Bot not found")
			return
		}
		jsonError(c, http.StatusInternalServerError, "Failed to query bot")
		return
	}

	var (
		resp string
		err  error
	)

	switch req.Command {
	case "play":
		if req.Value == "" {
			jsonError(c, http.StatusBadRequest, "value is required for play")
			return
		}
		resp, err = sendBotCommand(&bot, false, "play", req.Value)
	case "pause":
		resp, err = sendBotCommand(&bot, false, "pause")
	case "stop":
		resp, err = sendBotCommand(&bot, false, "stop")
	case "volume":
		if req.Value == "" {
			jsonError(c, http.StatusBadRequest, "value is required for volume")
			return
		}
		resp, err = sendBotCommand(&bot, false, "volume", req.Value)
	case "connect":
		target := strings.TrimSpace(req.Value)
		if target == "" {
			target = bot.ServerAddr
		}
		if target == "" {
			target = "127.0.0.1"
		}

		resp, err = sendBotCommand(&bot, true, "connect/to", target)
		if err == nil {
			var newBotInfo BotInfo
			if unmarshalErr := json.Unmarshal([]byte(resp), &newBotInfo); unmarshalErr == nil {
				bot.BotId = newBotInfo.ID
				_ = database.DB.Save(&bot).Error
			}
		}
	case "raw":
		parts := strings.Split(req.Value, "/")
		if len(parts) == 0 || strings.TrimSpace(parts[0]) == "" {
			jsonError(c, http.StatusBadRequest, "value is required for raw")
			return
		}
		cmd := parts[0]
		args := parts[1:]
		resp, err = sendBotCommand(&bot, false, cmd, args...)
	default:
		jsonError(c, http.StatusBadRequest, "Unsupported command")
		return
	}

	if err != nil {
		jsonError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "Executed", "bot_response": resp})
}
