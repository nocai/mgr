package key

import (
	"bytes"
	"strconv"
)

//type Key interface {
//	GetPage() int64
//	GetRows() int64
//	GetOrderBySql() string
//	GetLimitSql() string
//}

type Key struct {
	page int64
	rows int64

	Sort  []string
	Order []string
}

func (key Key) GetPage() int64 {
	return key.page
}

func (key Key) GetRows() int64 {
	return key.rows
}

func (key *Key) GetOrderBySql(alias string) string {
	if len(key.Sort) > 0 && len(key.Order) > 0 {
		var sql bytes.Buffer
		sql.WriteString(" Order by")
		for i := 0; i < len(key.Sort); i++ {
			s := key.Sort[i]
			o := key.Order[i]
			if s == "" || o == "" {
				if alias != "" {
					sql.WriteString(" ")
					sql.WriteString(alias)
					sql.WriteString(".")
					sql.WriteString("id desc")
					return sql.String()
				}
				sql.WriteString(" id desc")
				return sql.String()
			}

			sql.WriteString(" ")
			if alias != "" {
				sql.WriteString(alias)
				sql.WriteString(".")
			}
			sql.WriteString(s)
			sql.WriteString(" ")
			sql.WriteString(o)
			if i != len(key.Sort)-1 {
				sql.WriteString(",")
			}
		}
		return sql.String()
	}
	if alias != "" {
		return " " + alias + ".id desc"
	}
	return " id desc"
}

func (key *Key) GetLimitSql() string {
	if key.page > 0 && key.rows > 0 {
		startIndex := (key.page - 1) * key.rows
		return " limit " + strconv.FormatInt(startIndex, 10) + ", " + strconv.FormatInt(key.rows, 10)
	}
	return ""
}

func New(page, rows int64, sort, order []string) *Key {
	if len(sort) != len(order) {
		panic("Sort 与 Order 长度不相等")
	}
	return &Key{page: page, rows: rows, Sort: sort, Order: order}
}
