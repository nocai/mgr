package util

import (
	"bytes"
)

type Key struct {
	sql  bytes.Buffer
	args []interface{}
}

func NewKey(sql string) *Key {
	key := &Key{}
	key.sql.WriteString(sql)
	return key
}

func (key *Key) GetCountSql() string {
	return "select count(*) from (" + key.sql.String() + ") as t_t_t"
}

func (key *Key) AppendSql(sql string) *Key {
	key.sql.WriteString(sql)
	return key
}

func (key *Key) GetSql() string {
	return key.sql.String()
}

func (key *Key) AppendArg(arg interface{}) *Key {
	key.args = append(key.args, arg)
	return key
}

func (key *Key) GetArgs() []interface{} {
	return key.args
}
