package global

import (
	"fast_gin/config"
	"gorm.io/gorm"
)

var (
	//小写开头代表不可导出的，只能在当前包使用
	Config *config.Config
	DB     *gorm.DB
)
