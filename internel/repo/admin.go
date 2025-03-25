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
	Name     *string
	Password *string
	Email    *string
	Avatar   *string
	Verified *bool
}

func InitAdminRepo(tx *gorm.DB) *AdminRepo {
	return &AdminRepo{Repo: Repo[model.Admin]{DB: tx, Model: "Admin"}}
}

func (a *AdminRepo) IsUniqueName(name string) bool {
	_, ok, _ := a.GetByName(name, false, 0)
	return !ok
}

func (a *AdminRepo) IsUniqueEmail(email string) bool {
	_, ok, _ := a.GetByEmail(email, false, 0)
	return !ok
}

//func (a *AdminRepo) Create(options CreateAdminOptions) (model.Admin, bool, string) {
//	admin, err := utils.S2S[model.Admin](options)
//	if err != nil {
//		log.Logger.Warningf("Failed to convert options to model.Admin: %s", err)
//		return model.Admin{}, false, "Options2ModelError"
//	}
//	if res := a.DB.Model(&model.Admin{}).Create(&admin); res.Error != nil {
//		log.Logger.Warningf("Failed to create Admin: %s", res.Error)
//		return model.Admin{}, false, "CreateAdminError"
//	}
//	return admin, true, "Success"
//}

func (a *AdminRepo) getByUniqueKey(key string, value interface{}, preload bool, depth int) (model.Admin, bool, string) {
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
	res = model.GetPreload(res, a.Model, preload, depth).Find(&admin).Limit(1)
	if res.RowsAffected == 0 {
		return model.Admin{}, false, "AdminNotFound"
	}
	return admin, true, "Success"
}

func (a *AdminRepo) GetByID(id uint, preload bool, depth int) (model.Admin, bool, string) {
	return a.getByUniqueKey("id", id, preload, depth)
}

func (a *AdminRepo) GetByName(name string, preload bool, depth int) (model.Admin, bool, string) {
	return a.getByUniqueKey("name", name, preload, depth)
}

func (a *AdminRepo) GetByEmail(email string, preload bool, depth int) (model.Admin, bool, string) {
	return a.getByUniqueKey("email", email, preload, depth)
}

//func (a *AdminRepo) Count() (int64, bool, string) {
//	var count int64
//	res := a.DB.Model(&model.Admin{}).Count(&count)
//	if res.Error != nil {
//		log.Logger.Warningf("Failed to count Admins: %s", res.Error)
//		return 0, false, "CountModelError"
//	}
//	return count, true, "Success"
//}

//func (a *AdminRepo) GetAll(limit, offset int, preload bool, depth int) ([]model.Admin, int64, bool, string) {
//	var (
//		admins         = make([]model.Admin, 0)
//		count, ok, msg = a.Count()
//	)
//	if !ok {
//		return admins, count, false, msg
//	}
//	res := a.DB.Model(&model.Admin{})
//	res = model.GetPreload(res, a.Model, preload, depth).Find(&admins).Limit(limit).Offset(offset)
//	if res.Error != nil {
//		log.Logger.Warningf("Failed to get Admins: %s", res.Error)
//		return admins, count, false, "GetAdminError"
//	}
//	return admins, count, true, "Success"
//}

func (a *AdminRepo) Update(id uint, options UpdateAdminOptions) (bool, string) {
	var count uint
	data := utils.UpdateOptions2Map(options)
	for {
		count++
		if count > 10 {
			log.Logger.Warningf("Failed to update Admin: too many times failed due to optimistic lock")
			return false, "DeadLock"
		}
		admin, ok, msg := a.GetByID(id, false, 0)
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

//func (a *AdminRepo) Delete(idL ...uint) (bool, string) {
//	res := a.DB.Model(&model.Admin{}).Where("id IN ?", idL).Delete(&model.Admin{})
//	if res.Error != nil {
//		log.Logger.Warningf("Failed to delete Admin: %s", res.Error)
//		return false, "DeleteAdminError"
//	}
//	return true, "Success"
//}
