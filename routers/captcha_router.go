package routers

import (
	"fast_gin/api"
	"fast_gin/middleware"
	"github.com/gin-gonic/gin"
)

func CaptchaRouter(g *gin.RouterGroup) {
	app := api.App.CaptchaApi
	g.GET("captcha/generate", middleware.LimitMiddleware(10), app.GenerateView)
}
