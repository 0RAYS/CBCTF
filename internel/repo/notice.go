package repo

import (
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"CBCTF/internel/utils"
	"gorm.io/gorm"
)

type NoticeRepo struct {
	Repo[model.Notice]
}

type CreateNoticeOptions struct {
	ContestID uint
	AdminID   uint
	Title     string
	Contest   string
}

type UpdateNoticeOptions struct {
	Title   *string
	Contest *string
}

func InitNoticeRepo(tx *gorm.DB) *NoticeRepo {
	return &NoticeRepo{Repo: Repo[model.Notice]{DB: tx, Model: "Notice"}}
}

//func (n *NoticeRepo) Create(options CreateNoticeOptions) (model.Notice, bool, string) {
//	notice, err := utils.S2S[model.Notice](options)
//	if err != nil {
//		log.Logger.Warningf("Failed to convert options to model.Notice: %s", err)
//		return model.Notice{}, false, "Options2ModelError"
//	}
//	res := n.DB.Model(&model.Notice{}).Create(&notice)
//	if res.Error != nil {
//		log.Logger.Warningf("Failed to create Notice: %s", res.Error)
//		return model.Notice{}, false, "CreateNoticeError"
//	}
//	return notice, true, "Success"
//}

//func (n *NoticeRepo) getByUniqueKey(key string, value interface{}, preload bool, depth int) (model.Notice, bool, string) {
//	switch key {
//	case "id":
//		value = value.(uint)
//	default:
//		return model.Notice{}, false, "UnsupportedKey"
//	}
//	var notice model.Notice
//	res := n.DB.Model(&model.Notice{}).Where(key+" = ?", value)
//	res = model.GetPreload(res, model.Notice{}, preload, depth).Find(&notice).Limit(1)
//	if res.RowsAffected == 0 {
//		return model.Notice{}, false, "NoticeNotFound"
//	}
//	return notice, true, "Success"
//}

//func (n *NoticeRepo) GetByID(id uint, preload bool, depth int) (model.Notice, bool, string) {
//	return n.getByUniqueKey("id", id, preload, depth)
//}

func (n *NoticeRepo) Count(contestID uint) (int64, bool, string) {
	var count int64
	res := n.DB.Model(&model.Notice{}).Where("contest_id = ?", contestID).Count(&count)
	if res.Error != nil {
		log.Logger.Warningf("Failed to count Notices: %s", res.Error)
		return 0, false, "CountModelError"
	}
	return count, true, "Success"
}

func (n *NoticeRepo) GetAll(contestID uint, limit, offset int, preload bool, depth int) ([]model.Notice, int64, bool, string) {
	var (
		notices        = make([]model.Notice, 0)
		count, ok, msg = n.Count(contestID)
	)
	if !ok {
		return notices, count, false, msg
	}
	res := n.DB.Model(&model.Notice{}).Where("contest_id = ?", contestID)
	res = model.GetPreload(res, n.Model, preload, depth).Find(&notices).Limit(limit).Offset(offset)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Notices: %s", res.Error)
		return notices, 0, false, "GetNoticeError"
	}
	return notices, count, true, "Success"
}

func (n *NoticeRepo) Update(id uint, options UpdateNoticeOptions) (bool, string) {
	var count int
	data := utils.UpdateOptions2Map(options)
	for {
		count++
		if count > 10 {
			log.Logger.Warningf("Failed to update Notice: too many times failed due to optimistic lock")
			return false, "DeadLock"
		}
		notice, ok, msg := n.GetByID(id, false, 0)
		if !ok {
			return ok, msg
		}
		data["version"] = notice.Version + 1
		res := n.DB.Model(&model.Notice{}).Omit("id", "created_at", "updated_at", "deleted_at").
			Where("id = ? AND version = ?", id, notice.Version).Updates(data)
		if res.Error != nil {
			log.Logger.Warningf("Failed to update Notice: %s", res.Error)
			return false, "UpdateNoticeError"
		}
		if res.RowsAffected == 0 {
			continue
		}
		break
	}
	return true, "Success"
}

//func (n *NoticeRepo) Delete(idL ...uint) (bool, string) {
//	res := n.DB.Model(&model.Notice{}).Where("id IN ?", idL).Delete(&model.Notice{})
//	if res.Error != nil {
//		log.Logger.Warningf("Failed to delete Notice: %s", res.Error)
//		return false, "DeleteNoticeError"
//	}
//	return true, "Success"
//}
