package repo

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"CBCTF/internel/utils"
	"fmt"
	"gorm.io/gorm"
)

type AdminRepo struct {
	BasicRepo[model.Admin]
}

type CreateAdminOptions struct {
	Name     string
	Password string
	Email    string
	Avatar   model.AvatarURL
	Verified bool
}

func (c CreateAdminOptions) Convert2Model() model.Model {
	return model.Admin{
		Name:     c.Name,
		Password: c.Password,
		Email:    c.Email,
		Avatar:   c.Avatar,
		Verified: c.Verified,
	}
}

type UpdateAdminOptions struct {
	Name     *string
	Password *string
	Email    *string
	Avatar   *model.AvatarURL
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
	if u.Avatar != nil {
		options["avatar"] = *u.Avatar
	}
	if u.Verified != nil {
		options["verified"] = *u.Verified
	}
	return options
}

func InitAdminRepo(tx *gorm.DB) *AdminRepo {
	return &AdminRepo{
		BasicRepo: BasicRepo[model.Admin]{
			DB: tx,
		},
	}
}

func (a *AdminRepo) IsUniqueName(name string) bool {
	_, ok, _ := a.GetByUniqueKey("name", name, GetOptions{Selects: []string{"id"}})
	return !ok
}

func (a *AdminRepo) IsUniqueEmail(email string) bool {
	_, ok, _ := a.GetByUniqueKey("email", email, GetOptions{Selects: []string{"id"}})
	return !ok
}

func (a *AdminRepo) GetByName(name string, optionsL ...GetOptions) (model.Admin, bool, string) {
	options := GetOptions{}
	if len(optionsL) > 0 {
		options = optionsL[0]
	}
	return a.GetByUniqueKey("name", name, options)
}

func (a *AdminRepo) Delete(idL ...uint) (bool, string) {
	for _, id := range idL {
		admin, ok, msg := a.GetByID(id, GetOptions{Selects: []string{"name", "email"}})
		if !ok && msg != i18n.AdminNotFound {
			return false, msg
		}
		deletedName := fmt.Sprintf("%s_deleted_%s", admin.Name, utils.RandStr(6))
		deletedEmail := fmt.Sprintf("%s_deleted_%s", admin.Email, utils.RandStr(6))
		if ok, msg = a.Update(id, UpdateAdminOptions{
			Name:  &deletedName,
			Email: &deletedEmail,
		}); !ok {
			return false, msg
		}
	}
	if res := a.DB.Model(&model.Admin{}).Where("id IN ?", idL).Delete(&model.Admin{}); res.Error != nil {
		log.Logger.Warningf("Failed to delete Admin: %s", res.Error)
		return false, model.Admin{}.DeleteErrorString()
	}
	return true, i18n.Success
}
