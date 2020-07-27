package session

import (
	"gosky/orm/log"
	"reflect"
)

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

func (s *Session) CallMethod(method string, value interface{}) {
	fm := reflect.ValueOf(s.RefTable().Model).MethodByName(method)
	if value != nil {
		fm = reflect.ValueOf(value).MethodByName(method)
	}
	param := []reflect.Value{reflect.ValueOf(s)}
	if fm.IsValid() {
		if v := fm.Call(param); len(v) > 0 {//Call()可调用fm对应的reflect.ValueOf(s.RefTable().Model)中s.RefTable().Model的函数或方法
			if err, ok := v[0].Interface().(error); ok {
				log.Error(err)
			}
		}
	}
	return
}