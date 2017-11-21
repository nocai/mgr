package sqler

import (
	"bytes"
	"mgr/util/key"
	"strconv"
)

type Sqler struct {
	key *key.Key // Has a key

	sql   bytes.Buffer
	args  []interface{}
	alias string
}

func (this *Sqler) SetAlias(alias string) {
	this.alias = alias
}

func (this *Sqler) GetOrderBySql() string {
	if len(this.key.Sort) > 0 && len(this.key.Order) > 0 {
		var sql bytes.Buffer
		sql.WriteString(" order by")
		for i := 0; i < len(this.key.Sort); i++ {
			s := this.key.Sort[i]
			o := this.key.Order[i]
			if s == "" || o == "" {
				if this.alias != "" {
					sql.WriteString(" ")
					sql.WriteString(this.alias)
					sql.WriteString(".")
					sql.WriteString("id desc")
					return sql.String()
				}
				sql.WriteString(" id desc")
				return sql.String()
			}

			sql.WriteString(" ")
			if this.alias != "" {
				sql.WriteString(this.alias)
				sql.WriteString(".")
			}
			sql.WriteString(s)
			sql.WriteString(" ")
			sql.WriteString(o)
			if i != len(this.key.Sort)-1 {
				sql.WriteString(",")
			}
		}
		return sql.String()
	}
	if this.alias != "" {
		return " " + this.alias + ".id desc"
	}
	return " id desc"
}

func (this *Sqler) GetLimitSql() string {
	if this.key.GetPage() > 0 && this.key.GetRows() > 0 {
		startIndex := (this.key.GetPage() - 1) * this.key.GetRows()
		return " limit " + strconv.FormatInt(startIndex, 10) + ", " + strconv.FormatInt(this.key.GetRows(), 10)
	}
	return ""
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
		sql += sqler.GetOrderBySql() + sqler.GetLimitSql()
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
