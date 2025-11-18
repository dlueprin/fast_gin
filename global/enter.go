package global

import (
	"fast_gin/config"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const Version = "0.0.1"

var (
	//小写开头代表不可导出的，只能在当前包使用
	Config *config.Config
	DB     *gorm.DB
	Redis  *redis.Client
)
