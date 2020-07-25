package orm

import (
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"gosky/orm/session"
	"testing"
)

func OpenDB(t *testing.T) *Engine {//连接数据库
	t.Helper()
	engine, err := NewEngine("sqlite3", "orm.db")
	if err != nil {
		t.Fatal("failed to connect", err)
	}
	return engine
}

type User struct {
	Name string `geeorm:"PRIMARY KEY"`
	Age  int
}

func TestEngine_Transaction(t *testing.T) {//调用下面两个函数
	t.Run("rollback", func(t *testing.T) {
		transactionRollback(t)
	})
	t.Run("commit", func(t *testing.T) {
		transactionCommit(t)
	})
}

func transactionRollback(t *testing.T) {
	engine := OpenDB(t)//数据库的连接和最后的关闭
	defer engine.Close()

	s := engine.NewSession()//创建新的会话，如果表存在就删除
	_ = s.Model(&User{}).DropTable()
	_, err := engine.Transaction(func(s *session.Session) (result interface{}, err error) {
		_ = s.Model(&User{}).CreateTable()//一定会创建一张表并插入一条记录
		_, err = s.Insert(&User{"Tom", 18})
		//故意返回了一个自定义 error，最终事务回滚，表创建失败
		return nil, errors.New("Error")//返回一个error类型的错误
	})
	if err == nil || s.HasTable() {
		t.Fatal("failed to rollback")
	}
}

func transactionCommit(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()

	s := engine.NewSession()
	_ = s.Model(&User{}).DropTable()
	_, err := engine.Transaction(func(s *session.Session) (result interface{}, err error) {
		_ = s.Model(&User{}).CreateTable()
		_, err = s.Insert(&User{"Tom", 18})
		return
	})
	u := &User{}//并通过s.First()方法查询到插入的记录
	_ = s.First(u)
	if err != nil || u.Name != "Tom" {
		t.Fatal("failed to commit")
	}
}