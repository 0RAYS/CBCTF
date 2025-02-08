package model

import (
	"CBCTF/internal/config"
	"CBCTF/internal/utils"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
)

// Admin TODO 由于软删除, 即使数据被删除后 unique 字段仍会受到影响, 有待解决
type Admin struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"unique;not null" json:"name"`
	Password  string         `gorm:"not null" json:"-"`
	Email     string         `gorm:"unique;not null" json:"email"`
	Avatar    string         `json:"-"`
	Verified  bool           `gorm:"default:false" json:"verified"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (m Admin) MarshalJSON() ([]byte, error) {
	type Tmp Admin
	return json.Marshal(&struct {
		Tmp
		Avatar string `json:"avatar"`
	}{
		Tmp:    Tmp(m),
		Avatar: fmt.Sprintf("%s/%s", config.Env.Backend, strings.TrimPrefix(m.Avatar, "/")),
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
