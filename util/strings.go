package util

import (
	"strings"
	"fmt"
)

func getFieldName(fieldName string) string {
	var b  []string
	for i, s := range fieldName {
		if 'A' <= s && s <= 'Z' {
			s += 'a' - 'A'
			if i != 0 {
				b = append(b, "_")
			}
		}
		b = append(b, string(s))
	}

	r := strings.Join(b, "")
	fmt.Println(r)
	return r
}
