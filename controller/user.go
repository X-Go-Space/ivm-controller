package controller

import (
	"github.com/gin-gonic/gin"
	"ivm-controller/initEnv"
	"ivm-controller/service"
	"ivm-controller/utils"
	"ivm-controller/utils/errmsg"
)

func GetUsers(ctx *gin.Context) {
	data, err := service.GetUsers(ctx)
	if err != nil {
		initEnv.Logger.Error("get users fail, err: ", err)
		utils.Err(errmsg.ErrMsg["GET_USERS_FAIL"], errmsg.GET_USERS_FAIL, ctx)
		return
	}
	utils.OK(data, ctx)
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