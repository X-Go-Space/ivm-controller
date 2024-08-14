package service

import (
	"fmt"
	"github.com/google/uuid"
	"ivm-controller/initEnv"
	"ivm-controller/model"
)

func CheckUserDir(dirName string) bool{
	var userDir model.UserDirectory
	initEnv.Db.Where("name = ?", dirName).First(&userDir)
	if userDir.ID == "" {
		return false
	}
	return true
}

func CreateUserDirectory(userDir model.UserDirectory) error{
	if CheckUserDir(userDir.Name) {
		return fmt.Errorf("the userDir is exit")
	}
	userDir.ID = uuid.New().String()

	result:= initEnv.Db.Create(&userDir)
	if result.Error != nil {
		initEnv.Logger.Error("add userDir failed,err:", result.Error)
		return result.Error
	}
	return nil

}
