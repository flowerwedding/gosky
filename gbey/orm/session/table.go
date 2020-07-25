package session
//放置数据库表的相关
import (
	"fmt"
	"reflect"
	"strings"
	"what-unexpected-summer/gbey/gbey/orm/log"
	"what-unexpected-summer/gbey/gbey/orm/schema"
)

//Model()方法用于给refTable赋值。解析操作比较耗时，因此将解析的结果保存在成员变量refTable中，即使Model()被调用多次，如果传入的结构体名称不发生变化，就不会更新refTable的值
func (s *Session) Model(value interface{}) *Session {
	//空表或者型号不同的表就
	//型号不同是传进来的参数那个要操作的结构体和表本身属于的结构体类型不同
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable.Model) {
		//主要还是建表建字段吧
		s.refTable = schema.Parse(value, s.dialect)
	}
	return s
}

//该方法返回refTable的值，如果值不在就报错
//专门用了一个方法来返回一个成员变量的值，就为了判断为空就建表
func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		log.Error("Model is not set")
	}
	return s.refTable
}

//下面是数据库表的创建、删除和判断是否存在,用RefTable()返回数据库表和字段，调用原生SQL接口执行

//Raw().Exec()执行SQL原生语言，原先建的表和字段这里放入数据库
func (s *Session) CreateTable() error {
	table := s.RefTable()
	var columns []string
	for _, field := range table.Fields {
		columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}
	desc := strings.Join(columns, ",")
	_, err := s.Raw(fmt.Sprintf("CREATE TABLE %s (%s);", table.Name, desc)).Exec()//主要
	return err
}

func (s *Session) DropTable() error {
	_, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s", s.RefTable().Name)).Exec()
	return err
}

func (s *Session) HasTable() bool {
	sql, values := s.dialect.TableExistSQL(s.RefTable().Name)//之前就是那个定义在dialect结构体里面的方法
	row := s.Raw(sql, values...).QueryRow()
	var tmp string
	_ = row.Scan(&tmp)
	return tmp == s.RefTable().Name//就算找到也要再判断找到的是否等于表名
}