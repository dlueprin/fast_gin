package main

import (
	"fast_gin/core"
	"fast_gin/flags"
	"fast_gin/global"
	"fast_gin/routers"
	"fast_gin/service/cron_ser"
	"github.com/sirupsen/logrus"
)

func main() {
	core.InitLogger() //初始化日志设置
	flags.Parse()
	global.Config = core.ReadConfig()
	//global.Config.DB.Port = 3307
	//fmt.Println(global.Config.DB)
	//core.DumpConfig()
	global.DB = core.InitGorm()
	global.Redis = core.InitRedis()

	flags.Run() //根据是否执行命令行布尔操作来决定是否继续运行 web 程序，有时我只是想看看版本，那么我运行完初始化，信息配置好后，拿到信息返回给我，程序的使命就完成了，就不需要继续运行

	logrus.Info("[main] flags.Run() 执行完成，准备启动定时任务")

	//定时任务
	cron_ser.CronInit()

	logrus.Info("[main] 定时任务已初始化，准备启动路由服务")

	routers.Run()

	logrus.Info("[main] 程序正常退出")
}
