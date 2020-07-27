package gosky

import (
	"net/http"
	"path"
	"strings"
)

//留下绝对路径的后一半路径
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	//path.Join(path1,path2,…)路径片段使用特定的分隔符'\'连接起来形成路径，并规范化生成的路径。若任意一个路径片段类型错误，会报错。
	//path.resolve()把一个路径或路径片段的序列解析为一个绝对路径。
	absolutePath := path.Join(group.prefix, relativePath)
	//http.FileServer(fs)明确静态文件的根目录在fs，但是URL以绝对路径开头，如果有人请求 绝对路径+文件名，需要找到 fs+文件名 的文件。
	//因此从URL中过滤掉 绝对路径，用 fs 代替，返回参数为 fs+文件名 即从绝对路径变为相对路径
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		// 检查文件是否存在或者是否我们想要的
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

func (group *RouterGroup) Static(relativePath string, root string) {
	group.StaticFS(relativePath, http.Dir(root))
}

func (group *RouterGroup) StaticFS(relativePath string, root http.FileSystem) {
	//strings.Contains()里面是否含有字串
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static folder")
	}

	handler := group.createStaticHandler(relativePath, root)//相对路径
	//path.Join()将路径片段使用特定的分隔符（window：\）连接起来形成路径
	urlPattern := path.Join(relativePath, "/*filepath")//新的路由名字

	group.GET(urlPattern, handler)//注册路由，动态路由
}