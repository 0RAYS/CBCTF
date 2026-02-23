package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"

	"gorm.io/gorm"
)

type PermissionRepo struct {
	BaseRepo[model.Permission]
}

type CreatePermissionOptions struct {
	Name        string
	Resource    string
	Operation   string
	Description string
}

func (c CreatePermissionOptions) Convert2Model() model.Model {
	return model.Permission{
		Name:        c.Name,
		Resource:    c.Resource,
		Operation:   c.Operation,
		Description: c.Description,
	}
}

type UpdatePermissionOptions struct {
	Description *string
}

func (u UpdatePermissionOptions) Convert2Map() map[string]any {
	options := make(map[string]any)
	if u.Description != nil {
		options["description"] = *u.Description
	}
	return options
}

func InitPermissionRepo(tx *gorm.DB) *PermissionRepo {
	return &PermissionRepo{
		BaseRepo: BaseRepo[model.Permission]{
			DB: tx,
		},
	}
}

func (p *PermissionRepo) InitPermissions() model.RetVal {
	for _, permission := range model.Permissions {
		if _, ret := p.Insert(permission); !ret.OK && ret.Msg != i18n.Model.DuplicateKeyValue {
			return ret
		}
	}
	return model.SuccessRetVal()
}

func (p *PermissionRepo) HasAdminAccess(userID uint) bool {
	var count int64
	res := p.DB.Table("permissions").
		Joins("INNER JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("INNER JOIN roles ON role_permissions.role_id = roles.id").
		Joins("INNER JOIN `groups` ON roles.id = groups.role_id").
		Joins("INNER JOIN user_groups ON groups.id = user_groups.group_id").
		Where("user_groups.user_id = ? AND permissions.deleted_at IS NULL AND permissions.resource LIKE ?", userID, "admin:%").
		Count(&count)
	if res.Error != nil {
		return false
	}
	return count > 0
}

func (p *PermissionRepo) GetUserPermissions(userID uint) ([]string, model.RetVal) {
	var perms []model.Permission
	res := p.DB.Table("permissions").
		Distinct("permissions.*").
		Joins("INNER JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("INNER JOIN roles ON role_permissions.role_id = roles.id").
		Joins("INNER JOIN `groups` ON roles.id = groups.role_id").
		Joins("INNER JOIN user_groups ON groups.id = user_groups.group_id").
		Where("user_groups.user_id = ? AND permissions.deleted_at IS NULL", userID).
		Scan(&perms)
	if res.Error != nil {
		return nil, model.RetVal{Msg: i18n.Model.GetError, Attr: map[string]any{"Model": model.Permission{}.ModelName(), "Error": res.Error.Error()}}
	}
	names := make([]string, len(perms))
	for i, perm := range perms {
		names[i] = perm.Name
	}
	return names, model.SuccessRetVal()
}

func (p *PermissionRepo) CheckUserPermission(userID uint, permission string) (bool, model.RetVal) {
	var perm model.Permission
	res := p.DB.Table("permissions").
		Joins("INNER JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("INNER JOIN roles ON role_permissions.role_id = roles.id").
		Joins("INNER JOIN `groups` ON roles.id = groups.role_id").
		Joins("INNER JOIN user_groups ON groups.id = user_groups.group_id").
		Joins("INNER JOIN users ON user_groups.user_id = users.id").
		Where("users.deleted_at IS NULL AND permissions.deleted_at IS NULL AND users.id = ? AND permissions.name = ?", userID, permission).
		Limit(1).Find(&perm)
	if res.Error != nil {
		return false, model.RetVal{Msg: i18n.Model.GetError, Attr: map[string]any{"Model": perm.ModelName(), "Error": res.Error.Error()}}
	}
	if res.RowsAffected == 0 {
		return false, model.RetVal{Msg: i18n.Model.NotFound, Attr: map[string]any{"Model": perm.ModelName()}}
	}
	return true, model.SuccessRetVal()
}
