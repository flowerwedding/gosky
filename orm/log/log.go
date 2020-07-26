package log

import (
	"log"
	"os"
	"sync"
)

var (
	//info、error、disable为不同日志等级，输出不同颜色
	//log.New() 创建一个logger，参数：写入日志的地方，日志前缀，日志属性。os.Stdout 输出到控制它，log.LstdFlags|log.Lshortfile显示文件名和代码行号
	errorLog = log.New(os.Stdout, "\033[31m[error]\033[0m ", log.LstdFlags|log.Lshortfile)//红色
	infoLog  = log.New(os.Stdout, "\033[34m[info ]\033[0m ", log.LstdFlags|log.Lshortfile)//蓝色
	loggers  = []*log.Logger{errorLog, infoLog}
	mu       sync.Mutex
)

var (
	Error  = errorLog.Println
	Errorf = errorLog.Printf
	Info   = infoLog.Println
	Infof  = infoLog.Printf
)