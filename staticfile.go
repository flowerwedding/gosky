package gosky

import (
	"html/template"
	"net/http"
	"path"
	"strings"
)

//之前设计动态路由时，支持通配符 * 匹配多级子路径。
//静态文件路径是相对路径。映射到真实文件后，将文件返回，静态服务器就实现了。
//找到路径后，用 net/http 库返回。因此，需要将解析请求的地址，映射到服务器上文件的真实地址，交给http.FileServer处理。

func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	//path.Join(path1,path2,…)路径片段使用特定的分隔符'\'连接起来形成路径，并规范化生成的路径。若任意一个路径片段类型错误，会报错。
	//path.resolve()把一个路径或路径片段的序列解析为一个绝对路径。
	absolutePath := path.Join(group.prefix, relativePath)
	//relativePath用fs代替掉，在整个绝对路径里面前面的改成fs然后就变成相对路径
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		// Check if file exists and/or if we have permission to access it
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// 这个方法给用户调用
func (group *RouterGroup) StaticFS(relativePath string, root http.FileSystem) {
	//strings.Contains()里面是否含有字串
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static folder")
	}
	handler := group.createStaticHandler(relativePath, root)
	//path.Join()将路径片段使用特定的分隔符（window：\）连接起来形成路径
	urlPattern := path.Join(relativePath, "/*filepath")

	// Register GET handlers
	group.GET(urlPattern, handler)//调用GET方法
}

//root完整目录
func (group *RouterGroup) Static(relativePath string, root string) {
	group.StaticFS(relativePath, http.Dir(root))
}

func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {//和engine关联
	engine.funcMap = funcMap
}

func (engine *Engine) LoadHTMLGlob(pattern string) {//自定义模板渲染函数funcMap()
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))//ParseGlob方法，批量解析名字为pattern的文件
    //与HTML呈现器关联
}