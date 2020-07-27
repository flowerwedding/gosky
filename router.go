package gosky

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node//存储每种请求方式的Tire树根节点
	handlers map[string]HandlerFunc//存储每种请求方式的handlerfunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

//解析路由
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {//路由 {path : '*'}通常用于客户端404错误
				break
			}
		}
	}
	return parts
}

//新增路由，roots的key为方法名，value是tire树，存节点；handler的key是方法+名字，value是需要实现的功能
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)

	key := method + "-" + pattern
	_, ok := r.roots[method]//建立方法节点
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)//建立分路由节点
	r.handlers[key] = handler
}

//查找路由映射表，如果查到，就执行注册的处理方法
func (r *router) handle(c *Context) {
    n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern

		c.handlers = append(c.handlers, r.handlers[key])//
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
		    c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	    })
    }
    c.Next()//
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
		for index, part := range parts {//解析 : 和 * 两种匹配符的参数，
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}