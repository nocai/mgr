package userser

import (
	"github.com/astaxie/beego/orm"
	"time"
	"github.com/astaxie/beego"
	"fmt"
	"mgr/models"
	"mgr/models/service"
	"mgr/util/pager"
)

func InsertUser(user *models.User) error {
	if user == nil {
		beego.Error("user is nil")
		return service.ErrArgument
	}
	return insertUser(nil, user)
}
func insertUser(o orm.Ormer, user *models.User) error {
	if o == nil {
		o = orm.NewOrm()
	}

	if err := checkUser(user); err != nil {
		return err
	}
	user.CreateTime = time.Now()
	user.UpdateTime = time.Now()

	ormer := orm.NewOrm()
	id, err := ormer.Insert(user)
	if err != nil {
		beego.Error(err)
		return service.ErrInsert
	}
	beego.Debug(fmt.Sprintf("id = %v", id))
	return nil

}

func checkUser(user *models.User) error {
	if user == nil {
		beego.Error("user is nil")
		return service.ErrArgument
	}

	if user.Username == "" {
		beego.Error("user.Username is nil")
		return service.ErrArgument
	}
	if exist, err := IsExistOfUser(user); err != nil && exist {
		return service.ErrArgument
	}

	if user.Password == "" {
		beego.Error("user.Password is nil")
		return service.ErrArgument
	}
	return nil
}

// By Id
func GetUserById(id int64) (*models.User, error) {
	userKey := &models.UserKey{
		User :&models.User{
			Id : id,
		},
	}
	userSlice, err := FindUserByKey(userKey)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	if len(userSlice) == 0 {
		beego.Error(orm.ErrNoRows)
		return nil, orm.ErrNoRows
	} else if len(userSlice) > 1 {
		beego.Error(fmt.Sprintf("Data duplication: id = %d", id))
		panic(service.ErrDataDuplication)
	} else {
		return &userSlice[0], nil
	}
}

func IsExistOfUser(user *models.User) (bool, error) {
	userId := user.Id
	// 设置Id = 0，方便查询
	user.Id = 0
	userSlice, err := FindUserByKey(&models.UserKey{User:user})
	if err != nil {
		beego.Error(err)
		return false, err
	}

	for _, _user := range userSlice {
		if _user.Id != userId {
			beego.Info(fmt.Sprintf("user exist:user = %#v", _user))
			return true, nil
		}
	}
	// 将Id设置回来，不然role的数据不对
	user.Id = userId
	return false, nil
}

func GetUserByUsername(username string) (*models.User, error) {
	users, err := FindUserByKey(&models.UserKey{User:&models.User{Username:username}})
	if err != nil {
		beego.Error(err)
		return nil, err
	}

	if len(users) == 0 {
		return nil, orm.ErrNoRows
	} else if len(users) > 1 {
		beego.Error(service.ErrDataDuplication, fmt.Sprintf("username = %v, 重复数据 = %v", username, users))
	}
	return &users[0], nil
}

func PageUser(key *models.UserKey) (*pager.Pager, error) {
	total, err := CountUserByKey(key)
	if err != nil {
		return nil, err
	}
	users, err := FindUserByKey(key)
	if err != nil {
		return nil, err
	}
	return pager.New(key.Key, total, users), nil
}

func CountUserByKey(key *models.UserKey) (int64, error) {
	o := orm.NewOrm()
	sqler := key.NewSqler()

	var total int64
	err := o.Raw(sqler.GetCountSqlAndArgs()).QueryRow(&total)
	if err != nil {
		beego.Error(err)
		return 0, service.ErrQuery
	}
	return total, nil
}

func FindUserByKey(key *models.UserKey) ([]models.User, error) {
	o := orm.NewOrm()
	sqler := key.NewSqler()

	var users []models.User
	affected, err := o.Raw(sqler.GetSqlAndArgs()).QueryRows(&users)
	if err != nil {
		beego.Error(err)
		return users, service.ErrQuery
	}
	beego.Debug("affected = ", affected)
	if affected == 0 {
		beego.Debug(orm.ErrNoRows)
		return []models.User{}, nil
	}
	return users, nil

}

func DeleteUserById(id int64) error {
	if id == 0 {
		beego.Debug("id = %v", id)
		return service.ErrArgument
	}

	affected, err := orm.NewOrm().Delete(&models.User{Id:id})
	if err != nil {
		beego.Error(err)
		return service.ErrDelete
	}

	beego.Debug(fmt.Sprintf("affected = %v", affected))
	return nil
}
