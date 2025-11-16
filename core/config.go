package core

import (
	"fast_gin/config"
	"fast_gin/flags"
	"fast_gin/global"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
)

// ReadConfig 可能会有不同和环境，可能要区分，现在只是读文件
func ReadConfig() (cfg *config.Config) {
	cfg = new(config.Config)
	byteData, err := os.ReadFile(flags.Options.File) //现在就实现了可以自定义启动时选择配置文件，通过flag包
	if err != nil {
		logrus.Fatalf("读取配置文件失败：%s", err) //可能有不一样的错，总之报错先打印处理罢
		return
	}
	err = yaml.Unmarshal(byteData, cfg)
	if err != nil {
		logrus.Fatalf("配置文件格式错误：%s", err)
		return
	}
	logrus.Infof("%s 配置文件读取成功", flags.Options.File)
	return
}

// DumpConfig 写入配置文件，使得重启也会生效，防止在内存修改重启后不生效
func DumpConfig() {
	byteData, err := yaml.Marshal(global.Config)
	if err != nil {
		logrus.Errorf("配置文件转换错误：%s", err) //程序运行时不要用fatal
		return
	}
	err = os.WriteFile(flags.Options.File, byteData, 0666) //写入文件/写入的内容/权限
	if err != nil {
		logrus.Errorf("写入配置文件失败：%s", err)
		return
	}
	logrus.Infof("配置文件写入成功")
}
