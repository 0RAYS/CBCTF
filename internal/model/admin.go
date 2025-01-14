package model

import (
	"CBCTF/internal/utils"
	"gorm.io/gorm"
)

type Admin struct {
	ID       uint   `gorm:"primarykey"`
	Name     string `gorm:"unique;not null" json:"name"`
	Password string `gorm:"not null" json:"-"`
	Email    string `gorm:"unique;not null" json:"email"`
	gorm.Model
}

func InitAdmin(name string, password string, email string) Admin {
	return Admin{
		Name:     name,
		Password: utils.HashPassword(password),
		Email:    email,
	}
}
