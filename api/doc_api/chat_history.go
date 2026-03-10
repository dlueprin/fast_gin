package doc_api

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/model"
	"fast_gin/utils/res"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type HistoryRequest struct {
	BeforeID uint `form:"before_id"`
}

func (d DocApi) ChatHistoryView(c *gin.Context) {
	cr := middleware.GetBind[HistoryRequest](c)
	userID := middleware.GetAuth(c).UserID
	var list []model.ChatHistoryModel
	const pageSize = 5

	db := global.DB.Where("user_id = ?", userID).Order("id desc").Limit(pageSize)

	if cr.BeforeID > 0 {
		db = db.Where("id < ?", cr.BeforeID)
	}

	//这里查找完之后会触发钩子函数，将那个结构体字段填满
	if err := db.Find(&list).Error; err != nil {
		logrus.Errorf("查询失败: %v", err)
		res.FailWithMsg("查询失败", c)
		return
	}

	res.OkWithData(list, c)
}
