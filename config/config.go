package config

import (
	"fmt"
	"strconv"

	"github.com/spf13/viper"
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
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
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
