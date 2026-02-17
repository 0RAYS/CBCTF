package db

import (
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
