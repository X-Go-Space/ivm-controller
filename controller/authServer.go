package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
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
	var authServers []model.AuthServer
	result := initEnv.Db.Find(&authServers)
	if result.Error != nil {
		initEnv.Logger.Error("GetAuthServerList get auth servers from mysql failed, err:", result.Error)
		utils.Err(errmsg.ErrMsg["GET_AUTH_LIST_FAIL"], errmsg.GET_AUTH_LIST_FAIL, ctx)
	}

	for idx, _ := range authServers {
		err:= json.Unmarshal([]byte(authServers[idx].AuthConfigJson), &authServers[idx].AuthConfig)
		if err != nil {
			initEnv.Logger.Error("GetAuthServerList unmarshal failed, err:", result.Error)
			utils.Err(errmsg.ErrMsg["GET_AUTH_LIST_FAIL"], errmsg.GET_AUTH_LIST_FAIL, ctx)
		}
	}
	utils.OK(authServers, ctx)

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

