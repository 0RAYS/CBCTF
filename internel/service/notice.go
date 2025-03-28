package service

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"gorm.io/gorm"
)

func CreateNotice(tx *gorm.DB, contest model.Contest, form f.CreateNoticeForm, adminID uint) (model.Notice, bool, string) {
	repo := db.InitNoticeRepo(tx)
	return repo.Create(db.CreateNoticeOptions{
		ContestID: contest.ID,
		AdminID:   adminID,
		Title:     form.Title,
		Content:   form.Content,
	})
}

func UpdateNotice(tx *gorm.DB, notice model.Notice, form f.UpdateNoticeForm) (bool, string) {
	repo := db.InitNoticeRepo(tx)
	return repo.Update(notice.ID, db.UpdateNoticeOptions{
		Title:   form.Title,
		Content: form.Content,
	})
}

func DeleteNotice(tx *gorm.DB, notice model.Notice) (bool, string) {
	repo := db.InitNoticeRepo(tx)
	return repo.Delete(notice.ID)
}
