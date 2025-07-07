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
type User struct {
	Userid   int       `gorm:"primaryKey;not null"`
	Username string    `gorm:"unique;not null" json:"username" form:"username"`
	Password string    `gorm:"not null" json:"password" form:"password"`
	BirthDay time.Time `gorm:"column:birthday;type:date;not null" json:"birthday" form:"birthday"`
}

// JWT声明结构体
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// 生成JWT Token
func generateToken(username string) (string, error) {
	// 定义签名密钥
	var mySigningKey = []byte(os.Getenv("JWT_SECRET"))

	// 创建一个新的声明，设置一些自定义字段
	claims := Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(), // 设置过期时间为1小时
			Issuer:    "GoJWT",
		},
	}

	// 创建一个新的token对象，指定签名算法和声明
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 用签名密钥签名生成token字符串
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

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

	router := gin.Default()
	router.POST("/register", func(c *gin.Context) {
		var user User
		if err := c.ShouldBind(&user); err != nil {
			c.JSON(400, gin.H{"error": "invalid input"})
			return
		}

		if err := db.Where("username = ?", user.Username).First(&user).Error; err == nil {
			c.JSON(400, gin.H{"error": "username already exists"})
			return
		}
		birthdayTime, err2 := time.Parse("2006-01-02", user.BirthDay.Format("2006-01-02"))
		if err2 != nil {
			return
		}
		newUser := User{Username: user.Username, Password: user.Password, BirthDay: birthdayTime}
		if err := db.Create(&newUser).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to create user"})
			return
		}
		token, err := generateToken(user.Username)
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to generate token"})
			return
		}
		c.JSON(200, gin.H{"token": token})
	})

	router.POST("/login", func(c *gin.Context) {
		var user User
		if err := c.ShouldBind(&user); err != nil {
			c.JSON(400, gin.H{"error": "invalid input"})
			return
		}
		if user.Username == "" || user.Password == "" {
			c.JSON(400, gin.H{"error": "username and password are required"})
			return
		}
		if err := db.Where("username = ? AND password = ?", user.Username, user.Password).First(&user).Error; err != nil {
			c.JSON(401, gin.H{"error": "invalid username or password"})
			return
		}
		token, err := generateToken(user.Username)
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to generate token"})
			return
		}
		c.JSON(200, gin.H{"token": token})

	})

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
