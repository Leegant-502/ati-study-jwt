package utils

import (
	"github.com/dgrijalva/jwt-go"
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
