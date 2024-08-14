package controller

import (
	"github.com/gin-gonic/gin"
	"ivm-controller/initEnv"
	"ivm-controller/model"
	"ivm-controller/service"
	"ivm-controller/utils"
	"ivm-controller/utils/errmsg"
)

func CreateUserDirectory(ctx *gin.Context) {
    var userDir model.UserDirectory
	err:=ctx.ShouldBind(&userDir)
	if err != nil {
		initEnv.Logger.Error("user dir bind failed, err: ", err)
		utils.Err(errmsg.ErrMsg["CREATE_USER_DIRECTORY_FAILED"], errmsg.CREATE_USER_DIRECTORY_FAILED, ctx)
		return
	}

	err = service.CreateUserDirectory(userDir)
	if err != nil {
		initEnv.Logger.Error("user dir create failed, err: ", err)
		utils.Err(errmsg.ErrMsg["CREATE_USER_DIRECTORY_FAILED"], errmsg.CREATE_USER_DIRECTORY_FAILED, ctx)
		return
	}
	utils.OK(errmsg.ErrMsg["CREATE_USER_DIRECTORY_SUCCESS"], ctx)
}