package session

import "gosky/orm/log"

//封装事务的 Begin、Commit 和 Rollback 三个接口
//调用 s.db.Begin() 得到 *sql.Tx 对象，赋值给 s.tx，并且统一打印日志，方便定位问题。

//开始事务
func (s *Session) Begin() (err error) {
	log.Info("transaction begin")
	if s.tx, err = s.db.Begin(); err != nil {//调用 s.db.Begin() 得到 *sql.Tx 对象，赋值给 s.tx。
		log.Error(err)
		return
	}
	return
}

//无错误就提交事务
func (s *Session) Commit() (err error) {
	log.Info("transaction commit")
	if err = s.tx.Commit(); err != nil {
		log.Error(err)
	}
	return
}

//发生错误回滚事务
func (s *Session) Rollback() (err error) {
	log.Info("transaction rollback")
	if err = s.tx.Rollback(); err != nil {
		log.Error(err)
	}
	return
}