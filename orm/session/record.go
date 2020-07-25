package session
//用于实现记录增删改查相关的代码
import (
	"errors"
	"gosky/orm/clause"
	"reflect"
)

func (s *Session) Insert(values ...interface{}) (int64, error) {
	recordValues := make([]interface{}, 0)//初始化
	for _, value := range values {
		s.CallMethod(BeforeInsert, value)//钩子

		table := s.Model(value).RefTable()
		//insert函数里面的value[0]就是表名，value[1]就是所有字段名
		s.clause.Set(clause.INSERT, table.Name, table.FieldNames)//循环调用来构造每一个子句
		recordValues = append(recordValues, table.RecordValues(value))//将要记录的参数顺序平铺
	}

	s.clause.Set(clause.VALUES, recordValues...)//调用insert后还要调用value，这两个经常搭配一起使用
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)//按照传入顺序构造最终SQL语句
	result, err := s.Raw(sql, vars...).Exec()//调用 Raw().Exec() 方法，用数据库执行构造的SQL语句
	if err != nil {
		return 0, err
	}

	s.CallMethod(AfterInsert, nil)//钩子

	return result.RowsAffected()//返回它的参数
}

//传入一个切片指针，查询的结果都保存在切片中
//将平铺开的字段构造对象
func (s *Session) Find(values interface{}) error {
	s.CallMethod(BeforeQuery, nil)//钩子

	destSlice := reflect.Indirect(reflect.ValueOf(values))
	destType := destSlice.Type().Elem()//获取切片的单个元素类型
	//使用reflect.New()方法创建一个destType实例，作为Model()参数，构造表结构
	table := s.Model(reflect.New(destType).Elem().Interface()).RefTable()

	s.clause.Set(clause.SELECT, table.Name, table.FieldNames)//搜索，找到所有符合条件记录的rows
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}

	for rows.Next() {
		dest := reflect.New(destType).Elem()
		var values []interface{}
		//将dest所有字段铺平开，构造切片value
		for _, name := range table.FieldNames {
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}
		//rows.Scan() 将该行记录每一列的值依次赋值给 values 中的每一个字段。
		if err := rows.Scan(values...); err != nil {
			return err
		}
		//将 dest 添加到切片 destSlice 中。循环直到所有的记录都添加到切片 destSlice 中

		s.CallMethod(AfterQuery, dest.Addr().Interface())//AfterQuery 钩子可以操作每一行记录

		destSlice.Set(reflect.Append(destSlice, dest))
	}

	//当rows.Next()返回false，即所有行数据都已经遍历结束后，会自动调用rows.Close()方法。
	//事务结束，共享锁
	return rows.Close()
}

// support map[string]interface{}
// also support kv list: "Name", "Tom", "Age", 18, ....
//Update接受两种入参，平铺开来的键值对和map类型的键值对，如果不是，也会自动转换成map键值对
func (s *Session) Update(kv ...interface{}) (int64, error) {
	s.CallMethod(BeforeUpdate, nil)//钩子

	m, ok := kv[0].(map[string]interface{})
	if !ok {
		m = make(map[string]interface{})
		for i := 0; i < len(kv); i += 2 {
			m[kv[i].(string)] = kv[i+1]
		}
	}
	s.clause.Set(clause.UPDATE, s.RefTable().Name, m)
	sql, vars := s.clause.Build(clause.UPDATE, clause.WHERE)//都在该有的基础上多了给where
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}

	s.CallMethod(AfterUpdate, nil)//钩子

	return result.RowsAffected()
}

// Delete records with where clause
func (s *Session) Delete() (int64, error) {
	s.CallMethod(BeforeDelete, nil)//钩子

	s.clause.Set(clause.DELETE, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.DELETE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}

	s.CallMethod(AfterDelete, nil)//钩子

	return result.RowsAffected()
}

// Count records with where clause
func (s *Session) Count() (int64, error) {
	s.clause.Set(clause.COUNT, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.COUNT, clause.WHERE)
	row := s.Raw(sql, vars...).QueryRow()
	var tmp int64
	//Scan可以把数据库取出的字段值赋值给指定的数据结构&tmp,因为中间的参数空接口的切片，这就意味着可以传入任何值
	//有些字段类型无法转换成功，则会返回错误。因此在调用scan后都需要检查错误。
	if err := row.Scan(&tmp); err != nil {
		return 0, err
	}
	return tmp, nil
}

//链式调用
//链式调用是一种简化代码的方法，某个对象调用某个方法后，将该对象的引用或指针返回，即继续调用该对象的其他方法
//当某个对象需要一次调用多个方法来设置其属性，它就适合链式调用
//WHERE、LIMIT、ORDER BY 等查询条件语句非常适合链式调用，下面添加对应的方法

//就是把本来在客户端写的语句在内部封装成一个方法
func (s *Session) Limit(num int) *Session {
	s.clause.Set(clause.LIMIT, num)
	return s
}

func (s *Session) Where(desc string, args ...interface{}) *Session {
	var vars []interface{}
	s.clause.Set(clause.WHERE, append(append(vars, desc), args...)...)
	return s
}

func (s *Session) OrderBy(desc string) *Session {
	s.clause.Set(clause.ORDERBY, desc)
	return s
}

//First只返回一条记录
//根据传入的类型，利用反射构造切片，调用 Limit(1) 限制返回的行数，调用 Find 方法获取到查询结果。
func (s *Session) First(value interface{}) error {
	dest := reflect.Indirect(reflect.ValueOf(value))
	//reflect.New()创建对象的指针，reflect.Value.Elem() 来取得其实际的值、通过反射获取指针指向的元素类型
	destSlice := reflect.New(reflect.SliceOf(dest.Type())).Elem()
	//limit() 和 find()方法都是自己定义的，Addr()对可寻址的值返回地址
	if err := s.Limit(1).Find(destSlice.Addr().Interface()); err != nil {
		return err
	}
	if destSlice.Len() == 0 {//没找到
		return errors.New("NOT FOUND")
	}
	dest.Set(destSlice.Index(0))
	return nil
}