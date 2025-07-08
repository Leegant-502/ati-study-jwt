package main

import (
	"ati-study-jwt/repository"
	"ati-study-jwt/routes"
	"ati-study-jwt/service"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {

	db, err := repository.DBInit()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// 初始化各层
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	// 初始化路由
	router := gin.Default()

	// 设置路由
	routes.SetupRoutes(router, userService, db)

	// 启动服务器
	err = router.Run(":8080")
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
