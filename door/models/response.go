package models

import (
	"github.com/gin-gonic/gin"
	"net/http"
)


// ResponseData 自定义返回的消息结构
type ResponseData struct {
	Code    int      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}



func ResponseErrorWithMsg(ctx *gin.Context, code int, errMsg string) {
	rd := &ResponseData{
		Code:    code,
		Message: errMsg,
		Data:    nil,
	}
	ctx.JSON(http.StatusOK, rd)
}

func ResponseSuccess(ctx *gin.Context, code int,data interface{}) {
	rd := &ResponseData{
		Code:    code,
		Message: "success",
		Data:    data,
	}
	ctx.JSON(http.StatusOK, rd)
}
