package api

import (
	"time"

	"ati-study-jwt/model"
	"ati-study-jwt/repository"
	"ati-study-jwt/service"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func UserRegisterHandler(userRepo *repository.UserRepository) gin.HandlerFunc {
	userService := service.NewUserService(userRepo)

	return func(c *gin.Context) {
		var user struct {
			Username string `json:"username" form:"username"`
			Password string `json:"password" form:"password"`
			BirthDay string `json:"birthday" form:"birthday"`
		}
		if err := c.ShouldBind(&user); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}
		birthday, parseErr := time.Parse("2006-01-02", user.BirthDay)
		if parseErr != nil {
			c.JSON(400, gin.H{"error": "Invalid date format"})
			return
		}

		token, err := userService.Register(user.Username, user.Password, birthday)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, gin.H{"token": token})
	}
}

func LoginHandler(userRepo *repository.UserRepository) gin.HandlerFunc {
	var loginData service.LoginData
	userService := service.NewUserService(userRepo)
	return func(c *gin.Context) {
		if err := c.ShouldBind(&loginData); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		token, err := userService.Login(loginData.Username, loginData.Password)
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"token": token})
	}
}

func BirthdayHandler(userRepo *repository.UserRepository, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, exists := c.Get("username")
		if !exists {
			c.JSON(500, gin.H{"error": "failed to get username from token"})
			return
		}
		var user model.User
		if err := db.Where("username = ?", username).First(&user).Error; err != nil {
			c.JSON(404, gin.H{"error": "user not found"})
			return
		}
		c.JSON(200, gin.H{"birthday": user.BirthDay.Format("2006-01-02")})
	}
}
