package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

// 定义一个自定义的声明结构体
type User struct {
	Userid   int       `gorm:"primaryKey;not null"`
	Username string    `gorm:"unique;not null"`
	Password string    `gorm:"not null"`
	BirthDay time.Time `gorm:"column:birthday;type:date;not null"`
}

// JWT声明结构体
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// 生成JWT Token
func generateToken(username string) (string, error) {
	// 定义签名密钥
	var mySigningKey = []byte("secret")

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
	// 初始化数据库连接
	dsn := "postgres://leegant:woshi@localhost:5432/atifactory"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("连接数据库失败", err)
		return
	}

	router := gin.Default()
	router.POST("/register", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		birthday := c.PostForm("birthday")
		if username == "" || password == "" {
			c.JSON(400, gin.H{"error": "username and password are required"})
			return
		}
		var user User
		if err := db.Where("username = ?", username).First(&user).Error; err == nil {
			c.JSON(400, gin.H{"error": "username already exists"})
			return
		}
		birthdayTime, err2 := time.Parse("2006-01-02", birthday)
		if err2 != nil {
			return
		}
		newUser := User{Username: username, Password: password, BirthDay: birthdayTime}
		if err := db.Create(&newUser).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to create user"})
			return
		}
		token, err := generateToken(username)
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to generate token"})
			return
		}
		c.JSON(200, gin.H{"token": token})
	})

	router.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		if username == "" || password == "" {
			c.JSON(400, gin.H{"error": "username and password are required"})
			return
		}
		var user User
		if err := db.Where("username = ? AND password = ?", username, password).First(&user).Error; err != nil {
			c.JSON(401, gin.H{"error": "invalid username or password"})
			return
		}
		token, err := generateToken(username)
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
