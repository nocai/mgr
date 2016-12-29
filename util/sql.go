package util

type Sqler struct {
	dataSql  string
	countSql string
	args     []interface{}
}

func NewSqler(dataSql string) *Sqler {
	return &Sqler{dataSql:dataSql}
}

func (sqler *Sqler) GetCountSql() string {
	return "select count(*) from (" + sqler.dataSql + ") as t_t_t"
}

func (sqler *Sqler) AppendDataSql(dataSql string) *Sqler {
	sqler.dataSql += dataSql
	return sqler
}

func (sqler *Sqler) GetDataSql() string {
	return sqler.dataSql
}

func (sqler *Sqler) AppendArg(arg interface{}) *Sqler {
	sqler.args = append(sqler.args, arg)
	return sqler
}

func (sqler *Sqler) GetArgs() []interface{} {
	return sqler.args
}