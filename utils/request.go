package utils

import (
	"bytes"
	"fmt"
	"github.com/goccy/go-json"
	"io/ioutil"
	"ivm-controller/initEnv"
	"ivm-controller/model"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)
func getFieldValue(obj interface{}, fieldName string) interface{} {
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	field := val.FieldByName(fieldName)
	if !field.IsValid() {
		return nil
	}

	return field.Interface()
}

func structHasField(s interface{}, fieldName string) bool {
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr && !val.IsNil() {
		val = val.Elem()
	} else {
		return false
	}

	_, found := val.Type().FieldByName(fieldName)
	return found
}

func setFieldValue(obj interface{}, fieldName string, value interface{}) {
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	field := val.FieldByName(fieldName)
	if !field.IsValid() {
		fmt.Printf("Field %s not found in struct\n", fieldName)
		return
	}

	if field.CanSet() {
		field.Set(reflect.ValueOf(value))
	} else {
		fmt.Printf("Cannot set value for field %s\n", fieldName)
	}
}
// GetValue 对于KEY是以@开头的那么就要获取对应的值
func GetValue(key string, user *model.User, midData map[string]interface{}) string {
	var res interface{}
	if user != nil {
		res = getFieldValue(user, key)
		if res != nil {
			return fmt.Sprint(res)
		}
	}

	res = ReadNestedData(midData, key)
	if res != nil {
		return fmt.Sprint(res)
	}
	return ""
}
func JudgeRespSuccess(resp map[string]interface{}, condition model.SuccessCondition) bool {
	respValue := fmt.Sprint(ReadNestedData(resp, condition.ResponseFiled))
	switch condition.ResponseCondition {
	case "equal":
		if respValue == condition.ResponseResult {
			return true
		}
	case "noEqual":
		if respValue != condition.ResponseResult {
			return true
		}
	case "notNull":
		if respValue != "" {
			return true
		}
	case "isNull":
		if respValue == "" {
			return true
		}
	default:
		initEnv.Logger.Error("the ResponseCondition is not right")
		return false
	}
	return false
}
func RespMapResult(resp map[string]interface{}, successDataMap [][2]string, user *model.User) error {
	for _, dataMap := range successDataMap {
		mapKey, userKey := dataMap[0], dataMap[1]
		if !structHasField(user, userKey){
			initEnv.Logger.Error("mapdata user key is not exist")
			return fmt.Errorf("mapdata user key is not exist")
		}
		if strings.HasPrefix(mapKey, "@") {
			value := GetValue(mapKey[1:], nil, resp)
			if value == "" {
				initEnv.Logger.Error("get RespMapResult data failed")
				return fmt.Errorf("get RespMapResult data failed")
			}
			setFieldValue(user, userKey, value)
		}
	}
	return nil
}
func getHttp(authConfig model.AuthConfig, user *model.User, midData map[string]interface{}) (error, string) {
     reqUrl, err := url.Parse(authConfig.BaseUrl)
	 if err != nil {
		 initEnv.Logger.Error("get parse base url fail，err:", err)
		 return err, ""
	 }
	queryParams := reqUrl.Query()
	 for _, paramArray := range authConfig.Params{
		 key, value := paramArray[0], paramArray[1]
		 if strings.HasPrefix(value, "@") {
			 value = GetValue(value[1:], user, midData)
			 if value == "" {
				 initEnv.Logger.Error("get param get data failed")
				 return fmt.Errorf("get param get data failed"), ""
			 }
		 }
		 queryParams.Set(key, value)
	 }
	reqUrl.RawQuery = queryParams.Encode()
	// 创建一个新的请求
	req, err := http.NewRequest("GET", reqUrl.String(), bytes.NewBuffer([]byte(authConfig.Body)))
	if err != nil {
		initEnv.Logger.Error("create get req failed, err:", err)
		return err, ""
	}

	for _, headerArray := range authConfig.Headers{
		key, value := headerArray[0], headerArray[1]
		if strings.HasPrefix(value, "@") {
			value = GetValue(value[1:], user, midData)
			if value == "" {
				initEnv.Logger.Error("get header get data failed")
				return fmt.Errorf("get header get data failed"), ""
			}
			req.Header.Set(key, value)
		}
	}
	// 发起请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		initEnv.Logger.Error("send get request fail, err:", err)
		return err, ""
	}
	// 读取响应体内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		initEnv.Logger.Error("read data fail, err:", err)
		return err, ""
	}

	var respStruct map[string]interface{}
	err = json.Unmarshal(body, &respStruct)
	if err != nil {
		initEnv.Logger.Error("unmarshal data fail, err:", err)
		return err, ""
	}

	if !JudgeRespSuccess(respStruct, authConfig.SuccessCondition) {
		initEnv.Logger.Error("judge resp success fail")
		return err, ""
	}

	// w判断结束，如果登录成功，需要将返回的值进行映射
	err = RespMapResult(respStruct, authConfig.SuccessDataMap, user)
	if err != nil {
		initEnv.Logger.Error("resp map result fail")
		return err, ""
	}

	bodyData := string(body)
	return nil, bodyData
}

