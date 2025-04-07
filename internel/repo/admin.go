package repo

import (
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"CBCTF/internel/utils"
	"gorm.io/gorm"
)

type AdminRepo struct {
	Repo[model.Admin]
}

type CreateAdminOptions struct {
	Name     string
	Password string
	Email    string
	Avatar   string
	Verified bool
}

type UpdateAdminOptions struct {
	Name     *string `json:"name"`
	Password *string `json:"password"`
	Email    *string `json:"email"`
	Avatar   *string `json:"avatar"`
	Verified *bool   `json:"verified"`
}

func InitAdminRepo(tx *gorm.DB) *AdminRepo {
	return &AdminRepo{Repo: Repo[model.Admin]{DB: tx, Model: "Admin"}}
}

func (a *AdminRepo) IsUniqueName(name string) bool {
	_, ok, _ := a.GetByName(name, false)
	return !ok
}

func (a *AdminRepo) IsUniqueEmail(email string) bool {
	_, ok, _ := a.GetByEmail(email, false)
	return !ok
}

func (a *AdminRepo) getByUniqueKey(key string, value interface{}, preload bool, nestedL ...string) (model.Admin, bool, string) {
	switch key {
	case "name", "email":
		value = value.(string)
	case "id":
		value = value.(uint)
	default:
		return model.Admin{}, false, "UnsupportedKey"
	}
	var admin model.Admin
	res := a.DB.Model(&model.Admin{}).Where(key+" = ?", value)
	res = model.GetPreload(res, preload, nestedL...).Limit(1).Find(&admin)
	if res.RowsAffected == 0 {
		return model.Admin{}, false, "AdminNotFound"
	}
	return admin, true, "Success"
}

func (a *AdminRepo) GetByID(id uint, preload bool, nestedL ...string) (model.Admin, bool, string) {
	return a.getByUniqueKey("id", id, preload, nestedL...)
}

func (a *AdminRepo) GetByName(name string, preload bool, nestedL ...string) (model.Admin, bool, string) {
	return a.getByUniqueKey("name", name, preload, nestedL...)
}

func (a *AdminRepo) GetByEmail(email string, preload bool, nestedL ...string) (model.Admin, bool, string) {
	return a.getByUniqueKey("email", email, preload, nestedL...)
}

func (a *AdminRepo) Update(id uint, options UpdateAdminOptions) (bool, string) {
	var count uint
	data := utils.UpdateOptions2Map(options)
	for {
		count++
		if count > 10 {
			log.Logger.Warningf("Failed to update Admin: too many times failed due to optimistic lock")
			return false, "DeadLock"
		}
		admin, ok, msg := a.GetByID(id, false)
		if !ok {
			return ok, msg
		}
		data["version"] = admin.Version + 1
		res := a.DB.Model(&model.Admin{}).Omit("id", "created_at", "updated_at", "deleted_at").
			Where("id = ? AND version = ?", id, admin.Version).Updates(data)
		if res.Error != nil {
			log.Logger.Warningf("Failed to update Admin: %s", res.Error)
			return false, "UpdateAdminError"
		}
		if res.RowsAffected == 0 {
			continue
		}
		break
	}
	return true, "Success"
}
