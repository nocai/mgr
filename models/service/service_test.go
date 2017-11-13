package service

import (
	"testing"
	"fmt"
)

func TestNewError(t *testing.T) {
	err1 := NewError(MsgQuery, ErrArgument)
	err2 := NewError(MsgDataDuplication, err1)
	fmt.Println(err2)
}
