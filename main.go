package main

import (
	"ivm-controller/initEnv"
	"ivm-controller/routes"
)

/**
控制台页面
 */

func main() {
	initEnv.InitLogger()
	initEnv.InitDb()
	initEnv.InitRedis()
	routes.InitRouter()
}