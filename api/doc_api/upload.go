package doc_api

import (
	"fast_gin/global"
	"fast_gin/middleware"
	"fast_gin/model"
	"fast_gin/service/doc_ser"
	"fast_gin/utils/res"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"strings"
	"time"
)

func (d DocApi) DocUploadView(c *gin.Context) {
	//用gin的获取表单文件拿到文件控制权
	fileHeader, err := c.FormFile("file")
	if err != nil {
		logrus.Errorf("获取文件失败:%s", err)
		res.FailWithMsg("请选择文件", c)
	}
	//校验后缀
	ext := strings.ToLower(path.Ext(fileHeader.Filename))
	if ext != ".txt" && ext != ".md" && ext != ".pdf" && ext != ".docx" {
		logrus.Errorf("上传了非文档格式的文件：%s", ext)
		res.FailWithMsg("目前只支持pdf、docx、txt和md格式的文件", c)
	}
	basePath := path.Join("uploads", global.Config.Doc.Path)
	if _, err = os.Stat(basePath); os.IsExist(err) {
		os.MkdirAll(basePath, os.ModePerm)
	}
	fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), fileHeader.Filename)
	filePath := path.Join(basePath, fileName)

	if err = c.SaveUploadedFile(fileHeader, filePath); err != nil {
		logrus.Errorf("保存文件失败:%s", err)
		res.FailWithMsg("文件上传失败", c)
		return
	}
	claims := middleware.GetAuth(c)
	userID := claims.UserID
	var docService doc_ser.DocService
	docRecord := model.DocumentModel{
		Title:    fileHeader.Filename,
		Path:     filePath,
		FileType: ext,
		UserID:   userID,
	}
	if err = docService.UploadDoc(&docRecord); err != nil {
		logrus.Errorf("保存文件信息到数据库失败:%s", err)
		res.FailWithMsg("文件上传失败", c)
		return
	}

	//ai异步处理部分
	go docService.AsyncAnalyze(&docRecord, ext)

	res.OkWithData(gin.H{
		"id":    docRecord.ID,
		"title": docRecord.Title,
		"msg":   "文件上传成功，ai摘要正在后台生成中",
	}, c)

}
