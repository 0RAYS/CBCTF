package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"fmt"

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
		res := g.DB.Model(&model.Group{}).FirstOrCreate(&group, group)
		if res.Error != nil {
			return model.RetVal{Msg: i18n.Model.Group.GetError, Attr: map[string]any{"Error": res.Error.Error()}}
		}
		group, ret := g.GetByID(group.ID, GetOptions{Preloads: map[string]GetOptions{"Role": {}}})
		if !ret.OK {
			return ret
		}
		roleName, ok := model.DefaultGroupRoleMap[group.Name]
		if !ok {
			continue
		}
		role, ret := InitRoleRepo(g.DB).GetByUniqueField("name", roleName)
		if !ret.OK {
			return ret
		}
		if ret = g.Update(group.ID, UpdateGroupOptions{RoleID: new(role.ID)}); !ret.OK {
			return ret
		}
	}
	return model.SuccessRetVal()
}

func (g *GroupRepo) Delete(idL ...uint) model.RetVal {
	groupL, _, ret := g.List(-1, -1, GetOptions{
		Conditions: map[string]interface{}{"id": idL},
		Preloads:   map[string]GetOptions{"Users": {}},
	})
	if !ret.OK {
		if ret.Msg != i18n.Model.NotFound {
			return ret
		}
		return model.SuccessRetVal()
	}
	for _, group := range groupL {
		if ret = g.Update(group.ID, UpdateGroupOptions{
			Name: new(fmt.Sprintf("%s_deleted_%s", group.Name, utils.RandStr(6))),
		}); !ret.OK {
			return ret
		}
		for _, user := range group.Users {
			if ret = DeleteUserFromGroup(g.DB, user, group); !ret.OK {
				return ret
			}
		}
	}
	if res := g.DB.Model(&model.Group{}).Where("id IN ?", idL).Delete(&model.Group{}); res.Error != nil {
		log.Logger.Warningf("Failed to delete Group: %s", res.Error)
		return model.RetVal{Msg: i18n.Model.Group.DeleteError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}
