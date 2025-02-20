package db

import (
	f "CBCTF/internal/form"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"gorm.io/gorm"
)

func CreateNotice(tx *gorm.DB, contest model.Contest, form f.CreateNoticeForm) (model.Notice, bool, string) {
	notice := model.InitNotice(contest.ID, form)
	res := tx.Model(&model.Notice{}).Create(&notice)
	if res.Error != nil {
		return model.Notice{}, false, "CreateNoticeError"
	}
	return notice, true, "Success"
}

func GetNoticeByID(tx *gorm.DB, id uint) (model.Notice, bool, string) {
	var notice model.Notice
	res := tx.Model(&model.Notice{}).Where("id = ?", id).Find(&notice).Limit(1)
	if res.RowsAffected != 1 {
		return model.Notice{}, false, "NoticeNotFound"
	}
	return notice, true, "Success"
}

func UpdateNotice(tx *gorm.DB, id uint, updateDate map[string]interface{}) (bool, string) {
	res := tx.Model(&model.Notice{}).Where("id = ?", id).
		Omit("id", "created_at", "updated_at", "deleted_at").Updates(updateDate)
	if res.Error != nil {
		return false, "UpdateNoticeError"
	}
	return true, "Success"
}

func DeleteNotice(tx *gorm.DB, id uint) (bool, string) {
	res := tx.Model(&model.Notice{}).Where("id = ?", id).Delete(&model.Notice{})
	if res.Error != nil {
		return false, "DeleteNoticeError"
	}
	return true, "Success"
}

func GetNotices(tx *gorm.DB, limit, offset int, contestID uint) ([]model.Notice, int64, bool, string) {
	var notices []model.Notice
	var count int64
	res := tx.Model(&model.Notice{}).Where("contest_id = ?", contestID)
	if res.Count(&count).Error != nil {
		log.Logger.Warningf("Failed to count notices: %s", res.Error)
		return []model.Notice{}, 0, false, "UnknownError"
	}
	limit, offset = utils.TidyPaginate(int(count), limit, offset)
	res = res.Limit(limit).Offset(offset).Find(&notices)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get notices: %s", res.Error)
		return []model.Notice{}, 0, false, "GetNoticesError"
	}
	return notices, count, true, "Success"
}
