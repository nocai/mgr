package util

import (
	"strconv"
	"mgr/conf"
)

type PagerKey struct {
	page       int64
	rows       int64
	startIndex int64

	// 过期了，不要使用
	Data       map[string]interface{}
	sort       string
	order      string

	KeyWord    string

	Sqler
}

func (key *PagerKey) GetSql() string {
	if key.IsEmptySql() {
		return ""
	}
	key.appendOrderBy()
	key.appendLimit()
	return key.Sqler.GetSql()
}

func (key *PagerKey) appendOrderBy() {
	if key.sort != "" && key.order != "" {
		key.Sqler.AppendSql(" order by ")
		key.Sqler.AppendSql(key.sort)
		key.Sqler.AppendSql(" ")
		key.Sqler.AppendSql(key.order)
	}
}

func (key *PagerKey) appendLimit() {
	if key.startIndex >= 0 && key.rows > 0 {
		key.Sqler.AppendSql(" limit ")
		key.Sqler.AppendSql(strconv.FormatInt(key.startIndex, 10))
		key.Sqler.AppendSql(", ")
		key.Sqler.AppendSql(strconv.FormatInt(key.rows, 10))
	}
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
	Total    int64 `json:"total"`
	PageList interface{} `json:"rows"`
}

type Pager struct {
	Page      int64
	Rows      int64
	PageCount int64
	Pagination
}

// New
func NewPager(key *PagerKey, total int64, pageList interface{}) *Pager {
	var pageCount int64
	if total % key.rows == 0 {
		pageCount = total / key.rows
	} else {
		pageCount = total / key.rows + 1
	}

	return &Pager{Page:key.page, Rows:key.rows, PageCount:pageCount, Pagination:Pagination{Total:total, PageList:pageList}}
}
