package service

import (
	"errors"
	errors2 "github.com/pkg/errors"
	"mgr/conf"
)


var (
	ErrQuery  = errors.New(conf.MsgQuery)
	ErrInsert = errors.New(conf.MsgInsert)
	ErrUpdate = errors.New(conf.MsgUpdate)
	ErrDelete = errors.New(conf.MsgDelete)

	ErrArgument        = errors.New(conf.MsgArgument)
	ErrDataDuplication = errors.New(conf.MsgDataDuplication)
)

func NewError(msg string, err error) error {
	return errors2.Wrap(err, msg)
}
