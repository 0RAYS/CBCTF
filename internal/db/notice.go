package db

import (
	"CBCTF/internal/model"

	"gorm.io/gorm"
)

type NoticeRepo struct {
	BasicRepo[model.Notice]
}

type CreateNoticeOptions struct {
	ContestID uint
	Title     string
	Content   string
	Type      string
}

func (c CreateNoticeOptions) Convert2Model() model.Model {
	return model.Notice{
		ContestID: c.ContestID,
		Title:     c.Title,
		Content:   c.Content,
		Type:      c.Type,
	}
}

type UpdateNoticeOptions struct {
	Title   *string
	Content *string
	Type    *string
}

func (u UpdateNoticeOptions) Convert2Map() map[string]any {
	options := make(map[string]any)
	if u.Title != nil {
		options["title"] = *u.Title
	}
	if u.Content != nil {
		options["content"] = *u.Content
	}
	if u.Type != nil {
		options["type"] = *u.Type
	}
	return options
}

func InitNoticeRepo(tx *gorm.DB) *NoticeRepo {
	return &NoticeRepo{
		BasicRepo: BasicRepo[model.Notice]{
			DB: tx,
		},
	}
}
