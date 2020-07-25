package gbey

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

//错误处理作为中间件
//使用defer挂载上错误恢复的函数，在这个函数中调用recover()，捕获panic，并将堆栈信息打印在日志中，向用户返回Internal Server Error
func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {//在return之后运行，先赋值后放入堆栈 defer func(参数){}
			//panic 会中止当前执行的程序并退出
			//panic 会导致程序被中止，但是在退出前，会执行完defer的内容，因此用defer func(){}，defer任务执行完后，panic再继续抛出
			//recover()函数是go语言提供的，可以防止因为一个panic而导致整个程序终止，但是recover()函数只能在defer里面有用，然后程序会恢复正常
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
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