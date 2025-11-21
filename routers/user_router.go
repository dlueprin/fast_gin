package routers

import (
	"fast_gin/api"
	"github.com/gin-gonic/gin"
)

func UserRouter(g *gin.RouterGroup) {
	app := api.App.UserApi
	g.POST("users/login", app.LoginView) //经由路由组管理，api开头的来到这里，现在路径变成 api/users/login
}
