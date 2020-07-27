package gosky

import "html/template"

//自定义渲染函数
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

//加载模板
func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))//ParseGlob方法，批量解析名字为pattern的文件
}

func (c *Context) HTML(code int,  name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")//网页编码
	c.Status(code)

	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {//根据模板文件名选择模板进行渲染
		c.Fail(500, err.Error())
	}
}