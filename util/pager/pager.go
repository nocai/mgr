package pager

import (
	"github.com/astaxie/beego"
	"mgr/util/key"
)

type Pagination struct {
	Total    int64       `json:"total"`
	PageList interface{} `json:"rows"`
}

type Pager struct {
	Page      int64
	Rows      int64
	PageCount int64
	Pagination
}

// New
func New(k *key.Key, total int64, pageList interface{}) *Pager {
	beego.Error(k)
	if k.GetPage() == 0 || k.GetRows() == 0 {
		return &Pager{Page: 0, Rows: 0, PageCount: 0, Pagination: Pagination{Total: total, PageList: pageList}}
	}

	var pageCount int64
	if total%k.GetRows() == 0 {
		pageCount = total / k.GetRows()
	} else {
		pageCount = total/k.GetRows() + 1
	}

	return &Pager{Page: k.GetPage(), Rows: k.GetRows(), PageCount: pageCount, Pagination: Pagination{Total: total, PageList: pageList}}
}
