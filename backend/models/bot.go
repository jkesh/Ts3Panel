package models

import "gorm.io/gorm"

type MusicBot struct {
	gorm.Model
	Name string `json:"name"`

	ApiUrl string `json:"api_url"`

	ApiToken string `json:"api_token"`

	BotId int `json:"bot_id" gorm:"default:-1"`

	ServerAddr string `json:"server_addr"`
}
