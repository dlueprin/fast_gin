package model

import "time"

// Model model模板，让其他model继承，可以少写一点
type Model struct {
	ID        int       `gorm:"primaryKey" json:"id"` //json统一使用小驼峰格式
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// PageInfo 分页参数，BindQueryMiddleware会走到这里来并自动存放对应的参数
type PageInfo struct {
	Page  int    `form:"page"`  //页码
	Limit int    `form:"limit"` //每页数量限制
	Key   string `form:"key"`   //搜索关键字
	Order string `form:"order"`
}
