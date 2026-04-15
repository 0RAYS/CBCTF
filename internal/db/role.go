package db

import (
	"CBCTF/internal/config"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"fmt"
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
		res := r.DB.Model(&model.Role{}).FirstOrCreate(&role, role)
		if res.Error != nil {
			return model.RetVal{Msg: i18n.Model.Role.GetError, Attr: map[string]any{"Error": res.Error.Error()}}
		}
		savedRole, ret := r.GetByID(role.ID, GetOptions{Preloads: map[string]GetOptions{"Permissions": {}}})
		if !ret.OK {
			return ret
		}
		permissions, ok := model.DefaultRolePermissionMap[savedRole.Name]
		if !ok {
			continue
		}
		for _, permission := range permissions {
			perm, permRet := InitPermissionRepo(r.DB).GetByUniqueField("name", permission)
			if !permRet.OK {
				return permRet
			}
			if slices.ContainsFunc(savedRole.Permissions, func(permission model.Permission) bool {
				return permission.ID == perm.ID
			}) {
				continue
			}
			if ret = AssignPermissionToRole(r.DB, perm, savedRole); !ret.OK {
				return ret
			}
		}
	}
	return model.SuccessRetVal()
}

func (r *RoleRepo) GetFallbackRoleID(excludedRoleIDL ...uint) (uint, model.RetVal) {
	if config.Env.Registration.DefaultGroup != 0 {
		defaultGroup, ret := InitGroupRepo(r.DB).GetByID(config.Env.Registration.DefaultGroup)
		if ret.OK && defaultGroup.RoleID != 0 && !slices.Contains(excludedRoleIDL, defaultGroup.RoleID) {
			return defaultGroup.RoleID, model.SuccessRetVal()
		}
	}

	userRole, ret := r.GetByUniqueField("name", model.UserRoleName)
	if !ret.OK {
		return 0, ret
	}
	if slices.Contains(excludedRoleIDL, userRole.ID) {
		return 0, model.RetVal{
			Msg: i18n.Model.Role.GetError,
			Attr: map[string]any{"Error": "no fallback role available"},
		}
	}
	return userRole.ID, model.SuccessRetVal()
}

func (r *RoleRepo) Delete(idL ...uint) model.RetVal {
	roleL, _, ret := r.List(-1, -1, GetOptions{Conditions: map[string]interface{}{"id": idL}})
	if !ret.OK {
		if ret.Msg != i18n.Model.NotFound {
			return ret
		}
		return model.SuccessRetVal()
	}
	groupRepo := InitGroupRepo(r.DB)
	groupL, _, ret := groupRepo.List(-1, -1, GetOptions{
		Conditions: map[string]any{"role_id": idL},
	})
	if !ret.OK {
		return ret
	}
	fallbackRoleID, ret := r.GetFallbackRoleID(idL...)
	if !ret.OK {
		return ret
	}
	for _, group := range groupL {
		if ret = groupRepo.Update(group.ID, UpdateGroupOptions{RoleID: &fallbackRoleID}); !ret.OK {
			return ret
		}
	}
	for _, role := range roleL {
		if ret = r.Update(role.ID, UpdateRoleOptions{
			Name: new(fmt.Sprintf("%s_deleted_%s", role.Name, utils.RandStr(6))),
		}); !ret.OK {
			return ret
		}
	}
	if res := r.DB.Model(&model.Role{}).Where("id IN ?", idL).Delete(&model.Role{}); res.Error != nil {
		log.Logger.Warningf("Failed to delete Role: %s", res.Error)
		return model.RetVal{Msg: i18n.Model.Role.DeleteError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}
