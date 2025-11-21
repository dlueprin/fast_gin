package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware 认证中间件，就是一个视图
func AuthMiddleware(c *gin.Context) {
	fmt.Println("Auth请求")
	c.Next()
	fmt.Println("Auth相应")
}
