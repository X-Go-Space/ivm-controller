package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"ivm-controller/initEnv"
	"ivm-controller/utils/errmsg"
)

// IsLoginMiddleWare 登录验证中间件
func IsLoginMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 模拟登录状态，这里可以替换为实际的登录验证逻辑
		//var user model.User
		//_ = c.ShouldBindJSON(&user)
		userName:=c.PostForm("userName")
		// 在这里查询redis有没有过期，每次请求对应的接口都要携带对应的sid才可以
		isLoggedIn := true

		if isLoggedIn {
			// 如果已登录，继续执行下一个处理程序
			c.Next()
		} else {
			initEnv.Logger.WithFields(logrus.Fields{
				"username": userName,
			}).Error("当前用户未登录成功，请重新登录")
			c.Abort()
			c.JSON(errmsg.NOT_LOGIN, gin.H{
				"message": errmsg.GetErrMsg(errmsg.NOT_LOGIN),
			})
		}
	}
}