package image_api

import (
	"fast_gin/global"
	"fast_gin/utils/find"
	"fast_gin/utils/md5"
	"fast_gin/utils/random"
	"fast_gin/utils/res"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// 白名单，允许上传的文件后缀
var whiteList = []string{
	".png",
	".jpg",
	".jpeg",
	".webp",
	".gif",
	".svg",
	".ico",
	".bmp",
}

func (ImageApi) ImageUploadView(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		logrus.Errorf("获取文件失败：%s", err)
		res.FailWithMsg("请选择文件", c)
		return
	}
	//现在就可以获取文件头部、文件名、文件大小还有读文件了
	if fileHeader.Size > global.Config.Image.Size*1024*1024 {
		res.FailWithMsg("文件过大!", c)
	}
	//后缀判断,可能有大写的，所以转小写
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !find.Find(whiteList, ext) {
		res.FailWithMsg("文件格式错误!", c)
	}
	//处理同名文件
	//1。通过哈希值判断是否是相同文件
	//2.通过给同名不同哈希的文件后面加上三个包含数字和大小写字母的方式区别开

	//接收并保存文件
	fp := path.Join("uploads", global.Config.Image.Path, fileHeader.Filename)
	//无条件for当while使用，判断第二次生成的文件是否也重名
	for {
		//stat获取路径下的文件信息
		_, err1 := os.Stat(fp)
		if os.IsNotExist(err1) {
			//说明文件不存在
			break
		}
		//计算上传的图片哈希值和本地的哈希并判断是否一样，一样就返回原来的地址
		uploadFile, _ := fileHeader.Open()
		oldFile, _ := os.Open(fp)
		uploadFileHash := md5.MD5WithFile(uploadFile)
		oldFileHash := md5.MD5WithFile(oldFile)

		if uploadFileHash == oldFileHash {
			res.Ok("/"+fp, "upload success", c)
			return
		}
		//生成随机后缀
		fileNameNotExt := strings.TrimSuffix(fileHeader.Filename, ext)
		NewFileName := fmt.Sprintf("%s_%s%s", fileNameNotExt, random.RandStr(3), ext)
		fp = path.Join("uploads", global.Config.Image.Path, NewFileName)
	}

	c.SaveUploadedFile(fileHeader, fp)

	res.Ok("/"+fp, "upload success", c)
}
