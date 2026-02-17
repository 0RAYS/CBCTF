package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"

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
		if _, ret := r.Insert(role); !ret.OK && ret.Msg != i18n.Model.DuplicateKeyValue {
			return ret
		}
	}
	return model.SuccessRetVal()
}
