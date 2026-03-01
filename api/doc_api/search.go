package doc_api

import (
	"fast_gin/middleware"
	"fast_gin/utils/res"
	"github.com/gin-gonic/gin"
)

type DocSearchRequest struct {
	Query string `json:"query" binding:"required"`
}

func (d DocApi) DocSearchView(c *gin.Context) {
	//1、绑定关键词
	cr := middleware.GetBind[DocSearchRequest](c)
	//2、获取用户id
	userID := middleware.GetAuth(c).UserID
	//3、调用业务逻辑
	results, err := docService.Search(userID, cr.Query, 5)
	if err != nil {
		res.FailWithMsg("搜索失败:"+err.Error(), c)
		return
	}
	res.OkWithData(results, c)
}
