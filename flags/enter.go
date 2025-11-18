package flags

import (
	"fast_gin/global"
	"flag"
	"fmt"
)

// 可以通过go run main.go -h 查看使用了什么命令行参数
type FlagOptions struct {
	File    string
	Version bool
	DB      bool
}

var Options FlagOptions

// Parse 解析命令行参数并存到全局变量备用
// go run main.go -f settings_dev.yaml
// 启动的时候用这句可以指定配置文件，然后下面是设置默认的，就是没指定的时候怎么办
func Parse() { //解析的单词
	flag.StringVar(&Options.File, "f", "settings.yaml", "配置文件路径")
	flag.BoolVar(&Options.Version, "v", false, "打印当前版本") //如果有-v参数，就会把默认的值改为true,后面run函数就会执行相应的操作
	flag.BoolVar(&Options.DB, "db", false, "迁移表结构")
	flag.Parse() //这个是执行改变并写入的总操作
}

// Run 根据解析后的参数来运行相应操作
func Run() (ok bool) {
	if Options.DB {
		fmt.Println("表结构迁移")
		return true
	}
	if Options.Version {
		fmt.Println("当前后端版本：", global.Version)
		return true
	}
	return false
}
