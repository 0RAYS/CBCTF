package repo

import (
	"CBCTF/internel/i18n"
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
	FlagID      uint
	Value       string
	Solved      bool
	Score       float64
}

type UpdateSubmissionOptions struct {
	Solved *bool    `json:"solved"`
	Score  *float64 `json:"score"`
}

func InitSubmissionRepo(tx *gorm.DB) *SubmissionRepo {
	return &SubmissionRepo{Repo: Repo[model.Submission]{DB: tx, Model: "Submission"}}
}

func (s *SubmissionRepo) CountByKeyID(key string, id uint, solved bool) (int64, bool, string) {
	var count int64
	res := s.DB.Model(&model.Submission{}).Where(key+" = ?", id)
	if solved {
		res = res.Where("solved = ?", true)
	}
	res = res.Count(&count)
	if res.Error != nil {
		log.Logger.Warningf("Failed to count Submissions: %s", res.Error)
		return 0, false, i18n.CountModelError
	}
	return count, true, i18n.Success
}

func (s *SubmissionRepo) GetByKeyID(key string, id uint, limit, offset int, solved bool, preloadL ...string) ([]model.Submission, int64, bool, string) {
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
	res = preload(res, preloadL...).Limit(limit).Offset(offset).Find(&submissions)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Submissions: %s", res.Error)
		return submissions, count, false, i18n.GetSubmissionError
	}
	return submissions, count, true, i18n.Success
}

func (s *SubmissionRepo) Update(id uint, options UpdateSubmissionOptions) (bool, string) {
	var count uint
	data := utils.UpdateOptions2Map(options)
	for {
		count++
		if count > 10 {
			log.Logger.Warningf("Failed to update Submission: too many times failed due to optimistic lock")
			return false, i18n.DeadLock
		}
		submission, ok, msg := s.GetByID(id)
		if !ok {
			return ok, msg
		}
		data["version"] = submission.Version + 1
		res := s.DB.Model(&model.Submission{}).Where("id = ? AND version = ?", id, submission.Version).Updates(data)
		if res.Error != nil {
			log.Logger.Warningf("Failed to update Submission: %s", res.Error)
			return false, i18n.UpdateSubmissionError
		}
		if res.RowsAffected == 0 {
			continue
		}
		break
	}
	return true, i18n.Success
}
