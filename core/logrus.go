package core

//使用 go get github.com/sirupsen/logrus 下载
import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"sync"
)

// 这个空结构体是方法的载体，虽然看起来是空的，但调用logrus时会传入参数，如信息和level
// 主要是用来触发format的
type Mylog struct {
}

const (
	red    = 31
	yellow = 33
	blue   = 36
	gray   = 37
)

// Format logrus配置格式化，需要实现format接口
func (Mylog) Format(entry *logrus.Entry) ([]byte, error) {
	//根据日志基本设置颜色
	var leverColor int
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		leverColor = gray
	case logrus.WarnLevel:
		leverColor = yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		leverColor = red
	default:
		leverColor = blue
	}
	//缓冲区处理，将要输出但还没输出的日志写入缓冲区，方便处理
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	//自定义日期格式
	timestamp := entry.Time.Format("2006-01-02 15:04:05") //这个时间必须一致，不然就出错
	if entry.HasCaller() {
		//自定义文件路径
		funcval := entry.Caller.Function                                                 //获取当前函数调用栈信息，比如我是在main调用了logrus
		filevar := fmt.Sprintf("%s:%d", path.Base(entry.Caller.File), entry.Caller.Line) //通过文件名获取路径（方便终端点击导航），然后去掉路径只保留文件名，然后提取行号
		//自定义输出格式，\x1b[%dm是表示后面都用%d的颜色，再来一次设成0就回去了
		fmt.Fprintf(b, "[%s] \x1b[%dm[%s]\x1b[0m %s %s %s\n", timestamp, leverColor, entry.Level, filevar, funcval, entry.Message)
	}
	return b.Bytes(), nil
}

func InitLogger() {
	logrus.SetLevel(logrus.DebugLevel) //最低显示级别，普通是不会显示debug的
	logrus.SetReportCaller(true)       //来源显示
	logrus.SetFormatter(&Mylog{})      //使用自定义格式
	//logrus.SetFormatter(&logrus.JSONFormatter{}) //方便发送的json格式日志，比如发给日志中心
	logrus.AddHook(&MyHook{
		logPath: "logs", //写在这里的话虽然不能通过配置文件来配置，但是再main中的位置可以更灵活
	})
}

//日志的钩子函数，主要实现下面的三点；
//1、需要实现产生日志之后，把它写入到日志文件中去
//2、按天分片
//3、错误的日志单独存放
//logrus厉害的就是钩子函数

type MyHook struct {
	file    *os.File //日志文件的文件句柄，可以操作
	errFile *os.File //错误日志文件句柄
	date    string   //日志当前时间
	logPath string
	mu      sync.Mutex
}

// Fire 日志的钩子函数，要传指针
func (hook *MyHook) Fire(entry *logrus.Entry) error {
	//1.写入到文件
	//2.按时间分片，时间轮转
	//3.错误的日志单独存放
	hook.mu.Lock()
	//加锁防止并发
	defer hook.mu.Unlock()
	msg, _ := entry.String() //日志格式化：获取日志为字符串来写入文件，放外边主要是因为这个操作是通用且必要的
	date := entry.Time.Format("2006-01-02")
	if hook.date != date {
		//换时间换文件对象
		hook.rotateFile(date)
		hook.date = date
	}
	if entry.Level <= logrus.ErrorLevel {
		hook.errFile.Write([]byte(msg)) //单独写在结构体里
	}
	hook.file.Write([]byte(msg)) //文件写入
	return nil
}
func (hook *MyHook) rotateFile(timer string) error {
	if hook.file != nil {
		hook.file.Close()
	}
	if hook.file == nil { //文件初始化：判断不存在文件再去写文件，节省os开销，多次进来只会打开一次
		//创建目录
		logDir := fmt.Sprintf("%s/%s", hook.logPath, timer)
		os.MkdirAll(logDir, 0666)
		logPath := fmt.Sprintf("%s/info.log", logDir)
		//日志输出标准写法，比open功能更多，create：如果没有就创建，wronly只写，append：追加而不是覆盖,0代表8进制，666就是linux权限
		file, _ := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		hook.file = file
		//错误日志单独存放
		errLogPath := fmt.Sprintf("%s/err.log", logDir)
		errFile, _ := os.OpenFile(errLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		hook.errFile = errFile
	}
	return nil
}

func (*MyHook) Levels() []logrus.Level {
	return logrus.AllLevels //让所有的日志生效
}
