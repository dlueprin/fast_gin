package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// LimitMiddleware 限流中间件，是函数工厂，返回一个中间件类型
func LimitMiddleware(limit int) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("Limit请求")
		c.Next()
		fmt.Println("Limit相应")
	}
}
