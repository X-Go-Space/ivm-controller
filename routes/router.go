package routes

import (
	"github.com/gin-gonic/gin"
	"ivm-controller/controller"
	"ivm-controller/middleware"
	"ivm-controller/utils"
)

func InitRouter() {
	router := gin.Default()
	gin.SetMode(utils.AppMode)
	router.Use(middleware.Cors())
	router.Use(middleware.GlobalErrorInterceptor())

	// 后台对应接口
	consoleApi := router.Group("api/v1")
	consoleApi.Use(middleware.IsLoginMiddleWare())
	{
		consoleApi.GET("/getUsers", controller.GetUsers)
		consoleApi.POST("/addUser", controller.AddUser)
		consoleApi.POST("/createUserDirectory", controller.CreateUserDirectory)
		consoleApi.POST("/authServeCreate", controller.AuthServeCreate)
		consoleApi.GET("/getAuthServerById", controller.GetAuthServerById)
	}

	authApi := router.Group("api/v1")
	{
		authApi.POST("/login", controller.Login)
	}


	router.Run(utils.HttpPort)
}