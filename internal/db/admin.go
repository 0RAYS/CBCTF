package db

import (
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
)

// countAdmin 统计目前admin数量
func countAdmin() int64 {
	var count int64
	DB.Model(&model.User{}).Where("type = ?", "admin").Count(&count)
	return count
}

// initAdmin 如果没有管理员则创建一个管理员，admin/{random}/admin@0rays.club
func initAdmin() {
	if countAdmin() == 0 {
		CreateAdmin("admin", "admin@0rays.club")
	}
}

// CreateAdmin 创建管理员，本质上为创建用户，然后修改权限
func CreateAdmin(name string, email string) string {
	pwd := utils.RandomString()
	user, ok, msg := CreateUser(name, pwd, email)
	if !ok {
		log.Logger.Warningf("Failed to create admin: %s", msg)
	}
	ok, msg = UpdateUser(user, map[string]interface{}{"type": "admin"})
	if !ok {
		log.Logger.Warningf("Failed to update admin: %s", msg)
	}
	log.Logger.Infof("Create admin: %s/%s/%s", name, pwd, email)
	return pwd
}
