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

```
//发送一个原始的HTTP标头[Http Header]到客户端。
//标头是服务器以HTTP协议传HTML资料到浏览器前所送出的字串，在标头与HTML文件之间需要一行分隔。在送回HTML资料前，需要传完所有的标头。
func (c *Context) SetHeader(key string, value string)
```

```
//将和路由相关的东西从gbey里面提出来，方便下次加强
//Engine内部东西增多，原先的router只是其中的一小部分了
```

```
//router的roots里面的key是GET、POST方法，handler的key是完整的方法+名字
```

```
//实现Tire树的插入和查找功能，以运用到路由中去。
//之前的简单map存储路由表，使用map存储键值对，索引高效，但只能存储静态路由。
//动态路由就是一条路由规则可以匹配某一类型而非某一固定的路由。
//动态路由最常用的数据结构是前缀树（tire树），它的每一个节点的所有子路由都拥有相同的前缀。
//HTTP请求的路径刚好是由 / 分隔的多段构成，因此，每一段可以作为前缀树的一个节点。通过树结构查询，如果中间某一层的节点都不满足条件，就说明没有匹配到路由，查询结束。
```

```
//树节点上存储的信息//为了实现动态路由匹配，在普通的树基础上加上了isWild参数。
//例如在匹配 /p/go/doc/ 这个路由时，第一层节点，p 精准匹配到了p，第二层节点，go 模糊匹配到 :lang，那么将会把lang这个参数赋值为go，继续匹配下一层。
```

```
//便利子节点，寻找第一个，路由后面那个部分就是它或者是可以匹配任何值的精准匹配
//new()和make()区别
//内置函数 new 分配空间。传递给new 函数的是一个类型，不是一个值。返回值是 指向这个新分配的零值的指针。
//内建函数 make 分配并且初始化 一个 slice, 或者 map 或者 chan 对象。 并且只能是这三种对象。 和 new 一样，第一个参数是 类型，不是一个值。 但是make 的返回值就是这个类型（即使一个引用类型），而不是指针。
```

```
//路由最重要的是注册和匹配。
//开发服务时，注册路由规则，映射handler；访问时，匹配路由规则，查找到对应的handler。
//tire树支持节点的插入与查询。
//插入功能：递归查找每一层的节点，如果没有匹配到当前的part节点就新建一个。在匹配结束时，可以使用n.pattern == ”“ 来判断路由规则是否匹配成功。
//查询功能：同样递归查询每一层的节点，退出规则是，匹配到了 ' *' ，或者匹配到了第 len(parts) 层节点。
```

```
if len(parts) == height {
//height初值为0，就是匹配到了极限的下一层   
n.pattern = pattern
//然后就直接创建一个完整的路由pattern，parts是pattern的每一个一小部分，是用来查找用的，这里确定路由没找到，被匹配对象是 n 里面原来的值
```

```
part := parts[height]
//如果每一个路由的每一小段都一样就是找到了，part把每一小段都抽出来
child := n.matchChild(part)
//根据那一小段的名字去找它子节点上的值
if child == nil {
//没有匹配到档期那的part节点，就新建一个
```

```
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
```

```
//之前设计动态路由时，支持通配符 * 匹配多级子路径。
//静态文件路径是相对路径。映射到真实文件后，将文件返回，静态服务器就实现了。
//找到路径后，用 net/http 库返回。因此，需要将解析请求的地址，映射到服务器上文件的真实地址，交给http.FileServer处理。
```

```
//Execute一般与New创建的模板进行配合使用，默认去寻找该名称进行数据融合
//使用ParseFiles创建模板可以一次指定多个文件加载多个模板进来，但Execute不知道是哪个
//New还是ParseFiles创建模板都是可以使用ExecuteTemplate
```