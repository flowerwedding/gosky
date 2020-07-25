package schema
//schema实现对象和表的转换。给的一个任意的对象，转换为关系型数据库中的表结构。
//表名--结构体名    字段名和字段类型--成员变量和类型    额外的约束条件（例如非空、主键）--成员变量的Tag
//Tag 标签，字段定义后面带上一个字符串就也就是标签，就是当年蛇形的那个，Tag运行时能被reflection包读取
//Tag.Lookup返回值，和是否找到
import (
	"go/ast"
	"reflect"
	"what-unexpected-summer/gbey/gbey/orm/dialect"
)

//表的字段
type Field struct {
	Name string//字段名
	Type string//类型
	Tag  string//约束条件
}

//数据库表
type Schema struct {
	Model      interface{}//被映射的对象
	Name       string//表名
	Fields     []*Field//字段
	FieldNames []string//包含所有的字段名（列名）
	fieldMap   map[string]*Field//记录字段名和field的映射关系，方便之后无需遍历Fileds
}

//直接查字段，不遍历Fields
func (schema *Schema) GetField(name string) *Field {
	return schema.fieldMap[name]
}

//将任意对象解析为Schema实例
func Parse(dest interface{}, d dialect.Dialect) *Schema {
	//reflect.ValueOf()获取变量的值，如果是地址传递获取地址，值传递获取值
	//reflect.TypeOf()获取变量的类型
	//Value.Type()和Value.Kind()，变量的话这两个获取到的类型都一样，结构体的话Type获得底层的类型、而Kind获取表层的类型struct
	//value.Interface()获取变量的值，不过类型是interface.
	//reflect.Indirect(value),value是一个指针，判断是否指针，如果是指针就通过Elem()获取指针指向的变量值
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{//初始化表
		Model:    dest,
		Name:     modelType.Name(),//获取到的结构体名称做表名
		fieldMap: make(map[string]*Field),
	}

	//NumField()获取实例的字段个数
	for i := 0; i < modelType.NumField(); i++ {
		//通过下标i获取到特定的字段
		p := modelType.Field(i)
		//p.Anonymous判断不是匿名函数，ast.IsExported是否私有字段
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,//字段名
				//(Dialect).DataTypeOf() 转换为数据库的字段类型
				//reflect.New(p.Type)创建一个新的反射对象
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),//字段类型
			}
			//如果原来的字段后面有标签，那么就把标签放到新的结构体字段后面
			if v, ok := p.Tag.Lookup("orm"); ok {
				field.Tag = v
			}
			//新表的字段相关的成员都增加
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field
		}
	}
	return schema
}

//根据数据库中列的顺序，从对象中找到对应的值，所以需要按顺序平铺
func (schema *Schema) RecordValues(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var fieldValues []interface{}
	for _, field := range schema.Fields {
		//结构体反射：
		//Field() 根据索引，返回索引对应的结构体字段的信息
		//NumField() 返回结构体成员字段数量
		//FieldByName() 根据给定字符串返回字符串对应的结构体字段的信息
		//FieldByIndex() 多成员访问时，根据[]int提供的每个结构体的字段索引，返回字段的信息
		//FieldByNameFunc() 根据传入的匹配函数匹配需要的字段
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}
	return fieldValues
}