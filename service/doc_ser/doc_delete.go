package doc_ser

import (
	"fast_gin/global"
	"fast_gin/model"
	"github.com/sirupsen/logrus"
	"os"
)

func (s *DocService) DeleteDoc(idList []uint) (err error) {
	var docList []model.DocumentModel
	//1、找到对应的数据存到准备好的结构体，同时为后面删除找到对应的id（主键）防止误删，是常用的方法
	global.DB.Find(&docList, idList)
	//2、删除本地文件并记录删除失败的文件位置
	for _, doc := range docList {
		err = os.Remove(doc.Path)
		if err != nil {
			logrus.Errorf("删除文件失败%v,路径:%s", err, doc.Path)
		}
	}
	//3、删除数据库中数据。这里使用了gorm的主键删除，先传模型后传主键值（可以是列表）
	return global.DB.Delete(&model.DocumentModel{}, idList).Error
}
