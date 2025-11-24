package middleware

import (
	"fast_gin/service/redis_ser"
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
	if redis_ser.HasLogout(token) {
		logrus.Errorf("用户auth请求处token已登出")
		res.FailWithMsg("该账号已登出", c)
		c.Abort()
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
	if redis_ser.HasLogout(token) {
		logrus.Errorf("管理员auth请求处token已登出")
		res.FailWithMsg("该账号已登出", c)
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
	c.Set("claims", claims) //将角色信息放入上下文
	c.Next()
}

// GetAuth 也是为了视图方便用
func GetAuth(c *gin.Context) (cl *jwts.MyClaims) {
	cl = new(jwts.MyClaims)        //防止返回空指针，new了之后会返回空值
	_claims, ok := c.Get("claims") //这里虽然没报错，但是存的值的类型可能不是我们想要的，所以要断言一下
	if !ok {
		return
	}
	cl, ok = _claims.(*jwts.MyClaims)
	return
}
