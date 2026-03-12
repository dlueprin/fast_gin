package routers

import (
	"fast_gin/api"
	"fast_gin/api/doc_api"
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
		middleware.AuthMiddleware,
		middleware.BindQueryMiddleware[model.PageInfo],
		app.DocListView,
	)
	g.DELETE("doc",
		middleware.AuthMiddleware,
		middleware.BindJsonMiddleware[doc_api.DocDeleteRequest],
		app.DocDeleteView,
	)
	g.POST("doc/search",
		middleware.LimitMiddleware(1),
		middleware.AuthMiddleware,
		middleware.BindJsonMiddleware[doc_api.DocSearchRequest],
		app.DocSearchView,
	)
	g.POST("doc/chat",
		middleware.LimitMiddleware(1),
		middleware.AuthMiddleware,
		middleware.BindJsonMiddleware[doc_api.DocChatRequest],
		app.DocChatView,
	)
	g.GET("doc/chat/history",
		middleware.LimitMiddleware(30),
		middleware.AuthMiddleware,
		middleware.BindQueryMiddleware[doc_api.HistoryRequest],
		app.ChatHistoryView,
	)
	g.GET("doc/:doc_id",
		middleware.LimitMiddleware(10),
		middleware.AuthMiddleware,
		app.DocDetailView,
	)
}
