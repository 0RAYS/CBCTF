package model

import (
	"CBCTF/internal/form"
	"gorm.io/gorm"
	"time"
)

type Notice struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	ContestID uint           `json:"contest_id"`
	Title     string         `json:"title"`
	Content   string         `json:"content"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func InitNotice(contestID uint, f form.CreateNoticeForm) Notice {
	return Notice{
		Title:     f.Title,
		Content:   f.Content,
		ContestID: contestID,
	}
}
