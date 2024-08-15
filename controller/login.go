package controller

import (
	"github.com/gin-gonic/gin"
	"ivm-controller/initEnv"
	"ivm-controller/service"
	"ivm-controller/utils"
	"ivm-controller/utils/errmsg"
)

func Login(ctx *gin.Context) {
	var loginData map[string]interface{}
	err := ctx.BindJSON(&loginData)
	if err != nil {
		initEnv.Logger.Error("bind json failed,err:", err)
		utils.Err(errmsg.GetErrMsg(errmsg.ERROR), errmsg.ERROR, ctx)
		return
	}
	loginType := utils.ReadNestedData(loginData, "loginType")
	switch loginType {
	case "local":
		data, err := service.LoginByLocal(loginData, ctx)
		if err !=nil {
			initEnv.Logger.Error("login by local bind json failed,err:", err)
			utils.Err(errmsg.ErrMsg["LOGIN_BY_LOCAL_FAILED"], errmsg.LOGIN_BY_LOCAL_FAILED, ctx)
			break
		}
		utils.OK(data, ctx)
	case "http":
		data, err := service.LoginByHttp(loginData, ctx)
		if err !=nil {
			initEnv.Logger.Error("login by local bind json failed,err:", err)
			utils.Err(errmsg.ErrMsg["LOGIN_BY_HTTP_FAIL"], errmsg.LOGIN_BY_HTTP_FAIL, ctx)
			break
		}
		utils.OK(data, ctx)
	case "oauth2":
		data, err := service.LoginByOauth2(loginData, ctx)
		if err !=nil {
			initEnv.Logger.Error("login by local bind json failed,err:", err)
			utils.Err(errmsg.ErrMsg["LOGIN_BY_HTTP_FAIL"], errmsg.LOGIN_BY_HTTP_FAIL, ctx)
			break
		}
		utils.OK(data, ctx)
	default:
		initEnv.Logger.Error("login failed, unknow login type")
		utils.Err(errmsg.ErrMsg["UNKNOW_LOGIN_TYPE"], errmsg.UNKNOW_LOGIN_TYPE, ctx)
	}
}

func AuthConfig(ctx *gin.Context) {

}
