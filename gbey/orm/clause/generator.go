package clause
//查询语句等一般由多个子句构成，子句clause
//实现子句生成规则
import (
	"fmt"
	"strings"
)

//type定义函数类型
//函数类型相同：形参和返回值类型、个数、顺序都相同，形参名可以不同
//generator是以任何类型为参数，返回值为字符串和接口数组的函数类型
type generator func(values ...interface{}) (string, []interface{})

var generators map[Type]generator

func init() {
	generators = make(map[Type]generator)//初始化
	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[SELECT] = _select
	generators[LIMIT] = _limit
	generators[WHERE] = _where
	generators[ORDERBY] = _orderBy
	generators[UPDATE] = _update
	generators[DELETE] = _delete
	generators[COUNT] = _count
}

func genBindVars(num int) string {//拼接了很多问号
	var vars []string
	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}
	//将一系列字符串拼接成一个字符串，用 ， 分隔
	return strings.Join(vars, ", ")//当字符串数量大于3或者字符串来自切片，用string.Join拼接
}

//下面把SQL语句分解了

//向表中插入新的行
func _insert(values ...interface{}) (string, []interface{}) {
	// INSERT INTO $tableName ($fields)
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	//%v 默认格式
	return fmt.Sprintf("INSERT INTO %s (%v)", tableName, fields), []interface{}{}
}

//INSERT INTO table_name VALUES (value1, value2, value3,....)
func _values(values ...interface{}) (string, []interface{}) {
	// VALUES ($v1), ($v2), ...
	var bindStr string
	var sql strings.Builder
	var vars []interface{}
	sql.WriteString("VALUES ")
	for i, value := range values {
		v := value.([]interface{})
		if bindStr == "" {
			bindStr = genBindVars(len(v))
		}
		sql.WriteString(fmt.Sprintf("(%v)", bindStr))
		if i+1 != len(values) {
			sql.WriteString(", ")
		}
		vars = append(vars, v...)
	}
	return sql.String(), vars

}

func _select(values ...interface{}) (string, []interface{}) {
	// SELECT $fields FROM $tableName
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("SELECT %v FROM %s", fields, tableName), []interface{}{}
}

//SELECT * FROM table_name LIMIT i,n
//设定返回的记录数
func _limit(values ...interface{}) (string, []interface{}) {
	// LIMIT $num
	return "LIMIT ?", values
}

//SELECT column_name,column_name FROM table_name WHERE column_name operator value;
//过滤记录，包含条件
func _where(values ...interface{}) (string, []interface{}) {
	// WHERE $desc
	desc, vars := values[0], values[1:]
	return fmt.Sprintf("WHERE %s", desc), vars
}

//SELECT column_name,column_name FROM table_name ORDER BY column_name,column_name ASC|DESC;
//对结果集排序，默认升序
func _orderBy(values ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("ORDER BY %s", values[0]), []interface{}{}
}

//入参两个，第一个是表名，第二个是map类型，表示待更新的键值对
func _update(values ...interface{}) (string, []interface{}) {
	tableName := values[0]
	m := values[1].(map[string]interface{})
	var keys []string
	var vars []interface{}
	for k, v := range m {
		keys = append(keys, k+" = ?")
		vars = append(vars, v)
	}
	return fmt.Sprintf("UPDATE %s SET %s", tableName, strings.Join(keys, ", ")), vars
}

//入参，表名
func _delete(values ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("DELETE FROM %s", values[0]), []interface{}{}
}

//入参，表名，并复用 _select 生成器
func _count(values ...interface{}) (string, []interface{}) {
	return _select(values[0], []string{"count(*)"})
}