package routers

import (
	"fast_gin/api"
	"fast_gin/middleware"
	"github.com/gin-gonic/gin"
)

func UserRouter(g *gin.RouterGroup) {
	app := api.App.UserApi
	//加一个1秒限制一次的限流
	g.POST("users/login", middleware.LimitMiddleware(1), app.LoginView) //经由路由组管理，api开头的来到这里，现在路径变成 api/users/login
	g.GET("users", middleware.LimitMiddleware(10), middleware.AuthMiddleware, app.UserListView)
}
