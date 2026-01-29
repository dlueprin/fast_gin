package doc_api

import (
	"fast_gin/middleware"
	"fast_gin/utils/res"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type DocDeleteRequest struct {
	IDList []uint `json:"id_list" binding:"required" label:"文档id列表"`
}

func (d DocApi) DocDeleteView(c *gin.Context) {
	//1、中间件获取绑定参数
	cr := middleware.GetBind[DocDeleteRequest](c)
	//2、调用业务逻辑
	err := docService.DeleteDoc(cr.IDList)
	if err != nil {
		logrus.Errorf("删除文档失败：%v", err)
		res.FailWithMsg("删除文档失败", c)
		return
	}
	res.FailWithMsg("删除文档成功", c)
}
