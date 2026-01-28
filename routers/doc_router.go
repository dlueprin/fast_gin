package routers

import (
	"fast_gin/api"
	"fast_gin/middleware"
	"fast_gin/model"
	"github.com/gin-gonic/gin"
)

func DocRouter(g *gin.RouterGroup) {
	app := api.App.DocApi
	g.POST("doc/upload",
		middleware.LimitMiddleware(1),
		middleware.AuthMiddleware,
		app.DocUploadView, //gin会自动传上下文
	)
	g.GET("doc",
		middleware.BindQueryMiddleware[model.PageInfo],
		app.DocListView,
	)
}
