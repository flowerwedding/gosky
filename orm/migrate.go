package orm
//结构体变更时，数据库表的字段自动迁移
//仅支持字段新增和删除，不支持字段类型变更

//新增字段：ALTER TABLE table_name ADD COLUMN col_name, col_type;
//删除字段：CREATE TABLE new_table AS SELECT col1, col2, ... from old_table 从 old_table 中挑选需要保留的字段到 new_table 中
//         DROP TABLE old_table 删除 old_table
//         ALTER TABLE new_table RENAME TO old_table; 重命名 new_table 为 old_table
import (
	"fmt"
	"gosky/orm/log"
	"gosky/orm/session"
	"strings"
)

func difference(a []string, b []string) (diff []string) {//找切片A和B共有的切片
	mapB := make(map[string]bool)
	for _, v := range b {
		mapB[v] = true
	}
	for _, v := range a {
		if _, ok := mapB[v]; !ok {
			diff = append(diff, v)
		}
	}
	return
}

// Migrate table
func (engine *Engine) Migrate(value interface{}) error {
	_, err := engine.Transaction(func(s *session.Session) (result interface{}, err error) {//事务
		if !s.Model(value).HasTable() {//表不存在
			log.Infof("table %s doesn't exist", s.RefTable().Name)
			return nil, s.CreateTable()
		}
		table := s.RefTable()
		rows, _ := s.Raw(fmt.Sprintf("SELECT * FROM %s LIMIT 1", table.Name)).QueryRows()//多行查询
		columns, _ := rows.Columns()
		addCols := difference(table.FieldNames, columns)
		delCols := difference(columns, table.FieldNames)
		log.Infof("added cols %v, deleted cols %v", addCols, delCols)

		for _, col := range addCols {
			f := table.GetField(col)
			sqlStr := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s;", table.Name, f.Name, f.Type)
			if _, err = s.Raw(sqlStr).Exec(); err != nil {
				return
			}
		}

		if len(delCols) == 0 {
			return
		}
		tmp := "tmp_" + table.Name
		fieldStr := strings.Join(table.FieldNames, ", ")
		s.Raw(fmt.Sprintf("CREATE TABLE %s AS SELECT %s from %s;", tmp, fieldStr, table.Name))
		s.Raw(fmt.Sprintf("DROP TABLE %s;", table.Name))
		s.Raw(fmt.Sprintf("ALTER TABLE %s RENAME TO %s;", tmp, table.Name))
		_, err = s.Exec()
		return
	})
	return err
}