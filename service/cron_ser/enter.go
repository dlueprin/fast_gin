package cron_ser

import (
	"github.com/robfig/cron/v3"
	"time"
)

func CronInit() {
	timezone, _ := time.LoadLocation("Asia/Shanghai") //首字母大写
	//创建调度器实例，第一个参数使得支持秒，即6星，第二个参数是时区，不写默认是系统本地时区
	crontab := cron.New(cron.WithSeconds(), cron.WithLocation(timezone))
	//					  秒 分 时 日 月 星期
	crontab.AddFunc("0 11 21 * * *", HelloCron)
	crontab.Start()
}
