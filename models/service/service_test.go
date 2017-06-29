package service

import (
	"testing"
	"fmt"
)

func TestNewError(t *testing.T) {
	err1 := NewError(ErrArgument, "aaaaaaaaaaaaa")
	err2 := NewError(err1, "BBBBBBBBBBBB")
	fmt.Println(err2)
}
