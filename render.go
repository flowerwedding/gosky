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
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

//构造String/JSON/Data/HTML的快速方法
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")//纯文本格式
	c.Status(code)
	_, _ = c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) XML(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/xml")
	c.Status(code)
	encoder := xml.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) YAML(code int, format string, values ...interface{}) {
	c.SetHeader ("Content-Type", "application/x-yaml")
	c.Status(code)
	bytes, err := yaml.Marshal(c.Writer)
	if err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}

	_, _ = c.Writer.Write(bytes)
}
/*
把go.mod删除能运行，vendor里面下载不了包
func (c *Context) ProtoBuf(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/x-protobuf")
	c.Status(code)
	bytes, err := proto.Marshal(c.Writer.(proto.Message))
	if err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}

	_, _ = c.Writer.Write(bytes)
}*/

func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")//json格式
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) { //没有格式
	c.Status(code)
	_, _ = c.Writer.Write(data)
}