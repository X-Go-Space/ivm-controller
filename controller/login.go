package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"ivm-controller/initEnv"
	"ivm-controller/model"
	"ivm-controller/service"
	"ivm-controller/utils"
	"ivm-controller/utils/errmsg"
	"time"
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
			initEnv.Logger.Error("login by local failed,err:", err)
			utils.Err(errmsg.ErrMsg["LOGIN_BY_LOCAL_FAILED"], errmsg.LOGIN_BY_LOCAL_FAILED, ctx)
			break
		}
		utils.OK(data, ctx)
	case "http":
		data, err := service.LoginByHttp(loginData, ctx)
		if err !=nil {
			initEnv.Logger.Error("login by http failed,err:", err)
			utils.Err(errmsg.ErrMsg["LOGIN_BY_HTTP_FAIL"], errmsg.LOGIN_BY_HTTP_FAIL, ctx)
			break
		}
		utils.OK(data, ctx)
	case "oauth2":
		data, err := service.LoginByOauth2(loginData, ctx)
		if err !=nil {
			initEnv.Logger.Error("login by oauth2 failed,err:", err)
			utils.Err(errmsg.ErrMsg["LOGIN_BY_OAUTH2_FAIL"], errmsg.LOGIN_BY_OAUTH2_FAIL, ctx)
			break
		}
		utils.OK(data, ctx)
	case "qrcode":
		data, err := service.LoginByQrcode(loginData, ctx)
		if err !=nil {
			initEnv.Logger.Error("login by local bind json failed,err:", err)
			utils.Err(errmsg.ErrMsg["LOGIN_BY_QRCODE_FAIL"], errmsg.LOGIN_BY_QRCODE_FAIL, ctx)
			break
		}
		utils.OK(data, ctx)
	default:
		initEnv.Logger.Error("login failed, unknow login type")
		utils.Err(errmsg.ErrMsg["UNKNOW_LOGIN_TYPE"], errmsg.UNKNOW_LOGIN_TYPE, ctx)
	}
}

func AuthConfig(ctx *gin.Context) {
	// 获取认证服务器配置信息
	var authServers []model.AuthServer
	result := initEnv.Db.Find(&authServers)
	if result.Error != nil {
		initEnv.Logger.Error("get auth servers from mysql failed, err:", result.Error)
		utils.Err(errmsg.ErrMsg["GET_AUTH_CONFIG_FAIL"], errmsg.GET_AUTH_CONFIG_FAIL, ctx)
	}
	// oauth2需要下发获取code地址
	authConfigData := make([]gin.H, 0)
	for _, authServer := range authServers {
		var data  = make(gin.H)
		data["ID"] = authServer.Id
		data["authType"] = authServer.AuthType
		if authServer.AuthType == "authOauth2" {
			data ["getCodeUrl"] = authServer.GetCode
		}
		authConfigData = append(authConfigData, data)
	}
	utils.OK(authConfigData, ctx)
}

func GenerateQrCode(ctx *gin.Context) {
	QRcodeID := utils.RandID()

	// 将该ID设置到sess里面，初始化二维码为初始化状态
	err := initEnv.Redis.Set(ctx, QRcodeID, utils.QRcodeInit, 3*time.Minute).Err()
	if err != nil {
		initEnv.Logger.Error("generate qrcode id fail,err:", err)
		utils.Err(errmsg.ErrMsg["GENERATE_QRCODE_FAIL"], errmsg.GENERATE_QRCODE_FAIL, ctx)
		return
	}
	utils.OK(gin.H{
		"QRcodeID":QRcodeID,
	},ctx)
}

func GetQrcodeStatus(ctx *gin.Context) {
	QRcodeID := ctx.Query("qrcodeId")
	val, err := initEnv.Redis.Get(ctx, QRcodeID).Result()
	if err == redis.Nil {
		initEnv.Logger.Warn("qrcode is expired")
		utils.OK(gin.H{
			"status":utils.QRcodeExpired,
		},ctx)
		return
	}
	fmt.Println(val)
	if val == utils.QRcodeInit {
		initEnv.Logger.Warn("qrcode is init")
		utils.OK(gin.H{
			"status":utils.QRcodeInit,
		},ctx)
		return
	}

	if val == utils.QRcodeScan {
		initEnv.Logger.Warn("qrcode is scanned")
		utils.OK(gin.H{
			"status":utils.QRcodeScan,
		},ctx)
		return
	}

	// 如果是已确认，让前端携带qrCodeID进行认证上线
	if val == utils.QRcodeConfirm {
		utils.OK(gin.H{
			"status":utils.QRcodeConfirm,
		},ctx)
		return
	}
	initEnv.Logger.Error("the qrcode status is not right")
	utils.Err(errmsg.ErrMsg["QRCODE_STATUS_ERROR"], errmsg.QRCODE_STATUS_ERROR, ctx)
}

// MobileQrcodeLogin 移动端扫描状态
func MobileQrcodeLogin(ctx *gin.Context) {
	var data map[string]interface{}
	err := ctx.BindJSON(&data)
	if err != nil {
		initEnv.Logger.Error("MobileQrcodeLogin bind json failed, err:", err)
		utils.Err(errmsg.ErrMsg["MOBILE_LOGIN_SCAN_FAIL"], errmsg.MOBILE_LOGIN_SCAN_FAIL, ctx)
		return
	}
	QRcodeId := utils.ReadNestedData(data, "qrcodeId").(string)
	err = initEnv.Redis.Set(ctx, QRcodeId, utils.QRcodeScan, 3*time.Minute).Err()
	if err != nil {
		initEnv.Logger.Error("scan qrcode id fail,err:", err)
		utils.Err(errmsg.ErrMsg["MOBILE_LOGIN_SCAN_FAIL"], errmsg.MOBILE_LOGIN_SCAN_FAIL, ctx)
		return
	}
	utils.OK(gin.H{
		"message": errmsg.ErrMsg["MOBILE_LOGIN_SCAN_SUCCESS"],
	}, ctx)
}

func MobileQrcodeConfirm(ctx *gin.Context) {
	var data map[string]interface{}
	err := ctx.BindJSON(&data)
	if err != nil {
		initEnv.Logger.Error("MobileQrcodeLogin bind json failed, err:", err)
		utils.Err(errmsg.ErrMsg["MOBILE_LOGIN_CONFIRM_FAIL"], errmsg.MOBILE_LOGIN_CONFIRM_FAIL, ctx)
		return
	}
	QRcodeId := utils.ReadNestedData(data, "qrcodeId").(string)
	err = initEnv.Redis.Set(ctx, QRcodeId, utils.QRcodeConfirm, 3*time.Minute).Err()
	if err != nil {
		initEnv.Logger.Error("qrcode confirm fail,err:", err)
		utils.Err(errmsg.ErrMsg["MOBILE_LOGIN_CONFIRM_FAIL"], errmsg.MOBILE_LOGIN_CONFIRM_FAIL, ctx)
		return
	}
	sid := ctx.Request.Header.Get("sid")
	err = initEnv.Redis.Set(ctx, "SESS#ID#"+QRcodeId, sid, 3*time.Minute).Err()
	if err != nil {
		initEnv.Logger.Error("qrcode confirm fail get sid,err:", err)
		utils.Err(errmsg.ErrMsg["MOBILE_LOGIN_CONFIRM_FAIL"], errmsg.MOBILE_LOGIN_CONFIRM_FAIL, ctx)
		return
	}
	utils.OK(gin.H{
		"message": errmsg.ErrMsg["MOBILE_LOGIN_CONFIRM_SUCCESS"],
	}, ctx)
}
