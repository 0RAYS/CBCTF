package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"fmt"

	"gorm.io/gorm"
)

type AdminRepo struct {
	BaseRepo[model.Admin]
}

type CreateAdminOptions struct {
	Name     string
	Password string
	Email    string
	Picture  model.FileURL
	Verified bool
}

func (c CreateAdminOptions) Convert2Model() model.Model {
	return model.Admin{
		Name:     c.Name,
		Password: c.Password,
		Email:    c.Email,
		Picture:  c.Picture,
		Verified: c.Verified,
	}
}

type UpdateAdminOptions struct {
	Name     *string
	Password *string
	Email    *string
	Picture  *model.FileURL
	Verified *bool
}

func (u UpdateAdminOptions) Convert2Map() map[string]any {
	options := make(map[string]any)
	if u.Name != nil {
		options["name"] = *u.Name
	}
	if u.Password != nil {
		options["password"] = *u.Password
	}
	if u.Email != nil {
		options["email"] = *u.Email
	}
	if u.Picture != nil {
		options["picture"] = *u.Picture
	}
	if u.Verified != nil {
		options["verified"] = *u.Verified
	}
	return options
}

func InitAdminRepo(tx *gorm.DB) *AdminRepo {
	return &AdminRepo{
		BaseRepo: BaseRepo[model.Admin]{
			DB: tx,
		},
	}
}

func (a *AdminRepo) InitAdmin() model.RetVal {
	count, ret := a.Count()
	if !ret.OK {
		return ret
	}
	if count == 0 {
		pwd := utils.UUID()
		_, ret = a.Create(CreateAdminOptions{
			Name:     "admin",
			Password: utils.HashPassword(pwd),
			Email:    "admin@0rays.club",
		})
		if !ret.OK {
			return ret
		}
		log.Logger.Infof("Init Admin: Admin{ name: admin, password: %s, email: admin@0rays.club}", pwd)
	}
	return model.SuccessRetVal()
}

func (a *AdminRepo) Delete(idL ...uint) model.RetVal {
	adminL, _, ret := a.List(-1, -1, GetOptions{Conditions: map[string]any{"id": idL}})
	if !ret.OK {
		if ret.Msg != i18n.Model.NotFound {
			return ret
		}
		return model.SuccessRetVal()
	}
	for _, admin := range adminL {
		if ret = a.Update(admin.ID, UpdateAdminOptions{
			Name:  new(fmt.Sprintf("%s_deleted_%s", admin.Name, utils.RandStr(6))),
			Email: new(fmt.Sprintf("%s_deleted_%s", admin.Email, utils.RandStr(6))),
		}); !ret.OK {
			return ret
		}
	}
	if res := a.DB.Model(&model.Admin{}).Where("id IN ?", idL).Delete(&model.Admin{}); res.Error != nil {
		log.Logger.Warningf("Failed to delete Admin: %s", res.Error)
		return model.RetVal{Msg: i18n.Model.DeleteError, Attr: map[string]any{"Model": model.Admin{}.ModelName(), "Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}
