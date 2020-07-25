package dialect
//sql语句中的类型与Go语言中的不同，因此，需要将Go语言的类型映射为数据库的类型
//不同数据库的数据类型也不同，而orm往往需要兼容多种数据库，因此要将差异部分提取出来，实现最大程度的复用和解耦，这一部分叫dialect
//解耦：解专有的耦，变成通用的藕。
//dialect实现一些特点的SQL语句的转换
import "reflect"

var dialectsMap = map[string]Dialect{}

//这个接口里面包含两个方法，是方法
type Dialect interface {
	//Go是一种静态类型语言，即使不同变量有相同的相关类型，也不能相互赋值除非通过类型转换。
	//只要一个值实现了接口定义的方法，那么这个值就可以存储具体的值。
	//interface{} 代表一个空的方法集合并且满足任何值，只要这个值有零个或多个方法。
	//静态类型Java、C/C++、Golang，动态类型Python、Ruby。静态类型编译时发现错误，动态类型运行时发现错误
	//反射是一种检查接口变量的类型和值的机制
	//reflect.Value()通过反射获取值信息，是一些反射操作的重要类型。
	DataTypeOf(typ reflect.Value) string//用于将Go语言的类型转换为数据库的类型
	TableExistSQL(tableName string) (string, []interface{})//返回某个表是否存在SQL语句，参数是表名
}

func RegisterDialect(name string, dialect Dialect) {//注册dialect实例
	dialectsMap[name] = dialect
}

func GetDialect(name string) (dialect Dialect, ok bool) {//获取dialect实例
	dialect, ok = dialectsMap[name]
	return
}