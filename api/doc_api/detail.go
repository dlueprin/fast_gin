package doc_api

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/model"
	"fast_gin/utils/res"
	"fmt"
	"github.com/gin-gonic/gin"
)

func (d DocApi) DocDetailView(c *gin.Context) {
	docID := c.Param("doc_id")
	userID := middleware.GetAuth(c).UserID

	var doc model.DocumentModel
	if err := global.DB.Where("id = ? AND user_id = ?", docID, userID).First(&doc).Error; err != nil {
		res.FailWithMsg("文档不存在或无权访问", c)
		return
	}

	fileUrl := fmt.Sprintf("http://localhost:8080/%s", doc.Path)

	res.OkWithData(gin.H{
		"title":     doc.Title,
		"file_type": doc.FileType,
		"content":   doc.Content,
		"file_url":  fileUrl,
	}, c)
}