func postHttp(authConfig model.AuthConfig, user *model.User, midData map[string]interface{}) (error, string) {
	reqUrl, err := url.Parse(authConfig.BaseUrl)
	if err != nil {
		initEnv.Logger.Error("parse base url failed, err:", err)
		return err, ""
	}

	queryParams := reqUrl.Query()
	for _, paramArray := range authConfig.Params {
		key, value := paramArray[0], paramArray[1]
		if strings.HasPrefix(value, "@") {
			value = GetValue(value[1:], user, midData)
			if value == "" {
				initEnv.Logger.Error("post param get data failed")
				return fmt.Errorf("post param get data failed"), ""
			}
		}
		queryParams.Set(key, value)
	}
	reqUrl.RawQuery = queryParams.Encode()

	// 创建一个新的请求
	reqBody := []byte(authConfig.Body)
	req, err := http.NewRequest("POST", reqUrl.String(), bytes.NewBuffer(reqBody))
	if err != nil {
		initEnv.Logger.Error("create post request failed, err:", err)
		return err, ""
	}

	for _, headerArray := range authConfig.Headers {
		key, value := headerArray[0], headerArray[1]
		if strings.HasPrefix(value, "@") {
			value = GetValue(value[1:], user, midData)
			if value == "" {
				initEnv.Logger.Error("post header get data failed")
				return fmt.Errorf("post header get data failed"), ""
			}
		}
		req.Header.Set(key, value)
	}

	// 发起请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		initEnv.Logger.Error("send post request failed, err:", err)
		return err, ""
	}

	// 读取响应体内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		initEnv.Logger.Error("read data failed, err:", err)
		return err, ""
	}

	var respStruct map[string]interface{}
	err = json.Unmarshal(body, &respStruct)
	if err != nil {
		initEnv.Logger.Error("unmarshal data failed, err:", err)
		return err, ""
	}

	if !JudgeRespSuccess(respStruct, authConfig.SuccessCondition) {
		initEnv.Logger.Error("judge response success failed")
		return err, ""
	}
	// 判断结束，如果登录成功，需要将返回的值进行映射
	err = RespMapResult(respStruct, authConfig.SuccessDataMap, user)
	if err != nil {
		initEnv.Logger.Error("post resp map result fail")
		return err, ""
	}
	bodyData := string(body)
	return nil, bodyData
}

func SendRequest(authConfigs []model.AuthConfig, user *model.User) bool {
	// 根据传来的用户信息和请求配置进行发送数据
	var midData map[string]interface{}
	for _, authConfig := range authConfigs {
		if authConfig.RequestType == "GET" {
			err,result:= getHttp(authConfig, user, midData)
			if err != nil {
				initEnv.Logger.Error("send get http failed, err:", err)
				return false
			}
			err = json.Unmarshal([]byte(result), &midData)
			if err != nil {
				initEnv.Logger.Error("Unmarshal get http result failed, err:", err)
				return false
			}
		} else if authConfig.RequestType == "POST" {
			err,result:= postHttp(authConfig, user, midData)
			if err != nil {
				initEnv.Logger.Error("send post http failed, err:", err)
				return false
			}
			err = json.Unmarshal([]byte(result), &midData)
			if err != nil {
				initEnv.Logger.Error("Unmarshal post http result failed, err:", err)
				return false
			}
		}
	}
	return true
}
