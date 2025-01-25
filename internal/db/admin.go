package db

import (
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"context"
)

// CreateAdmin 创建管理员
func CreateAdmin(ctx context.Context, name string, password string, email string) (model.Admin, bool, string) {
	if !IsValidEmail(email) {
		return model.Admin{}, false, "InvalidEmail"
	}
	if !IsUniqueName(name, model.Admin{}) {
		return model.Admin{}, false, "AdminNameExists"
	}
	if !IsUniqueEmail(email) {
		return model.Admin{}, false, "EmailExists"
	}
	admin := model.InitAdmin(name, password, email)
	res := DB.WithContext(ctx).Model(&model.Admin{}).Create(&admin)
	if res.Error != nil {
		log.Logger.Errorf("Failed to create Admin: %s", res.Error.Error())
		return model.Admin{}, false, "CreateAdminError"
	}
	return admin, true, "Success"
}

func GetAdminByID(ctx context.Context, id uint) (model.Admin, bool, string) {
	var admin model.Admin
	res := DB.WithContext(ctx).Model(&model.Admin{}).Where("id = ?", id).Find(&admin)
	if res.RowsAffected != 1 {
		return model.Admin{}, false, "AdminNotFound"
	}
	return admin, true, "Success"
}

// DeleteAdmin 根据 id 删除 model.Admin
func DeleteAdmin(ctx context.Context, id uint) (bool, string) {
	res := DB.WithContext(ctx).Model(&model.Admin{}).Where("id = ?", id).Delete(&model.Admin{})
	if res.Error != nil {
		log.Logger.Errorf("Failed to delete Admin: %s", res.Error.Error())
		return false, "DeleteAdminError"
	}
	return true, "Success"
}

// UpdateAdmin 更新管理员, 使用 map 更新属性, 结构体会导致零值未更新, 对字段值的具体要求应当交给上层实现
func UpdateAdmin(ctx context.Context, id uint, updateData map[string]interface{}) (bool, string) {
	res := DB.WithContext(ctx).Model(&model.Admin{}).Where("id = ?", id).
		Omit("id", "created_at", "updated_at", "deleted_at").Updates(updateData)
	if res.Error != nil {
		log.Logger.Warningf("Failed to update Admin: %v", res.Error.Error())
		return false, "UpdateAdminError"
	}
	return true, "Success"
}

func VerifyAdmin(ctx context.Context, username string, password string) (model.Admin, bool, string) {
	var admin model.Admin
	res := DB.WithContext(ctx).Model(&model.Admin{}).Where("name = ? OR email = ?", username, username).
		Find(&admin).Limit(1)
	if res.RowsAffected != 1 {
		return model.Admin{}, false, "NameOrPasswordError"
	}
	if utils.CompareHashAndPassword(admin.Password, password) {
		return admin, true, "Success"
	}
	return model.Admin{}, false, "NameOrPasswordError"
}

func InitAdmin() {
	var count int64
	DB.Model(&model.Admin{}).Count(&count)
	if count == 0 {
		pwd := utils.RandomString()
		CreateAdmin(context.Background(), "admin", pwd, "admin@0rays.club")
		log.Logger.Infof("Init admin: admin/%s/admin@0rays.club", pwd)
	}
}

func ChangePasswordAdmin(ctx context.Context, id uint, oldPassword string, newPassword string) (bool, string) {
	admin, ok, msg := GetAdminByID(ctx, id)
	if !ok {
		return false, msg
	}
	if !utils.CompareHashAndPassword(admin.Password, oldPassword) {
		return false, "PasswordError"
	}
	if utils.CompareHashAndPassword(admin.Password, newPassword) {
		return false, "PasswordSame"
	}
	hash := utils.HashPassword(newPassword)
	if ok, msg := UpdateAdmin(ctx, id, map[string]interface{}{"password": hash}); !ok {
		return false, msg
	}
	return true, "Success"
}

func CountAdmins(ctx context.Context) int64 {
	var count int64
	DB.WithContext(ctx).Model(&model.Admin{}).Count(&count)
	return count
}

func GetAdmins(ctx context.Context, limit int, offset int) ([]model.Admin, int, bool, string) {
	if limit <= 0 {
		limit = -1
	}
	if offset <= 0 {
		offset = -1
	}
	var admins []model.Admin
	res := DB.WithContext(ctx).Model(&model.Admin{}).Limit(limit).Offset(offset).Find(&admins)
	if res.Error != nil {
		log.Logger.Errorf("Failed to get admins: %s", res.Error.Error())
		return nil, 0, false, "UnknownError"
	}
	return admins, len(admins), true, "Success"
}
