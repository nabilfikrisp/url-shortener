package database

import (
	"fmt"

	"github.com/nabilfikrisp/url-shortener/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	logMode := logger.Info
	if cfg.GoEnv == config.Production {
		logMode = logger.Silent
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logMode),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	fmt.Printf("Connected to %s in %s mode\n", db.Name(), cfg.GoEnv)
	return db, nil
}
