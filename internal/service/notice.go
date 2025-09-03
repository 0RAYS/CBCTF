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
	if ok, msg = db.InitContestRepo(tx).DiffUpdate(contest.ID, db.DiffUpdateContestOptions{NoticeCount: 1}); !ok {
		return model.Notice{}, false, msg
	}
	return notice, true, i18n.Success
}

func DeleteNotice(tx *gorm.DB, notice model.Notice) (bool, string) {
	if ok, msg := db.InitNoticeRepo(tx).Delete(notice.ID); !ok {
		return false, msg
	}
	return db.InitContestRepo(tx).DiffUpdate(notice.ID, db.DiffUpdateContestOptions{NoticeCount: -1})
}
