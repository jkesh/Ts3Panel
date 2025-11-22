package models

import "gorm.io/gorm"

type MusicBot struct {
	gorm.Model
	Name string `json:"name"`

	// TS3AudioBot 的 API 地址，例如 "http://127.0.0.1:58913"
	ApiUrl string `json:"api_url"`

	// 鉴权 Token (从 !api token 获取)
	ApiToken string `json:"api_token"` // JSON 中隐藏 Token，避免泄露给前端

	// 机器人实例 ID，TS3AudioBot 支持多开，默认通常是 0
	// 这就是报错缺少的那个字段
	BotId int `json:"bot_id" gorm:"default:0"`

	// 机器人连接的目标 TS3 服务器信息 (用于 connect 指令)
	ServerAddr string `json:"server_addr"`
}
