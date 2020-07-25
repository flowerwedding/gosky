package session
//Hook 钩子，提前在可能增加功能的地方预设一个钩子，当需要重新修改或者增加这个地方的逻辑时，把扩展的类和方法加载到这个点上。
//扩展点的选择，钩子应设在代码可能发生改变的地方，例如记录的增删改查前后。
import (
	"gosky/orm/log"
	"reflect"
)

// 钩子与结构体绑定，每个结构需要实现自己的钩子
const (
	BeforeQuery  = "BeforeQuery"
	AfterQuery   = "AfterQuery"

	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"

	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"

	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
)

//调用钩子，同样反射实现
//将 s *Session 作为入参调用。每一个钩子的入参类型均是 *Session
func (s *Session) CallMethod(method string, value interface{}) {
	fm := reflect.ValueOf(s.RefTable().Model).MethodByName(method)//使用 MethodByName 方法反射得到该对象的方法。
	if value != nil {
		fm = reflect.ValueOf(value).MethodByName(method)
	}
	param := []reflect.Value{reflect.ValueOf(s)}//参数
	if fm.IsValid() {//检查指定对象是否创建了实例
		if v := fm.Call(param); len(v) > 0 {//v是结果接受param的返回类型，Call()可调用函数或方法
			if err, ok := v[0].Interface().(error); ok {
				log.Error(err)
			}
		}
	}
	return
}