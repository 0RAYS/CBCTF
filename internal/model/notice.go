package model

import (
	f "CBCTF/internal/form"
	"gorm.io/gorm"
	"time"
)

type Notice struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	ContestID uint           `json:"contest_id"`
	Title     string         `json:"title"`
	Content   string         `json:"content"`
	CreatorID uint           `json:"creator_id"`
	CreatedAt time.Time      `json:"created"`
	UpdatedAt time.Time      `json:"updated"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func InitNotice(contestID uint, form f.CreateNoticeForm, creatorID uint) Notice {
	return Notice{
		Title:     form.Title,
		Content:   form.Content,
		ContestID: contestID,
		CreatorID: creatorID,
	}
}
