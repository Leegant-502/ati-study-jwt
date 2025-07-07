package routes

import (
	"ati-study-jwt/service"
	"github.com/gin-gonic/gin"
	"time"
)

func SetupRoutes(router *gin.Engine, userService *service.UserService) {
	// 用户注册路由
	router.POST("/register", func(c *gin.Context) {
		var user struct {
			Username string    `gorm:"unique;not null" json:"username" form:"username"`
			Password string    `gorm:"not null" json:"password" form:"password"`
			BirthDay time.Time `gorm:"column:birthday;type:date;not null" json:"birthday" form:"birthday"`
		}
		if err := c.ShouldBind(&user); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		token, err := userService.Register(user.Username, user.Password, user.BirthDay)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, gin.H{"token": token})
	})

	// 用户登录路由
	router.POST("/login", func(c *gin.Context) {
		var loginData service.LoginData
		if err := c.ShouldBindJSON(&loginData); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})

			return
		}

		token, err := userService.Login(loginData.Username, loginData.Password)
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"token": token})
	})
}
