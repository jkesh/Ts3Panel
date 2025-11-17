package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	App struct {
		JWTSecret string `mapstructure:"jwt_secret"`
		Port      string `mapstructure:"port"`
	} `mapstructure:"app"`

	// [新增] 数据库配置
	Database struct {
		Driver   string `mapstructure:"driver"`
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		DBName   string `mapstructure:"dbname"`
		SSLMode  string `mapstructure:"sslmode"`
		TimeZone string `mapstructure:"timezone"`
	} `mapstructure:"database"`

	TS3 struct {
		Protocol string `mapstructure:"protocol"`
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
	} `mapstructure:"ts3"`
}

var GlobalConfig *Config

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// 支持环境变量覆盖，例如 APP_DATABASE_PASSWORD
	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("[Config] Warning: %v, using defaults", err)
	}

	GlobalConfig = &Config{}
	if err := viper.Unmarshal(GlobalConfig); err != nil {
		log.Fatalf("[Config] Failed to parse config: %v", err)
	}

	// 设置默认值
	if GlobalConfig.Database.Driver == "" {
		GlobalConfig.Database.Driver = "postgres"
	}
	if GlobalConfig.Database.Host == "" {
		GlobalConfig.Database.Host = "127.0.0.1"
	}
	if GlobalConfig.Database.Port == 0 {
		GlobalConfig.Database.Port = 5432
	}
	if GlobalConfig.Database.SSLMode == "" {
		GlobalConfig.Database.SSLMode = "disable"
	}
}
