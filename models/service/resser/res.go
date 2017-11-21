package resser

import (
	"errors"
	"time"
)

// 资源类型
type ResType int

const (
	_ ResType = iota
	// 菜单
	ResType_Menu
	// 操作
	ResType_Button
)

const Pid_Default = -1

var (
	ErrResNameExist = errors.New("资源名称存在")
)


type Res struct {
	Id      int64   `json:"id"`
	ResName string  `json:"res_name"`
	Path    string  `json:"path"`
	ResType ResType `json:"res_type"`
	Seq     int     `seq`

	Pid int64 `json:"pid"`

	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

// 多字段唯一键
func (res *Res) TableUnique() [][]string {
	return [][]string{
		{"ResName"},
	}
}

// 多字段索引
func (res *Res) TableIndex() [][]string {
	return [][]string{
		{"ResName"},
	}
}
