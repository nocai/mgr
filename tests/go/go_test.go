package got

import (
	"fmt"
	"testing"
)

type NillBool *bool

func TestNillBool(t *testing.T) {
	var nb NillBool
	fmt.Println(nb)
}
