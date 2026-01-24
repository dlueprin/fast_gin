package doc_ser

import (
	"fast_gin/global"
	"fast_gin/model"
)

// 将文档内容传到数据库
func (s *DocService) UploadDoc(doc *model.DocumentModel) error {
	//然后这里用指针doc对象，存入数据库生成的id和创建时间什么的会自动填到doc对象里面
	return global.DB.Create(doc).Error
}
