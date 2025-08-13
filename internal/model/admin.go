package model

import (
	"CBCTF/internal/i18n"
)

// Admin 系统管理员
type Admin struct {
	Name     string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"name"`
	Password string    `gorm:"not null" json:"-"`
	Email    string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Avatar   AvatarURL `json:"avatar"`
	Verified bool      `gorm:"default:false" json:"verified"`
	BasicModel
}

func (a Admin) GetModelName() string {
	return "Admin"
}

func (a Admin) GetVersion() uint {
	return a.Version
}

func (a Admin) GetBasicModel() BasicModel {
	return a.BasicModel
}

func (a Admin) CreateErrorString() string {
	return i18n.CreateAdminError
}

func (a Admin) DeleteErrorString() string {
	return i18n.DeleteAdminError
}

func (a Admin) GetErrorString() string {
	return i18n.GetAdminError
}

func (a Admin) NotFoundErrorString() string {
	return i18n.AdminNotFound
}

func (a Admin) UpdateErrorString() string {
	return i18n.UpdateAdminError
}

func (a Admin) GetUniqueKey() []string {
	return []string{"id", "name", "email"}
}
