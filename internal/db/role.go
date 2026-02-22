package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"slices"

	"gorm.io/gorm"
)

type RoleRepo struct {
	BaseRepo[model.Role]
}

type CreateRoleOptions struct {
	Name        string
	Description string
}

func (c CreateRoleOptions) Convert2Model() model.Model {
	return model.Role{
		Name:        c.Name,
		Description: c.Description,
	}
}

type UpdateRoleOptions struct {
	Name        *string
	Description *string
}

func (u UpdateRoleOptions) Convert2Map() map[string]any {
	options := make(map[string]any)
	if u.Name != nil {
		options["name"] = *u.Name
	}
	if u.Description != nil {
		options["description"] = *u.Description
	}
	return options
}

func InitRoleRepo(tx *gorm.DB) *RoleRepo {
	return &RoleRepo{
		BaseRepo: BaseRepo[model.Role]{
			DB: tx,
		},
	}
}

func (r *RoleRepo) InitDefaultRoles() model.RetVal {
	for _, role := range model.DefaultRoles {
		role, ret := r.Insert(role)
		if !ret.OK && ret.Msg != i18n.Model.DuplicateKeyValue {
			return ret
		}
		role, ret = r.GetByID(role.ID, GetOptions{Preloads: map[string]GetOptions{"Permissions": {}}})
		if !ret.OK {
			return ret
		}
		if len(role.Permissions) > 0 {
			continue
		}
		permissions, ok := model.DefaultRolePermissionMap[role.Name]
		if !ok {
			continue
		}
		for _, permission := range permissions {
			perm, ret := InitPermissionRepo(r.DB).GetByUniqueKey("name", permission)
			if !ret.OK {
				return ret
			}
			if slices.ContainsFunc(role.Permissions, func(permission model.Permission) bool {
				return permission.ID == perm.ID
			}) {
				continue
			}
			if ret = AssignPermissionToRole(r.DB, perm, role); !ret.OK {
				return ret
			}
		}
	}
	return model.SuccessRetVal()
}
