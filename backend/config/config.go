package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"strings"
)

type AppConfig struct {
	JWTSecret string `mapstructure:"jwt_secret"`
	Port      string `mapstructure:"port"`
}

type DatabaseConfig struct {
	Driver   string `mapstructure:"driver"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
	TimeZone string `mapstructure:"timezone"`
}

type TS3Config struct {
	Protocol         string `mapstructure:"protocol"` // tcp | ssh | webquery
	Host             string `mapstructure:"host"`
	Port             int    `mapstructure:"port"`
	User             string `mapstructure:"user"`
	Password         string `mapstructure:"password"`
	ServerID         int    `mapstructure:"server_id"`
	APIKey           string `mapstructure:"api_key"`
	HTTPS            bool   `mapstructure:"https"`
	BasePath         string `mapstructure:"base_path"`
	FallbackProtocol string `mapstructure:"fallback_protocol"` // empty | tcp
	FallbackHost     string `mapstructure:"fallback_host"`
	FallbackPort     int    `mapstructure:"fallback_port"`
}

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	TS3      TS3Config      `mapstructure:"ts3"`
}

var GlobalConfig *Config

func LoadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// 支持环境变量覆盖，例如 APP_DATABASE_PASSWORD
	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("[Config] Warning: %v, using defaults", err)
	}

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	// 设置默认值
	if cfg.Database.Driver == "" {
		cfg.Database.Driver = "postgres"
	}
	if cfg.Database.Host == "" {
		cfg.Database.Host = "127.0.0.1"
	}
	if cfg.Database.Port == 0 {
		cfg.Database.Port = 5432
	}
	if cfg.Database.SSLMode == "" {
		cfg.Database.SSLMode = "disable"
	}
	if cfg.TS3.ServerID == 0 {
		cfg.TS3.ServerID = 1
	}
	if cfg.TS3.Protocol == "" {
		cfg.TS3.Protocol = "webquery"
	}
	cfg.TS3.Protocol = strings.ToLower(strings.TrimSpace(cfg.TS3.Protocol))
	cfg.TS3.FallbackProtocol = strings.ToLower(strings.TrimSpace(cfg.TS3.FallbackProtocol))
	if cfg.TS3.Port == 0 {
		switch cfg.TS3.Protocol {
		case "ssh":
			cfg.TS3.Port = 10022
		case "webquery":
			if cfg.TS3.HTTPS {
				cfg.TS3.Port = 10443
			} else {
				cfg.TS3.Port = 10080
			}
		default:
			cfg.TS3.Port = 10011
		}
	}
	if cfg.TS3.FallbackProtocol == "" {
		cfg.TS3.FallbackProtocol = "tcp"
	}
	if cfg.TS3.FallbackHost == "" {
		cfg.TS3.FallbackHost = cfg.TS3.Host
	}
	if cfg.TS3.FallbackPort == 0 && cfg.TS3.FallbackProtocol == "tcp" {
		cfg.TS3.FallbackPort = cfg.TS3.Port
	}
	if cfg.App.Port == "" {
		cfg.App.Port = ":8080"
	}

	GlobalConfig = cfg
	return nil
}
