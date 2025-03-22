package model

import (
	"CBCTF/internal/config"
	"CBCTF/internal/utils"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/plugin/optimisticlock"
	"strings"
	"time"
)

type Admin struct {
	ID        uint                   `gorm:"primarykey" json:"id"`
	Name      string                 `gorm:"index:idx_name_deleted,unique;not null" json:"name"`
	Password  string                 `gorm:"not null" json:"-"`
	Email     string                 `gorm:"index:idx_email_deleted,unique;not null" json:"email"`
	Avatar    string                 `json:"-"`
	Verified  bool                   `gorm:"default:false" json:"verified"`
	CreatedAt time.Time              `json:"-"`
	UpdatedAt time.Time              `json:"-"`
	DeletedAt gorm.DeletedAt         `gorm:"index;index:idx_name_deleted,unique;index:idx_email_deleted,unique" json:"-"`
	Version   optimisticlock.Version `json:"-" gorm:"default:1"`
}

// MarshalJSON 重写 MarshalJSON 方法, 使其返回完整的 URL
func (a Admin) MarshalJSON() ([]byte, error) {
	avatar := ""
	if strings.TrimPrefix(a.Avatar, "/") != "" {
		avatar = fmt.Sprintf("%s/%s", config.Env.Backend, strings.TrimPrefix(a.Avatar, "/"))
	}
	type Tmp Admin
	return json.Marshal(&struct {
		Tmp
		Avatar string `json:"avatar"`
	}{
		Tmp:    Tmp(a),
		Avatar: avatar,
	})
}

func InitAdmin(name string, password string, email string) Admin {
	return Admin{
		Name:     name,
		Password: utils.HashPassword(password),
		Email:    email,
		Avatar:   "",
		Verified: false,
	}
}
