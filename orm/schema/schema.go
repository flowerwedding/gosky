package schema

import (
	"go/ast"
	"gosky/orm/dialect"
	"reflect"
)

//表内容结构体
type Field struct {
	Name string//字段名
	Type string//类型
	Tag  string//约束条件
}

//表属性结构体
type Schema struct {
	Model      interface{}//被映射的对象
	Name       string//表名
	Fields     []*Field//字段
	FieldNames []string//所有的字段名（列名）
	fieldMap   map[string]*Field//字段名和field的映射关系
}

//直接查字段，不遍历Fields
func (schema *Schema) GetField(name string) *Field {
	return schema.fieldMap[name]
}

//根据结构体做表名新建表
func Parse(dest interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model:    dest,
		Name:     modelType.Name(),//获取到的结构体名称做表名
		fieldMap: make(map[string]*Field),
	}

	//modelType.NumField()获取实例的字段个数，modelType.Field(i)通过下标i获取到特定的字段
	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		//p.Anonymous判断不是匿名函数，ast.IsExported是否私有字段
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,
				//(Dialect).DataTypeOf() 转换为数据库的字段类型，reflect.New(p.Type)创建一个新的反射对象
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			//判断原来的字段后面是否有标签，有则把标签放到新的字段后面
			if v, ok := p.Tag.Lookup("orm"); ok {
				field.Tag = v
			}

			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field
		}
	}
	return schema
}

//平铺
func (schema *Schema) RecordValues(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var fieldValues []interface{}
	for _, field := range schema.Fields {
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}
	return fieldValues
}