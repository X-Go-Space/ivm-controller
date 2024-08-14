package middleware

import (
	"github.com/gin-gonic/gin"
	"ivm-controller/initEnv"
	"ivm-controller/utils/errmsg"
	"net/http"
)

// GlobalErrorInterceptor 是一个全局中间件，捕获所有 panic
func GlobalErrorInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			var errFlag any =nil
			if err := recover(); err != errFlag {
				initEnv.Logger.Error("panic err msg:", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": errmsg.GetErrMsg(errmsg.ERROR)})
			}
		}()
		c.Next()
	}
}