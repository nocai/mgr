package got

import (
	"fmt"
	"testing"
	"strings"
)

type NillBool *bool

func TestNillBool(t *testing.T) {
	var nb NillBool
	fmt.Println(nb)
}

func TestTitle(t *testing.T) {
	str := "testTitle"
	fmt.Println(strings.Title(str))
	fmt.Println(strings.Title("her royal highness"))
}

func TestToTitle(t *testing.T) {
	fmt.Println(strings.ToTitle("loud noises"))
	fmt.Println(strings.ToTitle("хлеб"))
}

