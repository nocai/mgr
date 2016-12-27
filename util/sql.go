package util

type Sqler struct {
	dataSql  string
	countSql string
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