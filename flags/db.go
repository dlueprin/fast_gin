package flags

import (
	"fast_gin/global"
	"fast_gin/model"
	"github.com/sirupsen/logrus"
)

// MigrateDB 迁移数据库,属于运维的内容，所以放在flags包中
func MigrateDB() {
	err := global.DB.AutoMigrate(
		&model.UserModel{},
	)
	if err != nil {
		logrus.Errorf("表结构迁移失败：%s", err)
		return
	}
	logrus.Infof("表结构迁移成功")
}
