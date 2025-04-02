package repo

import (
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"CBCTF/internel/utils"
	"gorm.io/gorm"
)

type SubmissionRepo struct {
	Repo[model.Submission]
}

type CreateSubmissionOptions struct {
	UsageID     uint
	ContestID   uint
	ChallengeID string
	TeamID      uint
	UserID      uint
	Value       string
	Solved      bool
	Score       float64
}

type UpdateSubmissionOptions struct {
	Solved *bool `json:"solved"`
}

func InitSubmissionRepo(tx *gorm.DB) *SubmissionRepo {
	return &SubmissionRepo{Repo: Repo[model.Submission]{DB: tx, Model: "Submission"}}
}

//func (s *SubmissionRepo) Create(options CreateSubmissionOptions) (model.Submission, bool, string) {
//	submission, err := utils.S2S[model.Submission](options)
//	if err != nil {
//		log.Logger.Warningf("Failed to conver options to model.Submission: %s", err)
//		return model.Submission{}, false, "Options2ModelError"
//	}
//	if res := s.DB.Model(&model.Submission{}).Create(&submission); res.Error != nil {
//		log.Logger.Warningf("Failed to create Submission: %s", res.Error)
//		return model.Submission{}, false, "CreateSubmissionError"
//	}
//	return submission, true, "Success"
//}

//func (s *SubmissionRepo) getByUniqueKey(key string, value interface{}, preload bool, depth int) (model.Submission, bool, string) {
//	switch key {
//	case "id":
//		value = value.(uint)
//	default:
//		return model.Submission{}, false, "UnsupportedKey"
//	}
//	var submission model.Submission
//	res := s.DB.Model(&model.Submission{}).Where(key+" = ?", value)
//	res = model.GetPreload(res, model.Notice{}, preload, depth).Find(&submission).Limit(1)
//	if res.RowsAffected == 0 {
//		return model.Submission{}, false, "SubmissionNotFound"
//	}
//	return submission, true, "Success"
//}

//func (s *SubmissionRepo) GetByID(id uint, preload bool, depth int) (model.Submission, bool, string) {
//	return s.getByUniqueKey("id", id, preload, depth)
//}

func (s *SubmissionRepo) CountByKeyID(key string, id uint, solved bool) (int64, bool, string) {
	var count int64
	res := s.DB.Model(&model.Submission{}).Where(key+" = ?", id)
	if solved {
		res = res.Where("solved = ?", true)
	}
	if res.Error != nil {
		log.Logger.Warningf("Failed to count Submissions: %s", res.Error)
		return 0, false, "CountModelError"
	}
	return count, true, "Success"
}

func (s *SubmissionRepo) GetAllByKeyID(key string, id uint, limit, offset int, preload bool, depth int, solved bool) ([]model.Submission, int64, bool, string) {
	var (
		submissions    = make([]model.Submission, 0)
		count, ok, msg = s.CountByKeyID(key, id, solved)
	)
	if !ok {
		return submissions, count, false, msg
	}
	res := s.DB.Model(&model.Submission{}).Where(key+" = ?", id)
	if solved {
		res = res.Where("solved = ?", true)
	}
	res = model.GetPreload(res, s.Model, preload, depth).Find(&submissions).Limit(limit).Offset(offset)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Submissions: %s", res.Error)
		return submissions, count, false, "GetSubmissionError"
	}
	return submissions, count, true, "Success"
}

func (s *SubmissionRepo) Update(id uint, options UpdateSubmissionOptions) (bool, string) {
	var count uint
	data := utils.UpdateOptions2Map(options)
	for {
		count++
		if count > 10 {
			log.Logger.Warningf("Failed to update Submission: too many times failed due to optimistic lock")
			return false, "DeadLock"
		}
		submission, ok, msg := s.GetByID(id, false, 0)
		if !ok {
			return ok, msg
		}
		data["version"] = submission.Version + 1
		res := s.DB.Model(&model.Submission{}).Omit("id", "created_at", "updated_at", "deleted_at").
			Where("id = ? AND version = ?", id, submission.Version).Updates(data)
		if res.Error != nil {
			log.Logger.Warningf("Failed to update Submission: %s", res.Error)
			return false, "UpdateSubmissionError"
		}
		if res.RowsAffected == 0 {
			continue
		}
		break
	}
	return true, "Success"
}

//func (s *SubmissionRepo) Delete(idL ...uint) (bool, string) {
//	res := s.DB.Model(&model.Submission{}).Where("id IN ?", idL).Delete(&model.Submission{})
//	if res.Error != nil {
//		log.Logger.Warningf("Failed to delete Submission: %s", res.Error)
//		return false, "DeleteSubmissionError"
//	}
//	return true, "Success"
//}
