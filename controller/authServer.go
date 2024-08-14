package controller

import (
	"github.com/gin-gonic/gin"
	"ivm-controller/initEnv"
	"ivm-controller/model"
	"ivm-controller/service"
	"ivm-controller/utils"
	"ivm-controller/utils/errmsg"
)

func AuthServeCreate(ctx *gin.Context)  {
	var authServer model.AuthServer
	err := ctx.ShouldBind(&authServer)
	if err != nil {
		initEnv.Logger.Error("bind auth server fail,err:", err)
		utils.Err(errmsg.ErrMsg["CREATE_AUTH_SERVER_FAIL"], errmsg.CREATE_AUTH_SERVER_FAIL, ctx)
		return
	}
	err = service.CreatAuthServer(authServer)
	if err != nil {
		initEnv.Logger.Error("create auth server fail,err:", err)
		utils.Err(errmsg.ErrMsg["CREATE_AUTH_SERVER_FAIL"], errmsg.CREATE_AUTH_SERVER_FAIL, ctx)
		return
	}
	utils.OK(errmsg.ErrMsg["CREATE_AUTH_SERVER_SUCCESS"], ctx)
}

func GetAuthServerList(ctx *gin.Context)  {

}
func GetAuthServerById(ctx *gin.Context)  {
	userDirectoryId := ctx.Query("id")
	data, err := service.GetAuthServerById(userDirectoryId)
	if err !=nil {
		utils.Err(errmsg.ErrMsg["GET_AUTH_SERVER_FAIL"], errmsg.GET_AUTH_SERVER_FAIL, ctx)
		return
	}
	utils.OK(data, ctx)
}

