package arrefser

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"mgr/models"
	"mgr/models/service"
	"time"
)

func FindAdminRoleRefByKey(key *models.AdminRoleRefKey) ([]models.AdminRoleRef, error) {
	o := orm.NewOrm()
	sqler := key.NewSqler()

	var refs []models.AdminRoleRef
	affected, err := o.Raw(sqler.GetSql(), sqler.GetArgs()).QueryRows(&refs)
	if err != nil {
		beego.Error(err)
		return []models.AdminRoleRef{}, service.NewError(service.MsgQuery, err)
	}
	beego.Debug("affected = ", affected)
	if affected == 0 {
		return []models.AdminRoleRef{}, nil
	}
	return refs, nil
}

func InsertAdminRoleRef(arRef *models.AdminRoleRef) error {
	if adminId := arRef.AdminId; adminId == 0 {
		beego.Error("adminId can't be 0")
		return service.ErrArgument
	}
	if roleId := arRef.RoleId; roleId != 0 {
		beego.Error("roleId can't be 0")
		return service.ErrArgument
	}

	now := time.Now()
	if arRef.CreateTime.IsZero() {
		arRef.CreateTime = now
	}
	if arRef.UpdateTime.IsZero() {
		arRef.UpdateTime = now
	}

	o := orm.NewOrm()
	id, err := o.Insert(arRef)
	if err != nil {
		beego.Error(err)
		return service.NewError(service.MsgInsert, err)
	}
	arRef.Id = id
	return nil

}

func GrantRole(adminId int64, roleIds []int64) []int64 {
	now := time.Now()


	o := orm.NewOrm()
	o.Begin()
	res, err := o.Raw("delete from t_mgr_admin_role_ref where admin_id = ?", adminId).Exec()
	//affected, err := o.Delete(&models.AdminRoleRef{AdminId:adminId})
	if err != nil {
		o.Rollback()
		panic(service.NewError(service.MsgDelete, err))
	}
	affected, err := res.RowsAffected()
	if err != nil {
		o.Rollback()
		panic(err)
	}
	beego.Info("affected = ", affected)

	var arIds []int64
	for _, roleId := range roleIds {
		arRef := &models.AdminRoleRef{
			AdminId:    adminId,
			RoleId:     roleId,
			CreateTime: now,
			UpdateTime: now,
		}

		id, err := o.Insert(arRef)
		if err != nil {
			o.Rollback()
			panic(service.NewError(service.MsgInsert, err))
		}
		arIds = append(arIds, id)
	}
	o.Commit()
	return arIds

}
