package db

import (
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"gorm.io/gorm"
)

// CreateAdmin 创建管理员
func CreateAdmin(tx *gorm.DB, name string, password string, email string) (model.Admin, bool, string) {
	if !IsValidEmail(email) {
		return model.Admin{}, false, "InvalidEmail"
	}
	if !IsUniqueName(tx, name, model.Admin{}) {
		return model.Admin{}, false, "AdminNameExists"
	}
	if !IsUniqueEmail(tx, email) {
		return model.Admin{}, false, "EmailExists"
	}
	admin := model.InitAdmin(name, password, email)
	res := tx.Model(&model.Admin{}).Create(&admin)
	if res.Error != nil {
		log.Logger.Warningf("Failed to create Admin: %s", res.Error)
		return model.Admin{}, false, "CreateAdminError"
	}
	return admin, true, "Success"
}

// GetAdminByID 根据 id 获取 model.Admin
func GetAdminByID(tx *gorm.DB, id uint) (model.Admin, bool, string) {
	var admin model.Admin
	res := tx.Model(&model.Admin{}).Where("id = ?", id).Find(&admin).Limit(1)
	if res.RowsAffected != 1 {
		return model.Admin{}, false, "AdminNotFound"
	}
	return admin, true, "Success"
}

func GetAdminByName(tx *gorm.DB, name string) (model.Admin, bool, string) {
	var admin model.Admin
	res := tx.Model(&model.Admin{}).Where("name = ?", name).Find(&admin)
	if res.RowsAffected != 1 {
		return model.Admin{}, false, "AdminNotFound"
	}
	return admin, true, "Success"
}

// GetAdmins 获取所有管理员
func GetAdmins(tx *gorm.DB) ([]model.Admin, int, bool, string) {
	var admins []model.Admin
	res := tx.Model(&model.Admin{}).Find(&admins)
	if res.Error != nil {
		log.Logger.Errorf("Failed to get admins: %s", res.Error)
		return nil, 0, false, "UnknownError"
	}
	return admins, len(admins), true, "Success"
}

// DeleteAdmin 根据 id 删除 model.Admin
func DeleteAdmin(tx *gorm.DB, id uint) (bool, string) {
	res := tx.Model(&model.Admin{}).Where("id = ?", id).Delete(&model.Admin{})
	if res.Error != nil {
		log.Logger.Warningf("Failed to delete Admin: %s", res.Error)
		return false, "DeleteAdminError"
	}
	return true, "Success"
}

// UpdateAdmin 更新管理员, 使用 map 更新属性, 结构体会导致零值未更新, 对字段值的具体要求应当交给上层实现
func UpdateAdmin(tx *gorm.DB, id uint, updateData map[string]interface{}) (bool, string) {
	res := tx.Model(&model.Admin{}).Where("id = ?", id).
		Omit("id", "created_at", "updated_at", "deleted_at").Updates(updateData)
	if res.Error != nil {
		log.Logger.Warningf("Failed to update Admin: %v", res.Error)
		return false, "UpdateAdminError"
	}
	return true, "Success"
}

// VerifyAdmin 验证管理员
func VerifyAdmin(tx *gorm.DB, username string, password string) (model.Admin, bool, string) {
	admin, ok, _ := GetAdminByName(tx, username)
	if !ok {
		return model.Admin{}, false, "NameOrPasswordError"
	}
	if utils.CompareHashAndPassword(admin.Password, password) {
		return admin, true, "Success"
	}
	return model.Admin{}, false, "NameOrPasswordError"
}

// InitAdmin 初始化管理员
func InitAdmin(tx *gorm.DB) {
	var count int64
	tx.Model(&model.Admin{}).Count(&count)
	if count == 0 {
		pwd := utils.UUID()
		CreateAdmin(tx, "admin", pwd, "admin@0rays.club")
		log.Logger.Infof("Init admin: admin/%s/admin@0rays.club", pwd)
	}
}

// ChangePasswordAdmin 修改管理员密码
func ChangePasswordAdmin(tx *gorm.DB, admin model.Admin, oldPassword string, newPassword string) (bool, string) {
	if !utils.CompareHashAndPassword(admin.Password, oldPassword) {
		return false, "PasswordError"
	}
	if utils.CompareHashAndPassword(admin.Password, newPassword) {
		return false, "PasswordSame"
	}
	hash := utils.HashPassword(newPassword)
	if ok, msg := UpdateAdmin(tx, admin.ID, map[string]interface{}{"password": hash}); !ok {
		return false, msg
	}
	return true, "Success"
}
