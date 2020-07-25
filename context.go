package gosky

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//构建JSON数据
//key是string类型，value是任意类型，interface{}是一个空接口，所有类型都实现了这个接口，所有可以代表所有类型
type H map[string]interface{}

//目前只包含了http.ResponseWriter和http.Request,另外提供了对Method和Path这两个常用属性的直接访问
//Context就像百宝箱，会包含很多东西，如动态路由的参数、中间件的信息等。它随着每个请求的出现而产生，请求的结束而销毁，和当前请求相关信息都在Context里面。
type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Path   string
	Method string

	//提供对参数路由的访问，解析后的参数存到Params中，且通过c.Param("lang")的方法获取对应的值
	Params map[string]string//新增

	// response info
	StatusCode int//状态码

	// middleware中间件
	handlers []HandlerFunc
	index    int

	// engine pointer页面渲染
	engine *Engine
}

//初始化
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,

		index:  -1,//新增，中间件
	}
}

//提供了访问Query和PostForm参数的方法
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

//设置状态码
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

//发送一个原始的HTTP标头[Http Header]到客户端。
//标头是服务器以HTTP协议传HTML资料到浏览器前所送出的字串，在标头与HTML文件之间需要一行分隔。在送回HTML资料前，需要传完所有的标头。
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

//提供了快速构造String/JSON/Data/HTML的快速方法
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")//纯文本格式
	c.Status(code)
	_, _ = c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")//json格式
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)//编码
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {//没有格式
	c.Status(code)
	_, _ = c.Writer.Write(data)
}

func (c *Context) HTML(code int,  name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")//网页编码
	c.Status(code)
	//_, _ = c.Writer.Write([]byte(html))

	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.Fail(500, err.Error())
	}
}

func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}

func (c *Context) Param(key string) string {//获取路由中对应的值，gin框架里面也本来就是这个作用
	value, _ := c.Params[key]
	return value
}

func (c *Context) Next() {//新增，中间件
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}