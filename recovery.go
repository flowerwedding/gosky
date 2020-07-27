package gosky

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				Email("2965502421@qq.com","error",trace(message))
				//错误500，内部服务器错误
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		c.Next()
	}
}

//用来触发panic的堆栈信息
//Callers用来返回调用栈的程序计数器，第0个Caller是Callers本身，第一个是上一层的trace，第二个是再上一层的defer func，因此，为了日志的简洁，跳过前三个Caller
func trace(message string) string {
	var pcs [32]uintptr//整形，地址
	//获取与当前堆栈记录相关链的调用栈踪迹
	//函数把当前go协程调用栈上的调用标识符填入切片pcs中，返回写入到pcs中的项数。实参skip为开始在pcs中记录之前所要跳过的栈帧数，0表示Callers自身的调用栈，1表示Callers所在的调用栈，返回写入p的项数
	//跳过前三个调用栈
	n := runtime.Callers(3, pcs[:])

	var str strings.Builder//字符串拼接类型
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		//通过runtime.FuncForPC(pc)获取对应的函数，再通过fn.FileLine(pc)获取到调用该函数的文件名和行号，打印在日志中。
		//获取一个标识调用栈识别符pc对应的调用栈
		//获取调用栈所调用的函数的名字：fn.Name()
		//获取调用栈所调用的函数的所在的源文件名称和行号：fn.FileLine(pc)，其中名字和fn.Name()一样。
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}