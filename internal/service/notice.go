package service

import (
	"CBCTF/internal/db"
	f "CBCTF/internal/form"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"

	"gorm.io/gorm"
)

func CreateNotice(tx *gorm.DB, contest model.Contest, form f.CreateNoticeForm) (model.Notice, bool, string) {
	notice, ok, msg := db.InitNoticeRepo(tx).Create(db.CreateNoticeOptions{
		ContestID: contest.ID,
		Title:     form.Title,
		Content:   form.Content,
		Type:      form.Type,
	})
	if !ok {
		return model.Notice{}, false, msg
	}
	if ok, msg = db.InitContestRepo(tx).Update(contest.ID, db.UpdateContestOptions{DiffNoticeCount: 1}); !ok {
		return model.Notice{}, false, msg
	}
	return notice, true, i18n.Success
}

func UpdateNotice(tx *gorm.DB, notice model.Notice, form f.UpdateNoticeForm) (bool, string) {
	return db.InitNoticeRepo(tx).Update(notice.ID, db.UpdateNoticeOptions{
		Title:   form.Title,
		Content: form.Content,
		Type:    form.Type,
	})
}

func DeleteNotice(tx *gorm.DB, notice model.Notice) (bool, string) {
	if ok, msg := db.InitNoticeRepo(tx).Delete(notice.ID); !ok {
		return false, msg
	}
	return db.InitContestRepo(tx).Update(notice.ID, db.UpdateContestOptions{DiffNoticeCount: -1})
}
