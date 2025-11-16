package flags

import (
	"flag"
)

type FlagOptions struct {
	File string
}

var Options FlagOptions

// Run
// go run main.go -f settings_dev.yaml
// 启动的时候用这句可以指定配置文件，然后下面是设置默认的，就是没指定的时候怎么办
func Run() {
	flag.StringVar(&Options.File, "f", "settings.yaml", "配置文件路径")
	flag.Parse()
}
