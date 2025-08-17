package database

import (
	"fmt"
	"log"
	"time"

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

	var db *gorm.DB
	var err error

	// wait for db to spin up if usign docker
	for range 10 {
		db, err = gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
			Logger: logger.Default.LogMode(logMode),
		})

		if err == nil {
			break
		}

		log.Print("Retrying DB Connection")
		time.Sleep(time.Second)
	}

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
