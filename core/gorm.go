package core

import (
	"fast_gin/global"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

// go get gorm.io/gorm
// go get gorm.io/driver/mysql
// go get gorm.io/driver/postgres
// go get github.com/glebarez/sqlite

func InitGorm() (db *gorm.DB) { //使用简单工厂模式：先用switch处理不同的数据库模式，后面统一用gorm，open，调用方无需关心处理逻辑，对外只暴露一个接口
	cfg := global.Config.DB   //通过全局变量获取配置项
	var dialector = cfg.Dsn() //数据库方言适配器，单词本身是方言的意思
	if dialector == nil {
		return
	}
	//switch cfg.Mode {
	//case config.DBMysqlMode:
	//	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
	//		cfg.User,
	//		cfg.Password,
	//		cfg.Host,
	//		cfg.Port,
	//		cfg.DBName,
	//	)
	//	dialector = mysql.Open(dsn) //驱动级别的open，用于创建数据库方言适配器
	//case config.DBPgsqlMode:
	//	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
	//		cfg.Host,
	//		cfg.User,
	//		cfg.Password,
	//		cfg.DBName,
	//		cfg.Port,
	//	)
	//	dialector = postgres.Open(dsn)
	//case config.DBSqliteMode:
	//	dialector = sqlite.Open(cfg.DBName)
	//default:
	//	logrus.Warnf("未配置数据库连接")
	//	return
	//}
	db, err := gorm.Open(dialector, &gorm.Config{ //数据库级别的open，根据前面的适配器真正建立连接
		DisableForeignKeyConstraintWhenMigrating: true, //禁止实体外键
	})
	if err != nil {
		logrus.Fatalf("数据库连接失败：%s", err)
	}
	//配置连接池
	sqlDB, err := db.DB() //获取数据库连接池管理器，是一个*sql.DB对象
	if err != nil {
		logrus.Fatalf("获取数据库连接池管理器失败，%s", err)
	}
	err = sqlDB.Ping() //验证和预热连接池，是第一个连接，会变成空闲连接
	if err != nil {
		logrus.Fatalf("数据库连接测试失败，%s", err)
	}
	//设置连接池，使用过的连接，会根据最大空闲数来选择是否关闭，放入连接池
	sqlDB.SetMaxIdleConns(10)  //设置最大空闲连接，备用，少了会频繁创建，多了会浪费资源
	sqlDB.SetMaxOpenConns(100) //设置最大连接数，包括空闲连接
	sqlDB.SetConnMaxLifetime(time.Hour)
	logrus.Infof("%s数据库连接成功", cfg.Mode)
	return
}
