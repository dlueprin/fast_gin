package common

import (
	"fast_gin/global"
	"fast_gin/model"
	"fmt"
	"gorm.io/gorm"
)

type QueryOption struct {
	model.PageInfo          //匿名字段/嵌入字段，相当于放了指定结构体的字段在里面，可以直接引用其中的字段，也叫字段提升field promotion
	Likes          []string //需要模糊查询的字段
	Where          *gorm.DB //高级查询条件
	Preloads       []string //预加载，外键关联查询时可能会用到
	Debug          bool     //gorm的显示sql语句的调试模式
}

func QueryList[T any](model any, option QueryOption) (list []T, count int64, err error) {
	list = make([]T, 0)
	query := global.DB.Where(model) //这里填入的model可以根据模型对象自动构建合理的空查询，初始化查询构建器,之前填的“”也是这个道理
	//模糊匹配，真是绕了好大一圈啊
	if option.Key != "" {
		if len(option.Likes) != 0 {
			likeQuery := global.DB.Where("")
			for _, column := range option.Likes {
				//用or连接所有模糊查询语句
				likeQuery.Or(
					fmt.Sprintf("%s like ?", column),
					fmt.Sprintf("%%%s%%", option.Key),
				)
			}
			query.Where(likeQuery)
		}
	}
	//预加载
	for _, preload := range option.Preloads {
		query = query.Preload(preload)
	}
	//分页查询
	if option.Page <= 0 {
		option.Page = 1
	}
	if option.Limit <= 0 {
		option.Limit = -1 //-1代表不分页
	}
	offset := (option.Page - 1) * option.Limit
	if option.Order == "" {
		option.Order = "created_at desc" //默认按创建时间倒序
	}
	//这里空查询是当作入口，节省代码的同时还可以控制是否debug，因为debug必须写在前面，还有可以后续添加其他控制
	db := global.DB.Where("")
	if option.Debug {
		db = db.Debug()
	}
	db.Where(query).Limit(option.Limit).Offset(offset).Order(option.Order).Find(&list)
	db.Model(model).Where(query).Count(&count) //注意统计查询次数是要指定模型的
	return
}
