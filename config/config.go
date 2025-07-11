package config

import (
	"github.com/spf13/viper"
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
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		panic("Failed to read config file: " + err.Error())
	}
	port, _ := strconv.Atoi(viper.GetString("database.port"))

	return &DatabaseConfig{
		Host:     viper.GetString("database.host"),
		Port:     port,
		User:     viper.GetString("database.user"),
		Password: viper.GetString("database.password"),
		Database: viper.GetString("database.name"),
	}
}
