package service

import (
	"errors"
	errors2 "github.com/pkg/errors"
)

const (
	MsgQuery           = "查询失败"
	MsgInsert          = "添加失败"
	MsgUpdate          = "更新失败"
	MsgDelete          = "删除失败"
	MsgArgument        = "无效参数"
	MsgDataDuplication = "数据重复"
)

var (
	ErrQuery  = errors.New(MsgQuery)
	ErrInsert = errors.New(MsgInsert)
	ErrUpdate = errors.New(MsgUpdate)
	ErrDelete = errors.New(MsgDelete)

	ErrArgument        = errors.New(MsgArgument)
	ErrDataDuplication = errors.New(MsgDataDuplication)
)

func NewError(msg string, err error) error {
	return errors2.Wrap(err, msg)
}
