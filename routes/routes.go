package routes

import (
	"ati-study-jwt/api"
	"ati-study-jwt/middleware"
	"ati-study-jwt/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, userRepo *repository.UserRepository, db *gorm.DB) {
	// 用户注册路由
	router.POST("/register", api.UserRegisterHandler(userRepo))

	// 用户登录路由
	router.POST("/login", api.LoginHandler(userRepo))

	// 需要JWT认证的路由组
	authRoutes := router.Group("/")
	authRoutes.Use(middleware.JWTAuthMiddleware())
	{
		// 生日查询路由
		authRoutes.GET("/birthday", api.BirthdayHandler(userRepo, db))
	}
}
