package log

import (
	"io/ioutil"
	"os"
)

const (
	InfoLevel = iota//自增
	ErrorLevel
	Disabled
)

func SetLevel(level int) {
	mu.Lock()
	defer mu.Unlock()

	for _, logger := range loggers {
		//log.SetOutput(file)    设置输出流
		logger.SetOutput(os.Stdout)
	}

	if ErrorLevel < level {
		//Discard 是一个 io.Writer 接口，调用它的 Write 方法将不做任何事情 。并且始终成功返回。
		//不打印该日志
		errorLog.SetOutput(ioutil.Discard)
	}
	if InfoLevel < level {
		infoLog.SetOutput(ioutil.Discard)
	}
}