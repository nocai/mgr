package test

import (
	"mgr/models"
	"time"
	"testing"
	"github.com/astaxie/beego"
)

var (
	now = time.Now()
)

func TestInsertUser(t *testing.T) {
	user := new(models.User)
	user.Username = "username1"
	user.Password = "password"
	user.CreateTime = now
	user.UpdateTime = now
	err := models.InsertUser(user)
	if err != nil {
		t.Error(err)
	}
	t.Log(user)
}

func TestGetUserById(t *testing.T) {
	ch := make(chan bool, 1)
	go func() {
		models.GetRoleById(15)
		beego.Error("GetRoleById")
		ch <- true
	}()
	user, err := models.GetUserById(1)
	beego.Error("GetUserById")
	if err != nil {
		t.Error(err)
	}
	t.Log(user)
	<-ch
}