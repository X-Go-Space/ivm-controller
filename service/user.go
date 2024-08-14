package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"ivm-controller/initEnv"
	"ivm-controller/model"
	"ivm-controller/utils"
	"ivm-controller/utils/errmsg"
	"time"
)

const (
	saltLength = 16
	day = 24*time.Hour
)

func CheckUser(userName string) bool{
	var user model.User
	initEnv.Db.Where("user_name = ?", userName).First(&user)
	if user.ID == "" {
		return false
	}
	return true
}

func AddUser(ctx *gin.Context) error {
	var newUser model.User
	err := ctx.BindJSON(&newUser)
	if err != nil {
		return err
	}
	if CheckUser(newUser.UserName) {
		return fmt.Errorf("the user name is exit")
	}
	newUser.ID = uuid.New().String()

	rawPassword, err := utils.Decrypt(newUser.Password)
	if err != nil {
		initEnv.Logger.Error("decrypt user password failed,err:", err)
		return err
	}
	newUser.PwdSalt, err = utils.GenerateSalt(saltLength)
	if err != nil {
		initEnv.Logger.Error("generate salt failed, err:", err)
		return err
	}
	newUser.Password = utils.HashPasswordWithSalt(rawPassword, newUser.PwdSalt)

	result := initEnv.Db.Create(&newUser)
	if result.Error != nil {
		initEnv.Logger.Error("add user failed,err:", result.Error)
		return result.Error
	}
	return nil
}

func LoginByLocal (loginData map[string]interface{}, ctx *gin.Context) (interface{}, error) {
	// 本地密码登录
	userName := utils.ReadNestedData(loginData, "userName")
	password := utils.ReadNestedData(loginData, "password").(string)
	password, err := utils.Decrypt(password)
	if err != nil {
		initEnv.Logger.Error("when login decrypt password failed, err:", err)
		return nil, err
	}
	var user model.User
	result := initEnv.Db.Where("is_local = ? AND user_name = ?", "1", userName).First(&user)
	if result.Error != nil{
		initEnv.Logger.Error("when login query mysql failed, err:", err)
		return nil, err
	}

	password = utils.HashPasswordWithSalt(password, user.PwdSalt)
	if password != user.Password {
		initEnv.Logger.Error("the password is not right")
		return nil, err
	}

	// 将用户信息设置到redis里面
	// 将sid下发下去
	sessionKey := utils.GenerateSessId(user.ID)
	sessData, err := json.Marshal(gin.H{
		"id": user.ID,
		"userName": user.UserName,
		"userDirectoryId": user.UserDirectoryId,
		"mobile": user.Mobile,
		"email": user.Email,
		"status": user.Status,
	})
	if err != nil {
		initEnv.Logger.Error("marshal sess data failed,err:", err)
		return nil, err
	}

	err = initEnv.Redis.Set(ctx, sessionKey, sessData, day).Err()

	if err != nil {
		initEnv.Logger.Error("set redis err, err:", err)
		return nil, err
	}

	return gin.H{
		"sid": user.ID,
		"msg": errmsg.ErrMsg["LOGIN_BY_LOCAL_SUCCESS"],
	}, nil

}