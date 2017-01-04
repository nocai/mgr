package util

import (
	"strconv"
	"mgr/conf"
	"bytes"
)

type PagerKey struct {
	page       int64
	rows       int64
	startIndex int64

	Data       map[string]interface{}
	sort       string
	order      string

	dataSql    string
	countSql   string
	args       []interface{}
}

func (key *PagerKey) GetCountSql() string {
	return "select count(*) from (" + key.dataSql + ") as t_t_t"
}

func (key *PagerKey) AppendDataSql(dataSql string) *PagerKey {
	key.dataSql += dataSql
	return key
}

func (key *PagerKey) GetDataSql() string {
	var buffer bytes.Buffer
	buffer.WriteString(key.dataSql)

	buffer.WriteString(" order by ")
	if key.sort != "" && key.order != "" {
		buffer.WriteString(key.sort)
		buffer.WriteString(" ")
		buffer.WriteString(key.order)
	} else {
		buffer.WriteString("id desc")
	}

	buffer.WriteString(" limit ")
	buffer.WriteString(strconv.FormatInt(key.startIndex, 10))
	buffer.WriteString(", ")
	buffer.WriteString(strconv.FormatInt(key.rows, 10))
	return buffer.String()
}

func (key *PagerKey) AppendArg(arg interface{}) *PagerKey {
	key.args = append(key.args, arg)
	return key
}

func (key *PagerKey) GetArgs() []interface{} {
	return key.args
}

func NewPagerKey(page, rows int64, data map[string]interface{}, sort, order string) *PagerKey {
	if page <= 0 {
		page = conf.Page
	}
	if rows <= 0 {
		rows = conf.Rows
	}

	if sort == "" {
		sort = "id"
	}
	if order == "" {
		order = "desc"
	}

	var startIndex int64;
	if page > 0 && rows > 0 {
		startIndex = (page - 1) * rows
	}
	return &PagerKey{page:page, rows: rows, startIndex:startIndex, Data: data, sort:sort, order:order}
}

type Pagination struct {
	Total int64 `json:"total"`
	Rows  interface{} `json:"rows"`
}

type Pager struct {
	Page      int64
	Rows      int64
	PageCount int64
	Pagination
}

// New
func NewPager(key *PagerKey, total int64, pageData interface{}) *Pager {
	var pageCount int64
	if total % key.rows == 0 {
		pageCount = total / key.rows
	} else {
		pageCount = total / key.rows + 1
	}

	return &Pager{Page:key.page, Rows:key.rows, PageCount:pageCount, Pagination:Pagination{Total:total, Rows:pageData}}
}
