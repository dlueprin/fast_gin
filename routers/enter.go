package routers

import (
	"fast_gin/global"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// go get github.com/gin-gonic/gin

func Run() {
	gin.SetMode(global.Config.System.Mode)
	r := gin.Default()
	r.Static("/uploads", "uploads") //获取静态文件，路径加里面的完整文件名

	//curl -X POST 127.0.0.1:8080/api/users/login
	g := r.Group("api") //创建api开头的路由组
	UserRouter(g)

	addr := global.Config.System.Addr()
	if global.Config.System.Mode == "release" {
		logrus.Infof("后端服务运行在：%s", addr)
	}

	r.Run(addr)
}
