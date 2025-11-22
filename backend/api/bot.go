package api

import (
	"Ts3Panel/database"
	"Ts3Panel/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// TS3AudioBot 返回的机器人信息结构
type BotInfo struct {
	Id     int    `json:"Id"`
	Name   string `json:"Name"`
	Status int    `json:"Status"` // 0=Offline, 1=Booting, 2=Connecting, 3=Online, 4=Transferring, 5=Stopping
}

// 用于前端显示的扩展结构
type BotResponse struct {
	models.MusicBot
	RealStatus  string `json:"real_status"` // Online, Offline, Playing...
	CurrentSong string `json:"current_song"`
}

// 发送指令给 TS3AudioBot
// isSystemCmd: true 表示是对整个程序的指令(如 connect)，false 表示是对具体 Bot 的指令(如 play)
func sendBotCommand(bot *models.MusicBot, isSystemCmd bool, cmd string, args ...string) (string, error) {
	var apiPath string
	if isSystemCmd {
		apiPath = fmt.Sprintf("%s/api/bot/%s", bot.ApiUrl, cmd)
	} else {
		apiPath = fmt.Sprintf("%s/api/bot/use/%d/%s", bot.ApiUrl, bot.BotId, cmd)
	}

	if len(args) > 0 {
		for _, v := range args {
			apiPath += "/" + url.PathEscape(v)
		}
	}

	req, err := http.NewRequest("GET", apiPath, nil)
	if err != nil {
		return "", err
	}

	q := req.URL.Query()
	q.Add("token", bot.ApiToken)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("bot api error (%d): %s", resp.StatusCode, bodyString)
	}

	return bodyString, nil
}

// AddBot 添加机器人
func AddBot(c *gin.Context) {
	var req models.MusicBot
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if req.ApiUrl == "" {
		req.ApiUrl = "http://127.0.0.1:58913"
	}
	// 默认 BotId 为 -1，表示未连接
	req.BotId = -1

	if err := database.DB.Create(&req).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to save bot"})
		return
	}
	c.JSON(200, gin.H{"msg": "Bot added", "data": req})
}

// ListBots 获取机器人列表（含实时状态同步）
func ListBots(c *gin.Context) {
	var dbBots []models.MusicBot
	database.DB.Find(&dbBots)

	var responseList []BotResponse

	// 1. 获取 TS3AudioBot 所有活动机器人列表
	// 我们假设第一个配置的 API 地址是主地址
	activeBotsMap := make(map[int]BotInfo)
	if len(dbBots) > 0 {
		// 调用 /api/bot/list 获取所有在线 Bot
		// 这里简单起见，使用第一条记录的 API 地址去查询
		listJson, err := sendBotCommand(&dbBots[0], true, "list")
		if err == nil {
			var activeList []BotInfo
			json.Unmarshal([]byte(listJson), &activeList)
			for _, b := range activeList {
				activeBotsMap[b.Id] = b
			}
		}
	}

	// 2. 合并状态
	for _, dbBot := range dbBots {
		res := BotResponse{MusicBot: dbBot}
		res.ApiToken = "******" // 隐藏 Token

		if val, ok := activeBotsMap[dbBot.BotId]; ok {
			// 机器人在运行中
			res.RealStatus = "Online"
			if val.Status == 3 { // Status 3 = Online
				// 可以进一步查询是否在播放: /api/bot/use/ID/song
				_, err := sendBotCommand(&dbBot, false, "song")
				if err == nil {
					res.RealStatus = "Playing"
					// 这里可以解析 songJson 获取歌名，暂时略过
				} else {
					res.RealStatus = "Idle"
				}
			}
		} else {
			res.RealStatus = "Offline"
		}
		responseList = append(responseList, res)
	}

	c.JSON(200, gin.H{"data": responseList})
}

// DeleteBot 删除
func DeleteBot(c *gin.Context) {
	id := c.Param("id")
	var bot models.MusicBot
	if err := database.DB.First(&bot, id).Error; err == nil {
		// 尝试停止并删除实例
		if bot.BotId >= 0 {
			sendBotCommand(&bot, true, "delete", fmt.Sprintf("%d", bot.BotId))
		}
	}
	database.DB.Delete(&models.MusicBot{}, id)
	c.JSON(200, gin.H{"msg": "Deleted"})
}

// SendCommand 执行指令
func SendCommand(c *gin.Context) {
	botID := c.Param("id")
	var req struct {
		Command string `json:"command"`
		Value   string `json:"value"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var bot models.MusicBot
	if err := database.DB.First(&bot, botID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Bot not found"})
		return
	}

	var err error
	var resp string

	switch req.Command {
	case "play":
		resp, err = sendBotCommand(&bot, false, "play", req.Value)
	case "pause":
		resp, err = sendBotCommand(&bot, false, "pause")
	case "stop":
		resp, err = sendBotCommand(&bot, false, "stop")
	case "volume":
		resp, err = sendBotCommand(&bot, false, "volume", req.Value)

	case "connect":
		// [关键] 连接指令
		target := bot.ServerAddr
		if target == "" {
			target = "127.0.0.1"
		}
		// 使用系统指令 connect/to/<ip>
		resp, err = sendBotCommand(&bot, true, "connect/to", target)

		// [关键] 如果连接成功，TS3AudioBot 会返回新 Bot 的信息 JSON
		// 我们解析它，拿到 ID，并更新数据库
		if err == nil {
			var newBotInfo BotInfo
			if jsonErr := json.Unmarshal([]byte(resp), &newBotInfo); jsonErr == nil {
				bot.BotId = newBotInfo.Id
				database.DB.Save(&bot) // 更新 BotId
			}
		}
	case "raw":
		// Value 格式: "cmd/arg1/arg2"
		// 我们需要把 Value 拆开传给 sendBotCommand
		parts := strings.Split(req.Value, "/")
		if len(parts) > 0 {
			cmd := parts[0]
			args := parts[1:]
			resp, err = sendBotCommand(&bot, false, cmd, args...)
		}
	}

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"msg": "Executed", "bot_response": resp})
}
