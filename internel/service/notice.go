package service

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"gorm.io/gorm"
)

func GetNotices(tx *gorm.DB, contest model.Contest, form f.GetModelsForm) ([]model.Notice, int64, bool, string) {
	repo := db.InitNoticeRepo(tx)
	return repo.GetAll(contest.ID, form.Limit, form.Offset, false, 0)
}
