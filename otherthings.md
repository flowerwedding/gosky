```
type HandlerFunc func(*Context)
// 定义结构体，定义路由映射的处理方法。
//这个是引擎map的value，也就是对应路由要实现的功能，里面的参数也就是平时用的net/http包里面的
//type HandlerFunc func(http.ResponseWriter, *http.Request)
//这个很重要！因为不然main里面就 不能用func
```

```
type Engine struct {
    router *router   //一张路由映射表router，key由请求方法和静态路由地址构成,针对相同的路由，请求方法不同，可以映射不同的处理方法，value是用户映射的处理方法
    *RouterGroup //engine是最顶层的分组，拥有RouterGroup所有的能力，所以和路由有关的函数都交给RouterGroup实现。   
    groups []*RouterGroup//存储所有的group   
    htmlTemplates *template.Template // for html render   funcMap 
    template.FuncMap   // for html render}
```

```
//分组控制，某一组路由需要相似的处理。
//此处分组控制以前缀为区分，并且支持分组的嵌套。作用在分组上的中间件可以作用在子分组上，子分组也可以有自己特有的中间件。
```

```
func New() *Engine {}
//出Default外的，是第二重要的函数，对Engine实例执行初始化并返回
//创建一个新的路由，也就是给引擎里面的所有参数都初始化
```

```
//Logger()：输出请求日志，并标准化日志格式
//Recovery()：异常捕获，防止出现panic导致服务崩溃，同时也将异常日志的格式化输出
```

```
//定义的路由注册进去
//还是用笨办法把路由名字、请求方式、执行函数传参进去
//各种方法都调用addRouter()函数，创建路由
//为了实现路由的映射，路由就是URL路径到函数的映射
GET()
```

```
//先定义类型handlerfunc，用来定义路由映射的处理方法。
//在engine中，添加一张路由映射表router，key由请求方法和静态路由地址构成，value是用户映射的处理方法。
//调用(*Engine).GET()方法时，会将路由和处理方法注册到映射表router中，(*Engine).Run()方法，是ListenAndServe的包装。
//Engine实现的ServiceHTTP方法的作业是，解析请求的路径，查找路由映射表，如果查到，就执行对应的方法；查不到就404.
```

```
strings.HasPrefix()函数用来检测字符串是否以指定的前缀开头。
```

```
//构建JSON数据
//key是string类型，value是任意类型，interface{}是一个空接口，所有类型都实现了这个接口，所有可以代表所有类型
type H map[string]interface{}
```

```
//目前只包含了http.ResponseWriter和http.Request,另外提供了对Method和Path这两个常用属性的直接访问
//Context就像百宝箱，会包含很多东西，如动态路由的参数、中间件的信息等。它随着每个请求的出现而产生，请求的结束而销毁，和当前请求相关信息都在Context里面。
type Context struct 
//提供对参数路由的访问，解析后的参数存到Params中，且通过c.Param("lang")的方法获取对应的值
Params map[string]string
```