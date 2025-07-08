package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

func PostgresConfig() *DatabaseConfig {
	// 加载.env文件
	err := godotenv.Load()
	if err != nil {
		return nil
	}

	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))

	return &DatabaseConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     port,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_NAME"),
	}
}
