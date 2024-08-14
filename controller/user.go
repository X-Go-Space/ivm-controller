package controller

import (
	"github.com/gin-gonic/gin"
	"ivm-controller/initEnv"
	"ivm-controller/service"
	"ivm-controller/utils"
	"ivm-controller/utils/errmsg"
)

func GetUsers(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "GetUsers",
	})
	initEnv.Logger.Error("This is an error message.")
}

func AddUser(ctx *gin.Context) {
	err := service.AddUser(ctx)
	if err != nil {
		initEnv.Logger.Error("add user failed, the err is: ", err)
		utils.Err(errmsg.GetErrMsg(errmsg.ADD_FAILED), errmsg.ADD_FAILED, ctx)
		return
	}
	utils.OK(errmsg.ErrMsg["ADD_SUCCESS"], ctx)
}