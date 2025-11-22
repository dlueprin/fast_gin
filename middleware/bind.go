package middleware

import (
	"fast_gin/utils/res"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func BindJsonMiddleware[T any](c *gin.Context) {
	var cr T
	//should同时完成了校验和绑定
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		logrus.Errorf("参数校验失败：%s", err)
		res.FailWithError(err, c)
		c.Abort()
		return
	}
	//set将参数绑定好的request结构体（包含用户名和密码）放入上下文的
	c.Set("request", cr)
	logrus.Debugf("参数绑定成功：%v", cr)
	return
}
func BindQueryMiddleware[T any](c *gin.Context) {
	var cr T
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		logrus.Errorf("参数校验失败：%s", err)
		res.FailWithError(err, c)
		c.Abort()
		return
	}
	c.Set("request", cr)
	logrus.Debugf("参数绑定成功：%v", cr)
	return
}

func BindUriMiddleware[T any](c *gin.Context) {
	var cr T
	err := c.ShouldBindUri(&cr)
	if err != nil {
		logrus.Errorf("参数校验失败：%s", err)
		res.FailWithError(err, c)
		c.Abort()
		return
	}
	c.Set("request", cr)
	logrus.Debugf("参数绑定成功：%v", cr)
	return
}

func GetBind[T any](c *gin.Context) (cr T) {
	//如果不存就会panic，因为设计的时候规定要先经过绑定中间件，所以这里如果没有request，必定是前面出问题了
	return c.MustGet("request").(T)
}
