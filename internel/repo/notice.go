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
	Content   string
	Type      string
}

type UpdateNoticeOptions struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
	Type    *string `json:"type"`
}

func InitNoticeRepo(tx *gorm.DB) *NoticeRepo {
	return &NoticeRepo{Repo: Repo[model.Notice]{DB: tx, Model: "Notice"}}
}

func (n *NoticeRepo) Count(contestID uint) (int64, bool, string) {
	var count int64
	res := n.DB.Model(&model.Notice{}).Where("contest_id = ?", contestID).Count(&count)
	if res.Error != nil {
		log.Logger.Warningf("Failed to count Notices: %s", res.Error)
		return 0, false, "CountModelError"
	}
	return count, true, "Success"
}

func (n *NoticeRepo) GetAll(contestID uint, limit, offset int, preloadL ...string) ([]model.Notice, int64, bool, string) {
	var (
		notices        = make([]model.Notice, 0)
		count, ok, msg = n.Count(contestID)
	)
	if !ok {
		return notices, count, false, msg
	}
	res := n.DB.Model(&model.Notice{}).Where("contest_id = ?", contestID)
	res = preload(res, preloadL...).Limit(limit).Offset(offset).Find(&notices)
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
		notice, ok, msg := n.GetByID(id)
		if !ok {
			return ok, msg
		}
		data["version"] = notice.Version + 1
		res := n.DB.Model(&model.Notice{}).Where("id = ? AND version = ?", id, notice.Version).Updates(data)
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
