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

type Admin struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"index:idx_name_deleted,unique;not null" json:"name"`
	Password  string         `gorm:"not null" json:"-"`
	Email     string         `gorm:"index:idx_email_deleted,unique;not null" json:"email"`
	Avatar    string         `json:"-"`
	Verified  bool           `gorm:"default:false" json:"verified"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index;index:idx_name_deleted,unique;index:idx_email_deleted,unique" json:"-"`
}

func (a *Admin) MarshalJSON() ([]byte, error) {
	type Tmp Admin
	return json.Marshal(&struct {
		*Tmp
		Avatar string `json:"avatar"`
	}{
		Tmp:    (*Tmp)(a),
		Avatar: fmt.Sprintf("%s/%s", config.Env.Backend, strings.TrimPrefix(a.Avatar, "/")),
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
