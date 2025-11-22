package database

import (
	"Ts3Panel/config"
	"Ts3Panel/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	cfg := config.GlobalConfig.Database

	var dialector gorm.Dialector

	if cfg.Driver == "sqlite" {
		// SQLite 模式（保留作为备选）
		dbFile := cfg.DBName
		if dbFile == "" {
			dbFile = "panel.db"
		}
		dialector = sqlite.Open(dbFile)
		log.Println("[Database] Using SQLite driver.")
	} else {
		// PostgreSQL 模式
		// dsn 格式: "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
			cfg.Host,
			cfg.User,
			cfg.Password,
			cfg.DBName,
			cfg.Port,
			cfg.SSLMode,
			cfg.TimeZone,
		)
		dialector = postgres.Open(dsn)
		log.Printf("[Database] Connecting to PostgreSQL at %s:%d...", cfg.Host, cfg.Port)
	}

	DB, err = gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Fatalf("[Database] Connection failed: %v", err)
	}

	// 自动迁移表结构
	if err := DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("[Database] Migration failed: %v", err)
	}
	if err := DB.AutoMigrate(&models.User{}, &models.MusicBot{}); err != nil { // [!code ++]
		log.Fatalf("[Database] Migration failed: %v", err)
	}
	log.Println("[Database] Initialized successfully.")

}
