package service

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/model"

	"gorm.io/gorm"
)

func ListNotices(tx *gorm.DB, contest model.Contest, form dto.ListModelsForm) ([]model.Notice, int64, model.RetVal) {
	return db.InitNoticeRepo(tx).List(form.Limit, form.Offset, db.GetOptions{
		Conditions: map[string]any{"contest_id": contest.ID},
		Sort:       []string{"id DESC"},
	})
}

func CreateNotice(tx *gorm.DB, contest model.Contest, form dto.CreateNoticeForm) (model.Notice, model.RetVal) {
	return db.InitNoticeRepo(tx).Create(db.CreateNoticeOptions{
		ContestID: contest.ID,
		Title:     form.Title,
		Content:   form.Content,
		Type:      form.Type,
	})
}

func UpdateNotice(tx *gorm.DB, notice model.Notice, form dto.UpdateNoticeForm) model.RetVal {
	return db.InitNoticeRepo(tx).Update(notice.ID, db.UpdateNoticeOptions{
		Title:   form.Title,
		Content: form.Content,
		Type:    form.Type,
	})
}

func DeleteNotice(tx *gorm.DB, notice model.Notice) model.RetVal {
	return db.InitNoticeRepo(tx).Delete(notice.ID)
}
