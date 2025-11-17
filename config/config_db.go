package config

import (
	"fmt"
	"github.com/glebarez/sqlite"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBMode string

const (
	DBMysqlMode  = "mysql"
	DBPgsqlMode  = "pgsql"
	DBSqliteMode = "sqlite"
)

type DB struct {
	Mode     DBMode `yaml:"mode"`
	DBName   string `yaml:"db_name"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

// 引用和结构优化，也是工厂方法，因为DB本来就在这里，放这里也很合理
func (db *DB) Dsn() gorm.Dialector {
	switch db.Mode {
	case DBMysqlMode:
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
			db.User,
			db.Password,
			db.Host,
			db.Port,
			db.DBName,
		)
		return mysql.Open(dsn) //驱动级别的open，用于创建数据库方言适配器
	case DBPgsqlMode:
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
			db.Host,
			db.User,
			db.Password,
			db.DBName,
			db.Port,
		)
		return postgres.Open(dsn)
	case DBSqliteMode:
		return sqlite.Open(db.DBName)
	default:
		logrus.Warnf("未配置数据库连接")
		return nil
	}
}
