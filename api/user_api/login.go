package user_api

import (
	"fast_gin/middleware"
	"fast_gin/utils/res"
	"fmt"
	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required" label:"用户名"`
	Password string `json:"password" binding:"required" label:"密码"`
}

func (UserApi) LoginView(c *gin.Context) {
	//var cr LoginRequest
	////参数绑定就是获取请求中各种格式的参数到自己定义的结构体里的
	//if err := c.ShouldBindJSON(&cr); err != nil {
	//	msg := validate.ValidateError(err)
	//	logrus.Errorf("参数校验失败：%s", msg)
	//	res.FailWithError(err, c)
	//	return
	//}
	cr := middleware.GetBind[LoginRequest](c)
	fmt.Println(cr.Username)
	res.OkWithMsg("登录页面", c)
	return
}
