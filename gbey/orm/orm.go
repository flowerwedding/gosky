package orm

import (
	"database/sql"
	"what-unexpected-summer/gbey/gbey/orm/dialect"
	"what-unexpected-summer/gbey/gbey/orm/log"
	"what-unexpected-summer/gbey/gbey/orm/session"
)

//交互前的准备工作，例如连接、测试数据库
//交互后的收尾工作，关闭连接

type Engine struct {//orm和用户交互的入口
	db *sql.DB
	dialect dialect.Dialect//新增
}

//连接数据库，返回 *sql.DB
//调用 db.Ping()，检查数据库是否能够正常连接
func NewEngine(driver, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)//连接数据库
	if err != nil {
		log.Error(err)
		return
	}

	if err = db.Ping(); err != nil { //db.Ping()函数来检测数据库的连接是否存在
		log.Error(err)
		return
	}

	dial, ok := dialect.GetDialect(driver)//新增，初始化的时候多了一个对象，最开始刚链接数据库的时候init函数有创建，这里是在原来的map基础上添加新值
	if !ok {
		log.Errorf("dialect %s Not Found", driver)
		return
	}

	//e = &Engine{db: db}//连接存在就初始化engine并返回
	e = &Engine{db: db, dialect: dial}
	log.Info("Connect database success")
	return
}

//断开与数据库连接
func (engine *Engine) Close() {
	if err := engine.db.Close(); err != nil {
		log.Error("Failed to close database")
	}
	log.Info("Close database success")
}

//通过Engine实例创建会话，进而与数据库交互
//连接数据库 >> defer关闭数据库 >> 创建会话
//会话(session)是通信双方从开始通信的到通信结束期间的一个上下文。这个上下文是一段位于服务器端的内存。
//会话和连接是同时建立的，两者是对同一件事情不同层次的描述。连接是物理上的客户端同服务器的通信链路，会话是逻辑上的用户同服务器的通信交互。
//连接到数据库用户开始到退出数据库结束就是会话的一个生命期，所以后面的session要复用。
func (engine *Engine) NewSession() *session.Session {
	//return session.New(engine.db)
	return session.New(engine.db, engine.dialect)
}

//NewEngine创建Engine实例时，获取driver对应的dialect
//NewSession创建Session实例时，传递dialect给构造函数New

type TxFunc func(*session.Session) (interface{}, error)

//函数作为函数参数，事务
func (engine *Engine) Transaction(f TxFunc) (result interface{}, err error) {
	s := engine.NewSession()
	if err := s.Begin(); err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {//异常处理
			_ = s.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			_ = s.Rollback() // err is non-nil; don't change it
		} else {
			err = s.Commit() // err is nil; if Commit returns error update err
		}
	}()

	return f(s)
}