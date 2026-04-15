package service

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"

	"gorm.io/gorm"
)

func GetRolePermissions(tx *gorm.DB, role model.Role) ([]model.Permission, model.RetVal) {
	role, ret := db.InitRoleRepo(tx).GetByID(role.ID, db.GetOptions{
		Preloads: map[string]db.GetOptions{"Permissions": {}},
	})
	if !ret.OK {
		return nil, ret
	}
	return role.Permissions, model.SuccessRetVal()
}

func ListRoles(tx *gorm.DB, form dto.ListModelsForm) ([]model.Role, int64, model.RetVal) {
	return db.InitRoleRepo(tx).List(form.Limit, form.Offset)
}

func CreateRole(tx *gorm.DB, form dto.CreateRoleForm) (model.Role, model.RetVal) {
	return db.InitRoleRepo(tx).Create(db.CreateRoleOptions{
		Name:        form.Name,
		Description: form.Description,
	})
}

func UpdateRole(tx *gorm.DB, role model.Role, form dto.UpdateRoleForm) model.RetVal {
	if role.Default && form.Name != nil {
		return model.RetVal{Msg: i18n.Model.Role.CannotUpdateDefault}
	}
	return db.InitRoleRepo(tx).Update(role.ID, db.UpdateRoleOptions{
		Name:        form.Name,
		Description: form.Description,
	})
}

func DeleteRole(tx *gorm.DB, role model.Role) model.RetVal {
	if role.Default {
		return model.RetVal{Msg: i18n.Model.Role.CannotDeleteDefault}
	}
	return db.WithTransactionDB(tx, func(tx2 *gorm.DB) model.RetVal {
		return db.InitRoleRepo(tx2).Delete(role.ID)
	})
}

func AssignPermission(tx *gorm.DB, role model.Role, form dto.AssignPermissionForm) model.RetVal {
	permission, ret := db.InitPermissionRepo(tx).GetByID(form.PermissionID)
	if !ret.OK {
		return ret
	}
	return db.AssignPermissionToRole(tx, permission, role)
}

func RevokePermission(tx *gorm.DB, role model.Role, form dto.AssignPermissionForm) model.RetVal {
	permission, ret := db.InitPermissionRepo(tx).GetByID(form.PermissionID)
	if !ret.OK {
		return ret
	}
	return db.RevokePermissionFromRole(tx, permission, role)
}
