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

func GetResourceUrlByResourceIds(ctx *gin.Context) {
	var resourceIdData map[string]interface{}
	err := ctx.ShouldBind(&resourceIdData)
	if err != nil {
		initEnv.Logger.Error("GetResourceByResourceIds bind resource fail,err:",err)
		utils.Err(errmsg.ErrMsg["GET_RESOURCE_FAIL"], errmsg.GET_RESOURCE_FAIL, ctx)
		return
	}
	resourceIds := utils.ReadNestedData(resourceIdData, "resourceIds").([]interface{})
	data:=make([]string, 0)
	for _, resourceId := range resourceIds {
		resourceId, _ := resourceId.(string)
		resourceSessId := utils.GenerateResourceId(resourceId)
		redirectUrl, err := initEnv.Redis.Get(ctx, resourceSessId).Result()
		if err != nil {
			initEnv.Logger.Error("GetResourceByResourceIds get resource fail,err:",err)
			continue
		}
		data = append(data, redirectUrl)
	}
	utils.OK(data, ctx)
}