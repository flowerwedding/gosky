package gosky

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
//	"github.com/golang/protobuf/proto"
    "gopkg.in/yaml.v2"
	"net/http"
)

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

func (c *Context) XML(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/xml")
	c.Status(code)
	encoder := xml.NewEncoder(c.Writer)//编码
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) YAML(code int, format string, values ...interface{}) {
	c.SetHeader ("Content-Type", "application/x-yaml")
	c.Status(code)
	bytes, err := yaml.Marshal(c.Data)
	if err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}

	_, _ = c.Writer.Write(bytes)
}
/*
func (c *Context) ProtoBuf(code int, obj interface{}) {//原理如下，包下载成功但无法导入
	c.SetHeader("Content-Type", "application/x-protobuf")
	c.Status(code)
	bytes, err := proto.Marshal(r.Data.(proto.Message))
	if err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}

	_, _ = c.Writer.Write(bytes)
}
*/

func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")//json格式
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)//编码
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) { //没有格式
	c.Status(code)
	_, _ = c.Writer.Write(data)
}

func (c *Context) HTML(code int,  name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")//网页编码
	c.Status(code)

	//Execute一般与New创建的模板进行配合使用，默认去寻找该名称进行数据融合
	//使用ParseFiles创建模板可以一次指定多个文件加载多个模板进来，但Execute不知道是哪个
	//New还是ParseFiles创建模板都是可以使用ExecuteTemplate
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {//根据模板文件名选择模板进行渲染
		c.Fail(500, err.Error())
	}
}