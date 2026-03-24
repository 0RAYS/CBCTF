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
	var exists bool
	res := p.DB.Raw(`
		SELECT EXISTS (
			SELECT 1
			FROM permissions
			INNER JOIN role_permissions ON permissions.id = role_permissions.permission_id
			INNER JOIN roles ON role_permissions.role_id = roles.id AND roles.deleted_at IS NULL
			INNER JOIN "groups" ON roles.id = "groups".role_id AND "groups".deleted_at IS NULL
			INNER JOIN user_groups ON "groups".id = user_groups.group_id
			INNER JOIN users ON user_groups.user_id = users.id AND users.deleted_at IS NULL
			WHERE user_groups.user_id = ? AND permissions.resource LIKE ? AND permissions.deleted_at IS NULL
		)
	`, userID, "admin:%").Scan(&exists)
	if res.Error != nil {
		return false
	}
	return exists
}

func (p *PermissionRepo) GetUserPermissions(userID uint) ([]string, model.RetVal) {
	var perms []string
	res := p.DB.Table("permissions").
		Distinct().
		Joins("INNER JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("INNER JOIN roles ON role_permissions.role_id = roles.id AND roles.deleted_at IS NULL").
		Joins(`INNER JOIN "groups" ON roles.id = "groups".role_id AND "groups".deleted_at IS NULL`).
		Joins(`INNER JOIN user_groups ON "groups".id = user_groups.group_id`).
		Joins("INNER JOIN users ON user_groups.user_id = users.id AND users.deleted_at IS NULL").
		Where("user_groups.user_id = ? AND permissions.deleted_at IS NULL", userID).
		Order("permissions.name ASC").
		Pluck("permissions.name", &perms)
	if res.Error != nil {
		return nil, model.RetVal{Msg: i18n.Model.Permission.GetError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return perms, model.SuccessRetVal()
}
