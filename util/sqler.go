package util

import (
	"bytes"
)

type Sqler struct {
	sql  bytes.Buffer
	args []interface{}
}

func NewKey(sql string) *Sqler {
	key := &Sqler{}
	key.sql.WriteString(sql)
	return key
}

func (sqler *Sqler) GetCountSql() string {
	return "select count(*) from (" + sqler.sql.String() + ") as t_t_t"
}

func (sqler *Sqler) AppendSql(sql string) *Sqler {
	sqler.sql.WriteString(sql)
	return sqler
}

func (sqler *Sqler) GetSql() string {
	return sqler.sql.String()
}

func (sqler *Sqler) AppendArg(arg interface{}) *Sqler {
	sqler.args = append(sqler.args, arg)
	return sqler
}

func (sqler *Sqler) GetArgs() []interface{} {
	return sqler.args
}

func (sqler *Sqler) IsEmptySql() bool {
	return sqler.sql.Len() == 0
}