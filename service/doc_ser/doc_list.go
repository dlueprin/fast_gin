package doc_ser

import (
	"fast_gin/model"
	"fast_gin/service/common"
)

func (s *DocService) GetDocList(page model.PageInfo) (list []model.DocumentModel, count int64, err error) {
	list, count, err = common.QueryList[model.DocumentModel](
		model.DocumentModel{},
		common.QueryOption{
			PageInfo: page,
			Likes:    []string{"title", "summary"}, //支持对标题和摘要进行模糊查询
		},
	)
	return
}
