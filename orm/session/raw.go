package session

import (
	"database/sql"
	"gosky/orm/clause"
	"gosky/orm/dialect"
	"gosky/orm/log"
	"gosky/orm/schema"
	"strings"
)

//session的核心功能是与数据库交互，因此增删改查都在里面
//这是对Exec()、Query()、QueryRow()三个原生方法的封装
//封装后不仅能够统一打印日志（例如执行的SQL语句和错误日志），而且还能清空(s *Session).sql和(s *Session).salVars两个变量。这样session就可以复用，即开启一次对话，执行多个sql语句。

type Session struct    {
	db      *sql.DB//使用sql/Open()方法链接数据库成功后返回的指针

	dialect  dialect.Dialect
	refTable *schema.Schema

	//当 tx 不为空时，则使用 tx 执行 SQL 语句，否则使用 db 执行 SQL 语句。
	tx       *sql.Tx//事务

	clause   clause.Clause

	sql     strings.Builder//拼接sql语句
	sqlVars []interface{}//sql占位符
}

//返回一个除db外的新的session切片
//func New(db *sql.DB) *Session {
//	return &Session{db: db}
//}

func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		db:      db,
		dialect: dialect,
	}
}

func (s *Session) Clear() {
	//Reset()重置
	s.sql.Reset()
	s.sqlVars = nil
	s.clause = clause.Clause{}
}

//进行与数据库链接，方便对数据库操作，这个函数是数据库和用户调用的操作数据库函数的中间层
//func (s *Session) DB() *sql.DB {
//	return s.db
//}

func (s *Session) Raw(sql string, values ...interface{}) *Session {//单行查询
	//写入文件（字符串）
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}

//Exec()执行除查询外的sql语句
func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	//调用之前log包里面的第一层错误，打印提示信息，也是请求参数
	log.Info(s.sql.String(), s.sqlVars)
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil {//数据库语句
		//调用之前log包里面的第二层错误
		log.Error(err)
	}
	return
}

//返回一条数据
func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)//数据库查询语句
}

//返回多条数据
func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

//数据库的最小函数集
type CommonDB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

var _ CommonDB = (*sql.DB)(nil)
var _ CommonDB = (*sql.Tx)(nil)

//当 tx 不为空时，则使用 tx 执行 SQL 语句，否则使用 db 执行 SQL 语句
func (s *Session) DB() CommonDB {
	if s.tx != nil {
		return s.tx
	}
	return s.db
}