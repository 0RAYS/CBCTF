package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"

	"gorm.io/gorm"
)

type GroupRepo struct {
	BaseRepo[model.Group]
}

type CreateGroupOptions struct {
	RoleID      uint
	Name        string
	Description string
}

func (c CreateGroupOptions) Convert2Model() model.Model {
	return model.Group{
		RoleID:      c.RoleID,
		Name:        c.Name,
		Description: c.Description,
	}
}

type UpdateGroupOptions struct {
	RoleID      *uint
	Name        *string
	Description *string
}

func (u UpdateGroupOptions) Convert2Map() map[string]any {
	options := make(map[string]any)
	if u.RoleID != nil {
		options["role_id"] = u.RoleID
	}
	if u.Name != nil {
		options["name"] = *u.Name
	}
	if u.Description != nil {
		options["description"] = *u.Description
	}
	return options
}

func InitGroupRepo(tx *gorm.DB) *GroupRepo {
	return &GroupRepo{
		BaseRepo: BaseRepo[model.Group]{
			DB: tx,
		},
	}
}
func (g *GroupRepo) InitDefaultGroups() model.RetVal {
	for _, group := range model.DefaultGroups {
		res := g.DB.Model(&model.Group{}).FirstOrCreate(&group, "name = ?", group.Name)
		if res.Error != nil {
			return model.RetVal{Msg: i18n.Model.GetError, Attr: map[string]any{"Model": group.ModelName(), "Error": res.Error.Error()}}
		}
		group, ret := g.GetByID(group.ID, GetOptions{Preloads: map[string]GetOptions{"Role": {}}})
		if !ret.OK {
			return ret
		}
		if group.RoleID > 0 {
			continue
		}
		roleName, ok := model.DefaultGroupRoleMap[group.Name]
		if !ok {
			continue
		}
		role, ret := InitRoleRepo(g.DB).GetByUniqueKey("name", roleName)
		if !ret.OK {
			return ret
		}
		if ret = g.Update(group.ID, UpdateGroupOptions{RoleID: new(role.ID)}); !ret.OK {
			return ret
		}
	}
	return model.SuccessRetVal()
}
