package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func OK(data interface{}, ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": data,
	})
}

func Err(errMsg string, errCode int, ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": errCode,
		"err_msg": errMsg,
	})
}