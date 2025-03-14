package db

import (
	f "CBCTF/internal/form"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"gorm.io/gorm"
)

// CreateNotice 创建公告
func CreateNotice(tx *gorm.DB, contest model.Contest, form f.CreateNoticeForm, creatorID uint) (model.Notice, bool, string) {
	notice := model.InitNotice(contest.ID, form, creatorID)
	res := tx.Model(&model.Notice{}).Create(&notice)
	if res.Error != nil {
		return model.Notice{}, false, "CreateNoticeError"
	}
	return notice, true, "Success"
}

// GetNoticeByID 根据 ID 获取公告
func GetNoticeByID(tx *gorm.DB, id uint) (model.Notice, bool, string) {
	var notice model.Notice
	res := tx.Model(&model.Notice{}).Where("id = ?", id).Find(&notice).Limit(1)
	if res.RowsAffected != 1 {
		return model.Notice{}, false, "NoticeNotFound"
	}
	return notice, true, "Success"
}

// UpdateNotice 更新公告, 使用 map 更新属性, 结构体会导致零值未更新, 对字段值的具体要求应当交给上层实现
func UpdateNotice(tx *gorm.DB, id uint, updateDate map[string]interface{}) (bool, string) {
	var count int
	for {
		count++
		if count > 10 {
			log.Logger.Warningf("Failed too many times to update user due to optimistic lock")
			return false, "FailedTooManyTimes"
		}
		var notice model.Notice
		res := tx.Model(&model.Notice{}).Where("id = ?", id).Find(&notice).Limit(1)
		if res.RowsAffected != 1 {
			return false, "NoticeNotFound"
		}
		res = tx.Model(&notice).Updates(updateDate)
		if res.Error != nil {
			log.Logger.Warningf("Failed to update Notice: %s", res.Error)
			return false, "UpdateNoticeError"
		}
		if res.RowsAffected == 0 {
			log.Logger.Debug("Failed to update notice due to optimistic lock")
			continue
		}
		break
	}
	return true, "Success"
}

// DeleteNotice 删除公告
func DeleteNotice(tx *gorm.DB, id uint) (bool, string) {
	res := tx.Model(&model.Notice{}).Where("id = ?", id).Delete(&model.Notice{})
	if res.Error != nil {
		return false, "DeleteNoticeError"
	}
	return true, "Success"
}

// GetNotices 获取公告
func GetNotices(tx *gorm.DB, limit, offset int, contestID uint) ([]model.Notice, int64, bool, string) {
	var notices []model.Notice
	var count int64
	res := tx.Model(&model.Notice{}).Where("contest_id = ?", contestID)
	if res.Count(&count).Error != nil {
		log.Logger.Warningf("Failed to count notices: %s", res.Error)
		return make([]model.Notice, 0), 0, false, "UnknownError"
	}
	limit, offset = utils.TidyPaginate(int(count), limit, offset)
	res = res.Limit(limit).Offset(offset).Find(&notices)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get notices: %s", res.Error)
		return make([]model.Notice, 0), 0, false, "GetNoticesError"
	}
	return notices, count, true, "Success"
}
