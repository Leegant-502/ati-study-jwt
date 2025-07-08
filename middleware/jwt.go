package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"os"
	"time"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// 生成JWT Token
func GenerateToken(username string) (string, error) {
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

// JWT 认证中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(401, gin.H{"error": "authorization header is required"})
			c.Abort()
			return
		}

		// 解析JWT Token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			c.JSON(500, gin.H{"error": "could not parse claims"})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("username", claims.Username)
		c.Next()
	}
}
