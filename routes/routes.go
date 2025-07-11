package routes

import (
	"ati-study-jwt/middleware"
	"ati-study-jwt/model"
	"ati-study-jwt/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

func SetupRoutes(router *gin.Engine, userService *service.UserService, db *gorm.DB) {
	// 用户注册路由
	router.POST("/register", func(c *gin.Context) {
		var user struct {
			Username string    `gorm:"unique;not null" json:"username" form:"username"`
			Password string    `gorm:"not null" json:"password" form:"password"`
			BirthDay time.Time `gorm:"column:birthday;type:date;not null" json:"birthday" form:"birthday"`
		}
		if err := c.ShouldBind(&user); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input", "details": err.Error()})
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
	})

	// 需要JWT认证的路由组
	authRoutes := router.Group("/")
	authRoutes.Use(middleware.JWTAuthMiddleware())
	{
		// 生日查询路由
		authRoutes.GET("/birthday", func(c *gin.Context) {
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
		})
	}
}
