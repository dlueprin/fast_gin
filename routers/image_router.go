package routers

import (
	"fast_gin/api"
	"fast_gin/middleware"
	"github.com/gin-gonic/gin"
)

func ImageRouter(g *gin.RouterGroup) {
	//获取实例以获得方法
	app := api.App.ImageApi
	g.POST("images/upload", middleware.LimitMiddleware(1), middleware.AuthMiddleware, app.ImageUploadView)
}
