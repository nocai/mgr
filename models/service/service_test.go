package service

import (
	"fmt"
	"testing"
)

func TestNewError(t *testing.T) {
	err1 := NewError(MsgQuery, ErrArgument)
	err2 := NewError(MsgDataDuplication, err1)
	fmt.Println(err2)
}
