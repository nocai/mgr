package test

import (
	"testing"
	"mgr/models"
	"fmt"
)


func TestLogin(t *testing.T) {
	admin, err := models.Login("username", "password")
	if err != nil {
		t.Error(err)
	}
	t.Log(admin)
}


func TestGetAdminByUserId(t *testing.T) {
	admin, err := models.GetAdminByUserId(1, true)
	if err != nil {
		t.Error(err)
	}
	t.Log(fmt.Printf("%+v", admin))
}