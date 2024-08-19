package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"ivm-controller/initEnv"
	"ivm-controller/model"
	"ivm-controller/utils"
	"ivm-controller/utils/errmsg"
)
func CheckResource(resourceName string) bool{
	var resource model.Resource
	initEnv.Db.Where("name = ?", resourceName).First(&resource)
	if resource.Id == "" {
		return false
	}
	return true
}
func CreateResource (resource model.Resource, ctx *gin.Context) (interface{}, error){
	if CheckResource(resource.Name) {
		initEnv.Logger.Error("the database have the resource or occurred error")
		return nil, fmt.Errorf("the database have the resource or occurred error")
	}
	resource.Id = uuid.New().String()
	// 创建的时候，将资源ID，塞进用户的会话信息里面，用户如果登录了，那么就塞进去，没登陆就不赛
	err := initEnv.Db.Create(&resource).Error
	if err != nil {
		initEnv.Logger.Error("add resource failed,err: ", err)
		return nil,err
	}
	// 同时将资源的redirectUrl给放进session里面，便于获取
	resourceRedisId := utils.GenerateResourceId(resource.Id)
	resourceData , err := json.Marshal(resource)
	if err!= nil {
		initEnv.Logger.Error("set redis marshal resource failed,err: ", err)
		return nil,err
	}
	err = initEnv.Redis.Set(ctx, resourceRedisId,resourceData, day).Err()
	if err != nil {
		initEnv.Logger.Error("set redis resource failed,err: ", err)
		return nil,err
	}
	userList := resource.UserList
	for _, user := range userList {
		userSid := utils.GenerateSessId(user.ID)
		var sessionData map[string]interface{}
		// 看一下Redis中是否有对应的用户
		val, err := initEnv.Redis.Get(ctx, userSid).Result()
	    if err != nil && err != redis.Nil {
			initEnv.Logger.Error("create resource get user session failed,err: ", err)
		    return nil, err
		} else if err == redis.Nil {
			continue
		}
		err = json.Unmarshal([]byte(val), &sessionData)
		if err != nil {
			initEnv.Logger.Error("create resource unmarshal  session data failed,err: ", err)
			return nil, err
		}
		userResources := utils.ReadNestedData(sessionData, "resources").([]interface{})
		userResources = append(userResources, resource.Id)
		utils.SetNestedValue(sessionData, "resources", userResources)

		sessionDataByte ,err := json.Marshal(&sessionData)
		if err != nil {
			initEnv.Logger.Error("create resource marshal  session data failed,err: ", err)
			return nil, err
		}
		err = initEnv.Redis.Set(ctx, userSid, sessionDataByte, day).Err()

		if err != nil {
			initEnv.Logger.Error("create resource set redis err, err:", err)
			return nil, err
		}
	}
	return errmsg.ErrMsg["CREATE_RESOURCE_SUCCESS"],nil

}
