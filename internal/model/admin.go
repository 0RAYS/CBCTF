package model

import (
	"CBCTF/internal/utils"
	"gorm.io/gorm"
	"time"
)

type Admin struct {
	ID        uint           `gorm:"primarykey"`
	Name      string         `gorm:"unique;not null" json:"name"`
	Password  string         `gorm:"not null" json:"-"`
	Email     string         `gorm:"unique;not null" json:"email"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func InitAdmin(name string, password string, email string) Admin {
	return Admin{
		Name:     name,
		Password: utils.HashPassword(password),
		Email:    email,
	}
}
