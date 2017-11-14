package sqler

import (
	"bytes"
	"mgr/util/key"
)

type Sqler struct {
	key *key.Key // Has a key

	sql  bytes.Buffer
	args []interface{}
}

func (sqler *Sqler) GetCountSqlAndArgs() (string, []interface{}) {
	sql := sqler.GetCountSql()
	args := sqler.GetArgs()
	return sql, args
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

func (sqler *Sqler) GetSqlAndArgs() (string, []interface{}) {
	sql := sqler.GetSql()
	args := sqler.GetArgs()
	return sql, args
}

func (sqler *Sqler) GetSql() string {
	if sqler.isEmptySql() {
		return ""
	}
	sql := sqler.sql.String()
	if sqler.key != nil {
		sql += sqler.key.GetOrderBySql("") + sqler.key.GetLimitSql()
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

func New(key *key.Key) *Sqler {
	return &Sqler{key: key}
}
