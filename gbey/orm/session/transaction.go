package session

import "what-unexpected-summer/gbey/gbey/orm/log"

//实现事务和SQL原生语句很接近。
//调用 db.Begin() 得到 *sql.Tx 对象，使用 tx.Exec() 执行一系列操作，如果发生错误，通过 tx.Rollback() 回滚，如果没有发生错误，则通过 tx.Commit() 提交
//封装事务的 Begin、Commit 和 Rollback 三个接口
//封装的另一个目的是统一打印日志，方便定位问题。

//tx := db.Begin() 开始事务
func (s *Session) Begin() (err error) {
	log.Info("transaction begin")
	if s.tx, err = s.db.Begin(); err != nil {//调用 s.db.Begin() 得到 *sql.Tx 对象，赋值给 s.tx。
		log.Error(err)
		return
	}
	return
}

//tx.Create(...) 在事务中做的一些数据库操作(使用tx,而不是db)

//tx.Commit() 提交事务
func (s *Session) Commit() (err error) {
	log.Info("transaction commit")
	if err = s.tx.Commit(); err != nil {
		log.Error(err)
	}
	return
}

//tx.Rollback() 发生错误时回滚事务
func (s *Session) Rollback() (err error) {
	log.Info("transaction rollback")
	if err = s.tx.Rollback(); err != nil {
		log.Error(err)
	}
	return
}