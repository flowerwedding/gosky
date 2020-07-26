package session

import (
	"database/sql"
	"gosky/orm/clause"
	"gosky/orm/dialect"
	"gosky/orm/log"
	"gosky/orm/schema"
	"strings"
)

//对Exec()、Query()、QueryRow()三个原生方法的封装，session会话与数据库进行交互
//封装后可统一打印日志，还可清空(s *Session).sql和(s *Session).salVars两个变量，以至session复用，即开启一次对话，执行多个sql语句。

type Session struct    {
	db      *sql.DB

	dialect  dialect.Dialect
	refTable *schema.Schema

	tx       *sql.Tx//事务

	clause   clause.Clause//子句

	sql     strings.Builder//拼接sql语句
	sqlVars []interface{}//sql占位符
}

func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		db:      db,
		dialect: dialect,
	}
}

func (s *Session) Clear() {
	s.sql.Reset()//重置
	s.sqlVars = nil
	s.clause = clause.Clause{}
}

//连接数据库
//支持事务，当 tx 不为空时，则使用 tx 执行 SQL 语句，否则使用 db 执行 SQL 语句
func (s *Session) DB() CommonDB {
	if s.tx != nil {
		return s.tx
	}
	return s.db
}

//调用SQL语句
func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}

//执行除查询外的sql语句
func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil {//SQL语句
		log.Error(err)
	}
	return
}

//单行查询
func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)
}

//多行查询
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