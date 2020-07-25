package gbey

import (
	"html/template"
	"log"
	"net/http"
	"path"
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

//创建一个新的路由，也就是给引擎里面的所有参数都初始化
func New() *Engine {
	//return &Engine{router: make(map[string]HandlerFunc)}//初始化一个新的map

	//return &Engine{router: newRouter()}

	engine := &Engine{router: newRouter()}
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

//之前设计动态路由时，支持通配符 * 匹配多级子路径。
//静态文件路径是相对路径。映射到真实文件后，将文件返回，静态服务器就实现了。
//找到路径后，用 net/http 库返回。因此，需要将解析请求的地址，映射到服务器上文件的真实地址，交给http.FileServer处理。

func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	//path.Join(path1,path2,…)路径片段使用特定的分隔符'\'连接起来形成路径，并规范化生成的路径。若任意一个路径片段类型错误，会报错。
	//path.resolve()把一个路径或路径片段的序列解析为一个绝对路径。
	absolutePath := path.Join(group.prefix, relativePath)
	//
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

// serve static files
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	// Register GET handlers
	group.GET(urlPattern, handler)
}

func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}


// Default use Logger() & Recovery middlewares
func Default() *Engine {
	engine := New()
	engine.Use(Logger(), Recovery())
	return engine
}