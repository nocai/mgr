package service

import (
	"errors"
)

var (
	ErrQuery = errors.New("查询失败")
	ErrInsert = errors.New("添加失败")
	ErrUpdate = errors.New("更新失败")
	ErrDelete = errors.New("删除失败")

	ErrArgument = errors.New("无效参数")
	ErrDataDuplication = errors.New("数据重复")
)
