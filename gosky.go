package gosky

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

// 定义结构体，定义路由映射的处理方法。
//这个是引擎map的value，也就是对应路由要实现的功能，里面的参数也就是平时用的net/http包里面的
//type HandlerFunc func(http.ResponseWriter, *http.Request)
//这个很重要！因为不然main里面就 不能用func
type HandlerFunc func(*Context)

//在engine中，添加一张路由映射表router，key由请求方法和静态路由地址构成
//针对相同的路由，请求方法不同，可以映射不同的处理方法，value是用户映射的处理方法
type Engine struct {
	//router map[string]HandlerFunc
	router *router

	*RouterGroup //engine是最顶层的分组，拥有RouterGroup所有的能力，所以和路由有关的函数都交给RouterGroup实现。
	groups []*RouterGroup//存储所有的group

	htmlTemplates *template.Template // for html render
	funcMap       template.FuncMap   // for html render
}

//分组控制，某一组路由需要相似的处理。
//此处分组控制以前缀为区分，并且支持分组的嵌套。作用在分组上的中间件可以作用在子分组上，子分组也可以有自己特有的中间件。

type RouterGroup struct {//group 分组
    prefix      string
    middlewares []HandlerFunc // support middleware
    parent      *RouterGroup  // support nesting嵌套

    //Group对象，需要有访问Router的能力，因此它里面有一个指针，指向Engine，通过Engine简介地访问各种接口。
    engine      *Engine       // all groups share a Engine instance
}

//出Default外的，是第二重要的函数，对Engine实例执行初始化并返回
//创建一个新的路由，也就是给引擎里面的所有参数都初始化
func New() *Engine {
	//return &Engine{router: make(map[string]HandlerFunc)}//初始化一个新的map

	//return &Engine{router: newRouter()}

	engine := &Engine{router: newRouter()}


	/*engine := &Engine{
	    RouterGroup : RouterGroup{//路由组
	        Handlers : nil,
	        basePath : "/",
	        root : true,
	    }
	    FuncMap : template.FuncMap{},
	    RedirectTailingSlash :true,//是否自动重定向
	    RedirectFixedPath : false,//是否尝试修复当前请求路径
	    HandleMethodNotAllowed : false,//判断当前路由是否允许调用其他方法
	    ForwardedByClientIP : true,//如果开启，尽可能返回客户端真实IP
	    AppEngine : defaultAppEngine,
	    UseRawPath : false,////如果开启，使用url.RawPath获取请求参数，否正url.Path
	    UnescapePathValues :true,//对路径值进行转义
	    MaxMultipartMemory : defaultMultipartMemory,//控制最大的文件上传大小
	    trees : make(methodTrees, 0 ,9),
	    delims : render.Delims{Left: "{{",Right:"}}"},//HTML模板左右界定符
	    secureJsonPrefix : "while(1);",
	}*/

	//所有路由规则都有他管，路由组和Engine实例形成一个关联的组件
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

/*
//新增一个路由，这个函数是真的增加了一个新的，方法就是GET、POST，其他一样
//引擎里面的参数是一个map，map的key就是这个路由的方法+名字，value就是它要实现的功能
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	//key := method + "-" + pattern
	//engine.router[key] = handler

	engine.router.addRoute(method, pattern, handler)
}

//就是平时用的GET、POST，pattern是路由的名字/secondday/helloname等，后面的那个handler参数就是跟在后面的函数，也就是这个函数要执行的功能
//当用户调用(*Engine).GET()方法时，会将路由和处理方法注册到映射表 router 中
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

//POST同理
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}*/

//定义组来创建新的路由组，所有组共享同一个引擎实例
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

//addRoute函数调用了group.engine.router.addRoute来实现路由的映射。
//因为engine继承了RouterGroup的所有属性和方法，又因为(*Engine).engine是指向自己的。所以既可以添加路由，也可以通过分组添加路由。
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)

	//为了实现路由的映射，路由就是URL路径到函数的映射
	group.engine.router.addRoute(method, pattern, handler)//其实group.engine就是原来的engine，它真正增加路由的方法给engine里面的router了
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {//这个函数还是老规矩
	group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

func (group *RouterGroup) DELETE(pattern string, handler HandlerFunc) {
	group.addRoute("DELETE", pattern, handler)
}

//(*Engine).Run()方法，是 ListenAndServe 的包装
//就是平时最后那个监听端口的方法，返回里面那个http的函数第一个参数是要监听的端口、第二个参数是根页面的处理函数，可以为nil
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

//解析请求的路径，查找路由映射表，如果查到，就执行注册的处理方法。如果查不到，就返回 404 NOT FOUND 。
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
/*	key := req.Method + "-" + req.URL.Path
	if handler, ok := engine.router[key]; ok {
		handler(w, req)
	} else {
		_, _ = fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}*/

	var middlewares []HandlerFunc//新增，中间件
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}

	c := newContext(w, req)//放到Context里面
	c.handlers = middlewares//新增，中间件
	c.engine = engine//页面渲染

	engine.router.handle(c)//原来
}

//先定义类型handlerfunc，用来定义路由映射的处理方法。
//在engine中，添加一张路由映射表router，key由请求方法和静态路由地址构成，value是用户映射的处理方法。
//调用(*Engine).GET()方法时，会将路由和处理方法注册到映射表router中，(*Engine).Run()方法，是ListenAndServe的包装。
//Engine实现的ServiceHTTP方法的作业是，解析请求的路径，查找路由映射表，如果查到，就执行对应的方法；查不到就404.

func (group *RouterGroup) Use(middlewares ...HandlerFunc) {//增加中间件的函数
	group.middlewares = append(group.middlewares, middlewares...)
}

// 这个函数是重点，调用New()创建默认Engine实例，初始化阶段引入Logger()和Recovery()中间件
func Default() *Engine {
	engine := New()
	//Logger()：输出请求日志，并标准化日志格式
	//Recovery()：异常捕获，防止出现panic导致服务崩溃，同时也将异常日志的格式化输出
	engine.Use(Logger(), Recovery())
	return engine
}