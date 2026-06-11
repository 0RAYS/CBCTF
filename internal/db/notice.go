package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"

	"gorm.io/gorm"
)

type NoticeRepo struct {
	BaseRepo[model.Notice]
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
		BaseRepo: BaseRepo[model.Notice]{
			DB: tx,
		},
	}
}

func (n *NoticeRepo) DeleteByContestID(contestIDL ...uint) model.RetVal {
	if len(contestIDL) == 0 {
		return model.SuccessRetVal()
	}
	var noticeIDL []uint
	if res := n.DB.Model(&model.Notice{}).Where("contest_id IN ?", contestIDL).Pluck("id", &noticeIDL); res.Error != nil {
		log.Logger.Warningf("Failed to get Notices by contest IDs %v: %s", contestIDL, res.Error)
		return model.RetVal{Msg: i18n.Model.Notice.DeleteError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return n.Delete(noticeIDL...)
}
