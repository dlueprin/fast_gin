package core

import (
	"fast_gin/config"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

// ReadConfig 可能会有不同和环境，可能要区分，现在只是读文件
func ReadConfig() (cfg *config.Config) {
	cfg = new(config.Config)
	byteData, err := os.ReadFile("settings.yaml")
	if err != nil {
		fmt.Println("读取配置文件失败：", err) //可能有不一样的错，总之报错先打印处理罢
		return
	}
	err = yaml.Unmarshal(byteData, cfg)
	if err != nil {
		fmt.Println("配置文件格式错误：", err)
		return
	}
	return
}
