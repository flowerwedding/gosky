```
//会话(session)是通信双方从开始通信的到通信结束期间的一个上下文。这个上下文是一段位于服务器端的内存。
//会话和连接是同时建立的，两者是对同一件事情不同层次的描述。连接是物理上的客户端同服务器的通信链路，会话是逻辑上的用户同服务器的通信交互。
//连接到数据库用户开始到退出数据库结束就是会话的一个生命期，在生命期内session要复用。
```

```
//NewEngine创建Engine实例时，获取driver对应的dialect
//NewSession创建Session实例时，传递dialect给构造函数New
```

```
//Go是一种静态类型语言，即使不同变量有相同的相关类型，也不能相互赋值除非通过类型转换。
//只要一个值实现了接口定义的方法，那么这个值就可以存储具体的值。
//interface{} 代表一个空的方法集合并且满足任何值，只要这个值有零个或多个方法。
//静态类型Java、C/C++、Golang，动态类型Python、Ruby。静态类型编译时发现错误，动态类型运行时发现错误
//反射是一种检查接口变量的类型和值的机制
//reflect.Value()通过反射获取值信息，是一些反射操作的重要类型。
```

```
user := (*User)(nil) 
user := &User{} 
user := new(User)
作用应该一样，用第一个可能是因为反射
```

```
//reflect.Value是将interface反射为go的值。
//Kind方法，用于确定它是什么类型，然后可以通过实际的类型方法比如 （Float或String）访问它实际的值. 如果需要改变它的值，可以调用对应的setter访问
//reflect.Struct是reflect 反射 struct 动态获取time.Time类型的值
```

```
//schema实现对象和表的转换。给的一个任意的对象，转换为关系型数据库中的表结构。
//表名--结构体名    字段名和字段类型--成员变量和类型    额外的约束条件（例如非空、主键）--成员变量的Tag
//Tag 标签，字段定义后面带上一个字符串就也就是标签，就是当年蛇形的那个，Tag运行时能被reflection包读取
//Tag.Lookup返回值，和是否找到
```

```
//reflect.ValueOf()获取变量的值，如果是地址传递获取地址，值传递获取值
//reflect.TypeOf()获取变量的类型
//Value.Type()和Value.Kind()，变量的话这两个获取到的类型都一样，结构体的话Type获得底层的类型、而Kind获取表层的类型struct
//value.Interface()获取变量的值，不过类型是interface.
//reflect.Indirect(value),value是一个指针，判断是否指针，如果是指针就通过Elem()获取指针指向的变量值
```

```
//modelType.NumField()获取实例的字段个数
//modelType.Field(i)通过下标i获取到特定的字段
//(Dialect).DataTypeOf() 转换为数据库的字段类型
//reflect.New(p.Type)创建一个新的反射对象
```

```
//type定义函数类型
//函数类型相同：形参和返回值类型、个数、顺序都相同，形参名可以不同
//generator是以任何类型为参数，返回值为字符串和接口数组的函数类型
```

```
//当字符串数量大于3或者字符串来自切片，用string.Join拼接
```

```
//结构体反射：
//Field() 根据索引，返回索引对应的结构体字段的信息
//NumField() 返回结构体成员字段数量
//FieldByName() 根据给定字符串返回字符串对应的结构体字段的信息
//FieldByIndex() 多成员访问时，根据[]int提供的每个结构体的字段索引，返回字段的信息
//FieldByNameFunc() 根据传入的匹配函数匹配需要的字段
```

```
//Scan可以把数据库取出的字段值赋值给指定的数据结构&tmp,因为中间的参数空接口的切片，这就意味着可以传入任何值
//有些字段类型无法转换成功，则会返回错误。因此在调用scan后都需要检查错误。
```

```
//链式调用
//链式调用是一种简化代码的方法，某个对象调用某个方法后，将该对象的引用或指针返回，即继续调用该对象的其他方法
//当某个对象需要一次调用多个方法来设置其属性，它就适合链式调用
//WHERE、LIMIT、ORDER BY 等查询条件语句非常适合链式调用，下面添加对应的方法
```

```
//reflect.New()创建对象的指针，reflect.Value.Elem() 来取得其实际的值、通过反射获取指针指向的元素类型
```

```
//将 s *Session 作为入参调用。每一个钩子的入参类型均是 *Session
```

```
//使用 MethodByName 方法反射得到该对象的方法。
//fm.Call()可调用fm对应的reflect.ValueOf(s.RefTable().Model)中s.RefTable().Model的函数或方法
```

```
//新增字段：ALTER TABLE table_name ADD COLUMN col_name, col_type;
//删除字段：CREATE TABLE new_table AS SELECT col1, col2, ... from old_table 从 old_table 中挑选需要保留的字段到 new_table 中
//         DROP TABLE old_table 删除 old_table
//         ALTER TABLE new_table RENAME TO old_table; 重命名 new_table 为 old_table
```