package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type GoEnv string

const (
	Production  GoEnv = "production"
	Development GoEnv = "development"
)

type Config struct {
	Port        string
	DatabaseURL string
	GoEnv       GoEnv
}

func Load() *Config {
	envValue := os.Getenv("GO_ENV")
	if envValue == "" {
		envValue = string(Development)
	}

	goEnv := GoEnv(envValue)
	if goEnv != Production {
		if err := godotenv.Load(); err != nil {
			panic(fmt.Sprintf("missing .env file: %v", err))
		}
	}

	return &Config{
		Port:        verifyEnv("PORT"),
		DatabaseURL: verifyEnv("DATABASE_URL"),
		GoEnv:       goEnv,
	}
}

func verifyEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("missing or empty required env var: %s", key))
	}
	return val
}
