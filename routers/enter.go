package routers

import (
	"fast_gin/global"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
)

// go get github.com/gin-gonic/gin

func Run() {
	gin.SetMode(global.Config.System.Mode)
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // 调试阶段允许所有，上线再改具体域名
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "token", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.Static("/uploads", "uploads") //获取静态文件，路径加里面的完整文件名

	//curl -X POST 127.0.0.1:8080/api/users/login
	g := r.Group("api") //创建api开头的路由组
	DocRouter(g)
	UserRouter(g)
	ImageRouter(g)
	CaptchaRouter(g)

	addr := global.Config.System.Addr()
	if global.Config.System.Mode == "release" {
		logrus.Infof("后端服务运行在：%s", addr)
	}

	logrus.Info("开始启动 HTTP 服务器...")
	err := r.Run(addr)
	if err != nil {
		logrus.Errorf("HTTP 服务器启动失败：%v", err)
		panic(err)
	}
	logrus.Info("HTTP 服务器已停止")
}
