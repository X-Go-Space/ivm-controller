package service

import (
	"fmt"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"ivm-controller/initEnv"
	"ivm-controller/model"
)

func CheckAuthServer(authServerName string) bool {
	var count int
	err := initEnv.Db.Raw("SELECT COUNT(*) FROM auth_server WHERE name = ?", authServerName).Scan(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}
func CreatAuthServer(authServer model.AuthServer) error{
	if CheckAuthServer(authServer.Name) {
		return fmt.Errorf("the authserver is exit")
	}

	authServer.Id = uuid.New().String()
	jsonData, err := json.Marshal(authServer.AuthConfig)
	if err != nil {
		return err
	}
	extData, err := json.Marshal(authServer.Ext)
	if err != nil {
		return err
	}

	authServer.AuthConfigJson = string(jsonData)
	authServer.ExtJson = string(extData)
	result:= initEnv.Db.Create(&authServer)
	if result.Error != nil {
		initEnv.Logger.Error("create authServer, err:", result.Error)
		return result.Error
	}
	return nil
}

func GetAuthServerById(id string) (interface{}, error) {
	var authServerData model.AuthServer
	err := initEnv.Db.Where("id = ?", id).First(&authServerData).Error
	if err != nil {
		initEnv.Logger.Error("get user by sql failed, err:", err)
		return nil, err
	}
	err = json.Unmarshal([]byte(authServerData.AuthConfigJson), &authServerData.AuthConfig)
	if err != nil {
		initEnv.Logger.Error("Unmarshal auth config failed, err:", err)
		return nil, err
	}
	return authServerData, nil
}