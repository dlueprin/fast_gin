package redis_ser

import (
	"context"
	"fast_gin/global"
	"fast_gin/utils/jwts"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

// 利用redis存放token黑名单实现jwt主动过期
func Logout(token string) {
	//检查头肯有效性并获取token中的信息
	claims, err := jwts.CheckToken(token)
	if err != nil {
		logrus.Errorf("登出失败，token检查失败 %s", err)
		return
	}
	//构建redis键名
	key := fmt.Sprintf("logout_%s", token)
	//获取过期时间差（比现在的时间
	sub := claims.ExpiresAt.Sub(time.Now())
	//将键（包含了token）和过期时间差保存到redis中，实现同步过期（因为过期了token失效也就不需要存了）
	//传键空值是因为节省
	//这里传 background是因为redis的set函数需要，而且现在没有上下文，其次背景上下文可以不取消不过期，防止上下文自己过期使得set失效
	_, err1 := global.Redis.Set(context.Background(), key, "", sub).Result()
	if err1 != nil {
		logrus.Errorf("登出失败，redis保存失败 %s", err1)
	}
}

func HasLogout(token string) bool {
	key := fmt.Sprintf("logout_%s", token)
	_, err := global.Redis.Get(context.Background(), key).Result()
	if err == nil {
		return true
	}
	return false
}
