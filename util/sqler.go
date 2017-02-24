package util

import "bytes"

type Sqler struct {
	*Key

	sql  bytes.Buffer
	args []interface{}
}


func (sqler *Sqler) GetCountSql() string {
	if sqler.isEmptySql() {
		return ""
	}
	return "select count(*) from (" + sqler.sql.String() + ") as t_t_t"
}

func (sqler *Sqler) AppendSql(sql string) *Sqler {
	sqler.sql.WriteString(sql)
	return sqler
}

func (sqler *Sqler) GetSql() string {
	if sqler.isEmptySql() {
		return ""
	}
	sql := sqler.sql.String()
	if sqler.Key != nil {
		sql += sqler.Key.getOrderBySql() + sqler.Key.getLimitSql()
	}
	return sql
}

func (sqler *Sqler) AppendArg(arg interface{}) *Sqler {
	sqler.args = append(sqler.args, arg)
	return sqler
}

func (sqler *Sqler) GetArgs() []interface{} {
	return sqler.args
}

func (sqler *Sqler) isEmptySql() bool {
	return sqler.sql.Len() == 0
}