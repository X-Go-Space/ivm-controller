package controller

import (
	"github.com/gin-gonic/gin"
	"ivm-controller/initEnv"
	"ivm-controller/model"
	"ivm-controller/service"
	"ivm-controller/utils"
	"ivm-controller/utils/errmsg"
)

func CreateResource(ctx *gin.Context) {
	var resource model.Resource
	err := ctx.ShouldBind(&resource)
	if err != nil {
		initEnv.Logger.Error("bind resource fail,err:",err)
		utils.Err(errmsg.ErrMsg["CREATE_RESOURCE_FAIL"], errmsg.CREATE_RESOURCE_FAIL, ctx)
		return
	}
	data ,err:= service.CreateResource(resource, ctx)
	if err != nil {
		initEnv.Logger.Error("create resource failed, err: ", err)
		utils.Err(errmsg.ErrMsg["CREATE_RESOURCE_FAIL"], errmsg.CREATE_RESOURCE_FAIL, ctx)
		return
	}
	utils.OK(data,ctx)
}