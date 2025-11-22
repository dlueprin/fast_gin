package res

import (
	"fast_gin/utils/validate"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int    `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

func Ok(data any, msg string, c *gin.Context) {
	c.JSON(200, Response{
		Code: 200,
		Data: data,
		Msg:  msg,
	})
}

// 一般用于列表获取信息的返回
func OkWithData(data any, c *gin.Context) {
	Ok(data, "success", c)
}
func OkWithMsg(msg string, c *gin.Context) {
	Ok(gin.H{}, msg, c)
}
func Fail(code int, msg string, c *gin.Context) {
	c.JSON(200, Response{
		Code: code,
		Data: gin.H{},
		Msg:  msg,
	})
}
func FailWithMsg(msg string, c *gin.Context) {
	Fail(7, msg, c)
}

// FailWithError 说是参数绑定很好用
func FailWithError(err error, c *gin.Context) {
	msg := validate.ValidateError(err)
	Fail(7, msg, c)
}
