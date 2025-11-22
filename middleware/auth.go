package middleware

import (
	"fast_gin/utils/jwts"
	"fast_gin/utils/res"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// AuthMiddleware 认证中间件，就是一个视图
func AuthMiddleware(c *gin.Context) {
	token := c.GetHeader("token")
	_, err := jwts.CheckToken(token)
	if err != nil {
		logrus.Errorf("用户auth请求处token验证失败：%s", err)
		res.FailWithMsg("token认证失败", c)
		c.Abort() //直接拦截，不执行后面的中间件了
		return
	}
	logrus.Info("用户auth请求处token验证成功")
	//直接放行
	c.Next()
}

func AdminMiddleware(c *gin.Context) {
	token := c.GetHeader("token")
	claims, err := jwts.CheckToken(token)
	if err != nil {
		logrus.Errorf("管理员auth请求处token验证失败：%s", err)
		res.FailWithMsg("token认证失败", c)
		c.Abort()
		return
	}
	if claims.RoleID != 1 {
		logrus.Errorf("管理员auth请求处身份认证失败，roleID:%d", claims.RoleID)
		res.FailWithMsg("角色认证失败", c)
		c.Abort()
		return
	}
	logrus.Info("管理员auth请求处token验证成功")
	c.Next()
}
