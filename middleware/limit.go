package middleware

import (
	"fast_gin/utils/res"
	"github.com/gin-gonic/gin"
	"time"
)

// LimitMiddleware 限流中间件，是函数工厂，返回一个中间件类型
// 经典漏桶：水以固定速率持续流出。
// 你的代码：水不会持续流出。而是在每次新请求到来时，一次性把所有 “过期” 的水（旧时间戳）全部倒掉，然后再判断当前水量是否超标。
// 不是漏桶，是漏桶算法思想的一种高效实现，具体来说，它是一个 “滑动时间窗口计数器”。
// 它通过记录和定期清理请求的时间戳，来控制在任意一个时间窗口内的请求总数，从而达到限流的目的。
func LimitMiddleware(limit int) gin.HandlerFunc {
	return NewLimiter(limit, 1*time.Second).Middleware
}
func NewLimiter(limit int, duration time.Duration) *Limiter {
	return &Limiter{
		limit:      limit,
		duration:   duration,
		timestamps: make(map[string][]int64),
	}
}

type Limiter struct {
	limit      int                //限制的请求数量
	duration   time.Duration      //时间窗口
	timestamps map[string][]int64 //请求的时间戳，是ip地址对应时间戳切片，或者说是时间戳列表
}

// 使用方法来实现中间件
func (l *Limiter) Middleware(c *gin.Context) {
	// 获取ip
	ip := c.ClientIP()
	//检查时间戳是否存在，不存在就创建
	if _, ok := l.timestamps[ip]; !ok {
		l.timestamps[ip] = make([]int64, 0) //开辟一个ip键的空时间戳切片
	}
	//当前秒级时间戳，代表从 Unix 纪元（1970 年 1 月 1 日 00:00:00 UTC）到当前时间的秒数
	now := time.Now().Unix()
	//删除过期的
	for i := 0; i < len(l.timestamps[ip]); i++ {
		if l.timestamps[ip][i] < now-int64(l.duration.Seconds()) {
			l.timestamps[ip] = append(l.timestamps[ip][:i], l.timestamps[ip][i+1:]...) //。。。是切片展开符，表示将后面的切片拆开再接上前面的，保持一致性
			i--
		}
	}
	//检查数量是否超出限制
	if len(l.timestamps[ip]) >= l.limit {
		res.FailWithMsg("访问过于频繁", c)
		c.Abort()
		return
	}
	//添加当前请求时间戳到切片
	l.timestamps[ip] = append(l.timestamps[ip], now)
	c.Next()
}
