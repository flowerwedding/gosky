package gosky

import (
	"net/http"
)

type H map[string]interface{}

const defaultMultipartMemory = 32 << 20 // 表单限制上传大小，默认 32 MB

type Context struct {
	Writer http.ResponseWriter
	Req    *http.Request

	Path   string
	Method string

	Params map[string]string//解析后的路由参数

	StatusCode int//状态码

	handlers []HandlerFunc//中间件
	index    int

	engine *Engine//页面渲染

	MaxMultipartMemory  int64//文件上传

	Render Redirect//重定向
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		//记录当前执行到第几个中间件，在中间件中调用Next方法，控制权就交给下个中间件，直到最后一个中间件，再从后往前，调用每个中间件在Next方法后定义的部分
		index:  -1,
		MaxMultipartMemory: defaultMultipartMemory,//文件上传
	}
}

//获取GET参数
func (c *Context) Query(key string) string {
	value := c.Req.URL.Query().Get(key)
	return value
}

func (c *Context) DefaultQuery(key string,defaultValue string) string {
	if value := c.Req.URL.Query().Get(key); value != ""{
		return value
	}
	return defaultValue
}

//获取POST参数
func (c *Context) PostForm(key string) string {
	value := c.Req.FormValue(key)
	return value
}

func (c *Context) DefaultPostForm(key string,defaultValue string) string {
	if value := c.Req.FormValue(key); value != ""{
		return value
	}
	return defaultValue
}

//获取路由中对应的值
func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

//s，也就是长度比中间件的总数大 1，因为有Handler在
func (c *Context) Next() {//中间件
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}