package doc_ser

import (
	"fast_gin/global"
	"fast_gin/model"
	"regexp"
)

//ALTER TABLE document_models ADD FULLTEXT INDEX ft_title_summary (title, summary) WITH PARSER ngram;
//ALTER TABLE document_models ADD FULLTEXT INDEX ft_content (content) WITH PARSER ngram;

func (s *DocService) GetDocList(page model.PageInfo) (list []model.DocumentModel, count int64, err error) {
	//1、初始化数据库对象，方便后面构建各种条件
	db := global.DB.Model(&model.DocumentModel{})

	if global.Config.System.Mode == "debug" {
		db = db.Debug()
	}

	//默认加载用户表
	db = db.Preload("User")

	//2、根据是否有关键词决定用什么搜索方法
	if page.Key != "" {
		//算分数
		db = db.Select(
			"*,(MATCH(title,summary) AGAINST(? IN NATURAL LANGUAGE MODE) + MATCH(content) AGAINST(? IN NATURAL LANGUAGE MODE)) AS score",
			page.Key, page.Key,
		)
		db = db.Where( //or链接两个ngram索引，都是相同关键词查文章内容和标题和摘要内容
			"MATCH(title,summary)	AGAINST(? IN NATURAL LANGUAGE MODE) OR MATCH(content) AGAINST(? IN NATURAL LANGUAGE MODE)",
			page.Key, page.Key,
		)
	}
	//else {
	//	list, count, err = common.QueryList[model.DocumentModel](
	//		model.DocumentModel{},
	//		common.QueryOption{
	//			PageInfo: page,
	//			//Likes:    []string{"title", "summary"}, //支持对标题和摘要进行模糊查询
	//		},
	//	)
	//	return
	//}

	//3、统计总数
	err = db.Count(&count).Error
	if err != nil || count == 0 {
		return nil, 0, err
	}

	if page.Order != "" {
		db = db.Order(page.Order)
	} else if page.Key != "" {
		db = db.Order("score desc")
	} else {
		db = db.Order("id desc")
	}

	//4、分页
	offset := (page.Page - 1) * page.Limit
	err = db.Limit(page.Limit).Offset(offset).Find(&list).Error

	//5、高亮搜索词
	if page.Key != "" {
		re := regexp.MustCompile("(?i)(" + regexp.QuoteMeta(page.Key) + ")")
		for i := range list {
			list[i].Summary = re.ReplaceAllString(list[i].Summary, "<em class='search-key'>$1</em>")
		}
	}
	return
}
