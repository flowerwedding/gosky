package dialect

import (
	"fmt"
	"reflect"
	"time"
)

type sqlite3 struct{}

//类型转换，使一个类型化的空指针
//user := (*User)(nil) user := &User{} user := new(User)，作用应该一样，用第一个可能是因为反射

var _ Dialect = (*sqlite3)(nil)

//init() 函数，包在第一次加载时，会将 sqlite3 的 dialect 自动注册到全局
func init() {
	RegisterDialect("sqlite3", &sqlite3{})
}

//将 Go 语言的类型映射为 SQLite 的数据类型
func (s *sqlite3) DataTypeOf(typ reflect.Value) string {
	//reflect.Value是将interface反射为go的值。Kind方法，用于确定它是什么类型，然后可以通过实际的类型方法比如 （Float或String）访问它实际的值. 如果需要改变它的值，可以调用对应的setter访问
	switch typ.Kind() {
	case reflect.Bool:
		return "bool"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uintptr:
		return "integer"
	case reflect.Int64, reflect.Uint64:
		return "bigint"
	case reflect.Float32, reflect.Float64:
		return "real"
	case reflect.String:
		return "text"
	case reflect.Array, reflect.Slice:
		return "blob"
	case reflect.Struct://reflect 反射 struct 动态获取time.Time类型的值
		if _, ok := typ.Interface().(time.Time); ok {
			return "datetime"
		}
	}
	panic(fmt.Sprintf("invalid sql type %s (%s)", typ.Type().Name(), typ.Kind()))
}

//在 SQLite 中判断表 tableName 是否存在的 SQL 语句
func (s *sqlite3) TableExistSQL(tableName string) (string, []interface{}) {
	args := []interface{}{tableName}
	return "SELECT name FROM sqlite_master WHERE type='table' and name = ?", args
}