# gosky框架

2020年7月18日至7月27日24点，重庆邮电大学红岩网校工作站web研发部后端暑假大作业。

该框架模仿gin框架和grom框架，基本架构可分为MVC三层。

## 功能简介

### controller设计

#### 快速入门

```
package main

import "gosky"

func main() {
    r := gosky.Default()
    r.GET("/ping", func(c *gosky.Context) {
        c.JSON(200, gosky.H{
            "message": "pong",
        })
    })
    r.Run() // listen and serve on 0.0.0.0:8080
}
```

#### 使用GET,POST,DELETE

```
router.GET("/someGet", getting)

router.POST("/somePost", posting)

router.DELETE("/someDelete", deleting)
```

#### 获取路径中的参数

```
name := c.Param("name")
```

#### 获取GET参数

```
firstname := c.DefaultQuery("firstname", "Guest")
 
lastname := c.Query("lastname")
```

#### 获取POST参数

```
message := c.PostForm("message")

nick := c.DefaultPostForm("nick", "anonymous")
```

#### 上传文件

```
//单个文件
file, _ := c.FormFile("file")

//多个文件
form, _ := c.MultipartForm()
files := form.File["upload[]"]

//上传到指定路径
c.SaveUploadedFile(file, dst)
```

#### 路由分组

```
v1 := r.Group("/v1")
{
	v1.GET("/", func(c *Context) {
		c.HTML(http.StatusOK, "b" ,"<h1>Hello world</h1>")
	})

	v1.GET("/hello", func(c *Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})
}
```

#### 使用中间件

```
r.Use(gee.Logger())

v2 := r.Group("/v2")
v2.Use(onlyForV2()) // v2 group middleware
{
    ...
}
```

#### 重定向

```
c.Redirect(http.StatusMovedPermanently, "https://mail.qq.com")
```

#### 设置并获取session

#### 错误处理

### model设计

#### 数据库连接

```
engine, _ := orm.NewEngine("sqlite3", "orm.db")

defer engine.Close()
```

#### 数据库迁移

 支持字段的新增与删除，不支持字段类型变更等

```
engine.Migrate(&User{})
```

#### 表操作

```
//数据库该结构的表模板初始化
s := NewSession().Model(&User{})

//数据库表删除
_ = s.DropTable()

//数据库表创建
_ = s.CreateTable()

//判断数据库表是否存在
if !s.HasTable() {
	t.Fatal("Failed to create table User")
}
```

#### 记录操作

```
//新增记录user
affected, err := s.Insert(user)

//删除记录
affected, _ := s.Where("Name = ?", "Tom").Delete()

//更新记录
affected, _ := s.Where("Name = ?", "Tom").Update("Age", 30)

//查询记录放入users
var users []User
err := s.Find(&users)

//查询第一条记录
s.First(&user)
```

#### Limit

指定要检索的记录数

```
err := s.Limit(1).Find(&users)
```

#### Order

 在从数据库检索记录时指定顺序，将重排序设置为`true`以覆盖定义的条件 

```
u := &User{}
_ = s.OrderBy("Age DESC").First(u)
```

#### Count

 获取模型的记录数 

```
count, _ := s.Count()
```

#### 日志处理

orm内置简易log库，默认情况下，支持日志分级、颜色区分、打印对应的文件名和行号。

```
t.Fatal("expect 2, but got", count)
```

#### 事务

```
s := engine.NewSession()

// 开始事务
err := s.Begin()

// ...

// 发生错误时回滚事务
err = s.Rollback()

// 或提交事务
err = s.Commit()

//事务封装
_, err := engine.Transaction(func(s *session.Session) (result interface{}, err error) {
    ...
    return err
})
```

### view设计

#### 静态文件处理

```
router.Static("/assets", "./assets")

router.StaticFS("/more_static", http.Dir("my_file_system"))
```

#### XML、JSON、YAML和ProtoBuf 渲染（输出格式）

```
c.JSON(http.StatusOK, msg)

c.XML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})

c.YAML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})

c.ProtoBuf(http.StatusOK, data)
```

#### HTML渲染

```
r.SetFuncMap(template.FuncMap{
	"formatAsDate": formatAsDate,
})

r.LoadHTMLGlob("templates/*")

r.GET("/", func(c *gosky.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
        "title": "Main website",
    })
})
```

## 参考文档

 [Go 语言编程之旅](http://product.dangdang.com/28982027.html)

[gin-gonic/gin](https://github.com/gin-gonic/gin)

[7天用Go从零实现Web框架Gee教程](https://geektutu.com/post/gee.html)

[golang reflect 反射包](https://www.jianshu.com/p/1333fd84e3be)

[SQlite常用命令](https://www.runoob.com/sqlite/sqlite-commands.html)