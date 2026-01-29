package doc_api

import (
	"fast_gin/middleware"
	"fast_gin/model"
	"fast_gin/utils/res"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (d DocApi) DocListView(c *gin.Context) {
	//1、获取中间件绑定好的分页参数
	//_page, _ := c.Get("request")//这是自己写的中间件放进上下文的东西，现在在视图这里取出来
	//page := _page.(model.PageInfo)

	page := middleware.GetBind[model.PageInfo](c) //使用封装好的getBind方法

	//2、调用服务层对应服务
	list, count, err := docService.GetDocList(page)
	if err != nil {
		logrus.Errorf("获取文档列表失败%v", err)
		res.FailWithMsg("获取文档列表失败", c)
		return
	}

	res.OkWithList(list, count, c)
}
