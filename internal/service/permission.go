package service

import (
	"gorm.io/gorm"

	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/model"
)

func ListPermissions(tx *gorm.DB, form dto.ListModelsForm) ([]model.Permission, int64, model.RetVal) {
	return db.InitPermissionRepo(tx).List(form.Limit, form.Offset)
}

func UpdatePermission(tx *gorm.DB, permission model.Permission, form dto.UpdatePermissionForm) model.RetVal {
	return db.InitPermissionRepo(tx).Update(permission.ID, db.UpdatePermissionOptions{
		Description: form.Description,
	})
}
