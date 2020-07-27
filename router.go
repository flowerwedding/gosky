package gosky

import (
	"net/http"
	"strings"
)

//将和路由相关的东西从gbey里面提出来，方便下次加强
//Engine内部东西增多，原先的router只是其中的一小部分了

type router struct {
	roots    map[string]*node//新增，存储每种请求方式的Tire树根节点
	handlers map[string]HandlerFunc//存储每种请求方式的handlerfunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),//新增之后，已初始化就要初始化两个
		handlers: make(map[string]HandlerFunc),
	}
}

func parsePattern(pattern string) []string {//新增
	vs := strings.Split(pattern, "/")//解析路由，把路由以 / 为节点但分割成多个小段

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {//路由 {path : '*'}通常用于客户端404错误
				break
			}
		}
	}
	return parts//就是把一个完整的路由分成一段一段作为返回值
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {//增加路由
	parts := parsePattern(pattern)

	key := method + "-" + pattern//map 的 key 还是老的 key，roots里面的key是GET、POST方法，handler的key是完整的方法+名字
	_, ok := r.roots[method]//判断这个方法的map是否存在，没有就新增
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)//roots：确保这个方法的map有了后就新增节点，主要存的还是名字
	r.handlers[key] = handler//handler：直接新增路由，存储它想实现的函数，存的是映射
}

func (r *router) handle(c *Context) {
    n, params := r.getRoute(c.Method, c.Path)//找这个想要的路由是否存在
	if n != nil {
		c.Params = params//如果存在，放到map里面
		key := c.Method + "-" + n.pattern

		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
		    c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	    })
    }
    c.Next()
}

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]//该路由的方法存在

	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)//该路由的名字存在

	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {//解析 : 和 * 两种匹配符的参数，返回一个map
			if part[0] == ':' {
				//例如：/p/go/doc 匹配到 /p/:lang/doc 解析 {lang: "go"}
				//node路由里面存储的是 :lang ，动态路由，而不是它具体的名字
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				//strings.Join()，以后面那个字符为间隔，字符串拼接
				//例如：/static/css/geetutu.css 匹配到 /static/*filepath 解析 {filepath: "css/geektutu.css"}
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}