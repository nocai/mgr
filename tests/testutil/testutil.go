package testutil

import (
	"mgr/conf"
	"mgr/util/key"
)

func GetTestKey() *key.Key {
	return key.New(conf.Page, conf.Rows, []string{}, []string{})
}
