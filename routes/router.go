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
		consoleApi.GET("/getAuthServerList", controller.GetAuthServerList)
		consoleApi.POST("/createResource", controller.CreateResource)
	}

	authApi := router.Group("api/v1")
	{
		authApi.POST("/login", controller.Login)
		authApi.GET("/authConfig", controller.AuthConfig)
		authApi.GET("/generateQrCode", controller.GenerateQrCode)
		authApi.GET("/getQrcodeStatus", controller.GetQrcodeStatus)
		authApi.POST("/mobileQrcodeLogin", controller.MobileQrcodeLogin)
		authApi.POST("/mobileQrcodeConfirm", controller.MobileQrcodeConfirm)
	}


	router.Run(utils.HttpPort)
}