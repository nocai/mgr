package key

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

func New(page, rows int64, sort, order []string) *Key {
	if len(sort) != len(order) {
		panic("Sort 与 Order 长度不相等")
	}
	return &Key{page: page, rows: rows, Sort: sort, Order: order}
}
