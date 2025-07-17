package repository

import (
	"ati-study-jwt/config"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DBInit() (*gorm.DB, error) {
	postgresConfig := config.PostgresConfig()
	db, err := gorm.Open(postgres.Open(
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable&parseTime=True",
			postgresConfig.User,
			postgresConfig.Password,
			postgresConfig.Host,
			postgresConfig.Port,
			postgresConfig.Database),
	), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	fmt.Println("Database initialized successfully")
	return db, nil
}
