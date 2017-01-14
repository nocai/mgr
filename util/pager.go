package util

import (
	"strconv"
	"mgr/conf"
)

type PagerKey struct {
	page       int64
	rows       int64
	startIndex int64

	Data       map[string]interface{}
	sort       string
	order      string

	Key
}

func (key *PagerKey) GetSql() string {
	key.sql.WriteString(" order by ")
	if key.sort != "" && key.order != "" {
		key.sql.WriteString(key.sort)
		key.sql.WriteString(" ")
		key.sql.WriteString(key.order)
	} else {
		key.sql.WriteString("id desc")
	}

	key.sql.WriteString(" limit ")
	key.sql.WriteString(strconv.FormatInt(key.startIndex, 10))
	key.sql.WriteString(", ")
	key.sql.WriteString(strconv.FormatInt(key.rows, 10))
	return key.sql.String()
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
