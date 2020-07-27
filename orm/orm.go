package orm

import (
	"database/sql"
	"gosky/orm/dialect"
	"gosky/orm/log"
	"gosky/orm/session"
)

type Engine struct {
	db *sql.DB
	dialect dialect.Dialect
}

//交互前的准备工作，如连接、测试数据。
func NewEngine(driver, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)//连接数据库，返回 *sql.DB
	if err != nil {
		log.Error(err)
		return
	}

	if err = db.Ping(); err != nil { //检查数据库是否能够正常连接
		log.Error(err)
		return
	}

	dial, ok := dialect.GetDialect(driver)//获取dialect实例，初始化时有创建
	if !ok {
		log.Errorf("dialect %s Not Found", driver)
		return
	}

	e = &Engine{db: db, dialect: dial}
	log.Info("Connect database success")
	return
}

//交互后的收尾工作，断开与数据库连接
func (engine *Engine) Close() {
	if err := engine.db.Close(); err != nil {
		log.Error("Failed to close database")
	}
	log.Info("Close database success")
}

//通过Engine实例创建会话，进而与数据库交互
func (engine *Engine) NewSession() *session.Session {
	return session.New(engine.db, engine.dialect)
}

type TxFunc func(*session.Session) (interface{}, error)

//事务封装
func (engine *Engine) Transaction(f TxFunc) (result interface{}, err error) {
	s := engine.NewSession()
	if err := s.Begin(); err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {//异常处理
			_ = s.Rollback()
			panic(p)
		} else if err != nil {
			_ = s.Rollback()
		} else {
			err = s.Commit()
		}
	}()

	return f(s)
}