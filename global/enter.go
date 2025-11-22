package global

import (
	"fast_gin/config"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const Version = "0.0.5"

// 这些都是基础设施级别的全局变量，所以单独建包存放他们
var (
	//小写开头代表不可导出的，只能在当前包使用
	Config *config.Config
	DB     *gorm.DB
	Redis  *redis.Client
)
