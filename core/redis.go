package core

import (
	"context"
	"fast_gin/global"
	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
	"github.com/sirupsen/logrus"
	"time"
)

func InitRedis() (client *redis.Client) {
	cfg := global.Config.Redis
	if cfg.Addr == "" {
		logrus.Warnf("未配置redis连接")
		return
	}
	client = redis.NewClient(&redis.Options{
		Addr:        cfg.Addr,
		Password:    cfg.Password,
		DB:          cfg.DB,
		DialTimeout: time.Second,
		MaintNotificationsConfig: &maintnotifications.Config{
			Mode: maintnotifications.ModeDisabled,
		},
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		logrus.Errorf("redis连接失败：%s", err)
		return
	}
	logrus.Infof("redis连接成功")
	return
}
