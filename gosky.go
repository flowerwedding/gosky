package gosky

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

type HandlerFunc func(*Context)

type Engine struct {
	router *router//路由映射表

	*RouterGroup //分组
	groups []*RouterGroup

	htmlTemplates *template.Template //模板渲染
	funcMap       template.FuncMap
}

type RouterGroup struct {
    prefix      string
    middlewares []HandlerFunc // 中间件
    parent      *RouterGroup  // 嵌套

    //Group对象，需要有访问Router的能力，因此它里面有一个指向Engine的指针，通过Engine简介地访问各种接口。
    engine      *Engine
}

//创建默认的Engine实例，引入日志和错误中间件
func Default() *Engine {
	engine := New()
	engine.Use(Logger(), Recovery())
	return engine
}

//对Engine实例初始化，Engine实例像引擎，关联整个应用的运行扽、路由……
func New() *Engine {
	engine := &Engine{router: newRouter()}
    engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

//定义组来创建新的路由组，所有组共享同一个引擎实例
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,//核心
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

//添加路由
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp//核心
	log.Printf("Route %4s - %s", method, pattern)

	group.engine.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

func (group *RouterGroup) DELETE(pattern string, handler HandlerFunc) {
	group.addRoute("DELETE", pattern, handler)
}

//监听端口，这里是 ListenAndServe 的包装
//http.ListenAndServe函数第一个参数是要监听的端口、第二个参数是根页面的处理函数，可以为nil
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

//将中间件应用到某个Group
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

//解析请求的路径，查找路由映射表
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {//通过前缀判断属于哪个Group，拥有哪些中间件
			middlewares = append(middlewares, group.middlewares...)
		}
	}

	c := newContext(w, req)
	c.handlers = middlewares
	c.engine = engine//页面渲染

	engine.router.handle(c)
}