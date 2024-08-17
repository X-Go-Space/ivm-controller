package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"gorm.io/gorm"
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
func LoginByHttp (loginData map[string]interface{}, ctx *gin.Context) (interface{}, error) {
	// HTTP登录
	userName := utils.ReadNestedData(loginData, "userName")
	password := utils.ReadNestedData(loginData, "password").(string)
	password, err := utils.Decrypt(password)
	if err != nil {
		initEnv.Logger.Error("when login decrypt password failed, err:", err)
		return nil, err
	}
	var user model.User
	result := initEnv.Db.Where("is_local = ? AND user_name = ?", "0", userName).First(&user)
	if result.Error != nil{
		initEnv.Logger.Error("when login query mysql failed, err:", err)
		return nil, err
	}
	if user.ID == "" {
		initEnv.Logger.Error("the user is not in the database")
		return nil, fmt.Errorf("the user is not in the database")
	}
	user.Password = password
	// 查到了数据库的用户信息
	// 需要根据认证服务器的配置，将用户信息给发送出去
	row:= initEnv.Db.Raw(utils.GET_AUTH_CONFIG_FROM_USER_AND_USER_DIRECTORY, user.ID).Row()
	if err != nil {
		initEnv.Logger.Error("get auth config failed,err:", err)
		return nil,err
	}
	var authConfigJson string
	if row != nil {
		row.Scan(&authConfigJson)
	}

	var authConfig []model.AuthConfig
	err = json.Unmarshal([]byte(authConfigJson), &authConfig)
	if err != nil {
		initEnv.Logger.Error("get auth config failed,err:", err)
		return nil,err
	}

	httpRes := utils.SendRequest(authConfig, &user)
	if !httpRes {
		initEnv.Logger.Error("http login failed")
		return nil, fmt.Errorf("http login failed")
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
		"msg": errmsg.ErrMsg["LOGIN_BY_HTTP_SUCCESS"],
	}, nil
}

func LoginByOauth2(loginData map[string]interface{}, ctx *gin.Context)(interface{}, error) {
	// 根据id找到对应的authConfig
	var oauth2AuthServer model.AuthServer
	code := utils.ReadNestedData(loginData,"code").(string)
	id := utils.ReadNestedData(loginData, "authId").(string)
	fmt.Println(code, id)
	result := initEnv.Db.Where("id = ?", id).First(&oauth2AuthServer)
	if result.Error != nil {
		initEnv.Logger.Error("get oauth2 auth server fail,err:", result.Error)
		return nil, result.Error
	}
	if oauth2AuthServer.Id == "" {
		initEnv.Logger.Error("the oauth2 auth server is empty")
		return nil, fmt.Errorf("the user is not in the database")
	}
	var authConfig []model.AuthConfig
	err := json.Unmarshal([]byte(oauth2AuthServer.AuthConfigJson), &authConfig)
	if err != nil {
		initEnv.Logger.Error("Unmarshal oauth2 auth server is failed,err:", err)
		return nil, err
	}

	// 拿到code，然后根据code从第三方去换信息
	var user model.User
	user.Code = code
	oauth2Res := utils.SendRequest(authConfig, &user)
	if !oauth2Res {
		initEnv.Logger.Error("oauth2 login failed")
		return nil, fmt.Errorf("oauth2 login failed")
	}
	// user 已经被绑定了，找一下该用户是否存在在该数据库里面
	result = initEnv.Db.Where("user_name = ?", user.UserName).First(&user)
	if result.Error !=nil &&result.Error == gorm.ErrRecordNotFound {
		initEnv.Logger.Error("oauth2 login failed, user is not in the database")
		return nil, fmt.Errorf("oauth2 login failed, user is not in the database")
	} else if result.Error !=nil {
		initEnv.Logger.Error("oauth2 login find user failed")
		return nil, fmt.Errorf("oauth2 login find user failed")
	}
	// 找到改用户了
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
		initEnv.Logger.Error("oauth2 marshal sess data failed,err:", err)
		return nil, err
	}

	err = initEnv.Redis.Set(ctx, sessionKey, sessData, day).Err()

	if err != nil {
		initEnv.Logger.Error("oauth2 set redis err, err:", err)
		return nil, err
	}

	return gin.H{
		"sid": user.ID,
		"msg": errmsg.ErrMsg["LOGIN_BY_OAUTH2_SUCCESS"],
	}, nil
}

func LoginByQrcode(data map[string]interface{},ctx *gin.Context) (interface{}, error) {
	QRcodeId := utils.ReadNestedData(data, "qrcodeId").(string)

	val, err := initEnv.Redis.Get(ctx, "SESS#ID#"+QRcodeId).Result()
	if err == redis.Nil {
		initEnv.Logger.Warn("LoginByQrcode qrcode is expired")
		return gin.H{
			"status":utils.QRcodeExpired,
		}, err
	} else if err != nil {
		return nil,err
	}

	userInfo, err := initEnv.Redis.Get(ctx, "USER#SESS#"+val).Result()
	if err == redis.Nil {
		initEnv.Logger.Warn("LoginByQrcode qrcode get sid fail")
		return gin.H{
			"status":utils.QRcodeExpired,
		}, err
	} else if err != nil {
		return nil,err
	}
	// 找到改用户了
	sessionKey := utils.GenerateSessId(QRcodeId)
	err = initEnv.Redis.Set(ctx, sessionKey, userInfo, day).Err()
	if err != nil {
		initEnv.Logger.Error("qrcode set redis err, err:", err)
		return nil,err
	}

	return gin.H{
		"sid": sessionKey,
		"msg": errmsg.ErrMsg["LOGIN_BY_QRCODE_SUCCESS"],
	}, nil
}