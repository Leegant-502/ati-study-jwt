package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

// 定义一个自定义的声明结构体

// JWT声明结构体

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	// 初始化数据库连接
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("连接数据库失败", err)
		return
	}

	router.GET("/birthday", func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(401, gin.H{"error": "authorization header is required"})
			return
		}

		// 解析JWT Token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})
		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "invalid token"})
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			c.JSON(500, gin.H{"error": "could not parse claims"})
			return
		}

		var user User
		if err := db.Where("username = ?", claims.Username).First(&user).Error; err != nil {
			c.JSON(404, gin.H{"error": "user not found"})
			return
		}

		c.JSON(200, gin.H{"birthday": user.BirthDay.Format("2006-01-02")})
	})

	err = router.Run(":8080")
}
