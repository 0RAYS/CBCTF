package repo

import (
	"CBCTF/internel/log"
	"CBCTF/internel/model"
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

func (s *SubmissionRepo) CountByTeam(teamID uint) (int64, bool, string) {
	var count int64
	res := s.DB.Model(&model.Submission{}).Where("team_id = ?", teamID).Count(&count)
	if res.Error != nil {
		log.Logger.Warningf("Failed to count Submissions: %s", res.Error)
		return 0, false, "CountModelError"
	}
	return count, true, "Success"
}

func (s *SubmissionRepo) GetAllByTeam(teamID uint, limit, offset int, preload bool, depth int) ([]model.Submission, int64, bool, string) {
	var (
		submissions    = make([]model.Submission, 0)
		count, ok, msg = s.CountByTeam(teamID)
	)
	if !ok {
		return submissions, count, false, msg
	}
	res := s.DB.Model(&model.Submission{}).Where("team_id = ?", teamID)
	res = model.GetPreload(res, s.Model, preload, depth).Find(&submissions).Limit(limit).Offset(offset)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Submissions: %s", res.Error)
		return submissions, count, false, "GetSubmissionError"
	}
	return submissions, count, true, "Success"
}

//func (s *SubmissionRepo) Delete(idL ...uint) (bool, string) {
//	res := s.DB.Model(&model.Submission{}).Where("id IN ?", idL).Delete(&model.Submission{})
//	if res.Error != nil {
//		log.Logger.Warningf("Failed to delete Submission: %s", res.Error)
//		return false, "DeleteSubmissionError"
//	}
//	return true, "Success"
//}
