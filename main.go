package main

import (
	"fast_gin/core"
	"fast_gin/flags"
	"fast_gin/global"
	"github.com/sirupsen/logrus"
)

func main() {
	core.InitLogger() //初始化日志设置
	flags.Run()
	global.Config = core.ReadConfig()
	//global.Config.DB.Port = 3307
	//fmt.Println(global.Config.DB)
	//core.DumpConfig()
	global.DB = core.InitGorm()
	global.Redis = core.InitRedis()

	logrus.Infof("你好")
	logrus.Debugf("你好")
	logrus.Warnf("你好")
	logrus.Errorf("你好")
}
