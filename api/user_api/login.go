package user_api

import (
	"fast_gin/utils/res"
	"github.com/gin-gonic/gin"
)

func (UserApi) LoginView(c *gin.Context) {
	res.OkWithMsg("登录页面", c)
	return
}
