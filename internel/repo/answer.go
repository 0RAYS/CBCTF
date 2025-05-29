package repo

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"CBCTF/internel/utils"
	"gorm.io/gorm"
)

type AnswerRepo struct {
	Repo[model.Answer]
}

type CreateAnswerOptions struct {
	TeamID uint
	FlagID uint
	Value  string
	Solved bool
}

type UpdateAnswerOptions struct {
	Value  *string `json:"value"`
	Solved *bool   `json:"solved"`
}

func InitAnswerRepo(tx *gorm.DB) *AnswerRepo {
	return &AnswerRepo{
		Repo: Repo[model.Answer]{
			DB: tx, Model: "Answer",
			CreateError:   i18n.CreateAnswerError,
			DeleteError:   i18n.DeleteAnswerError,
			GetError:      i18n.GetAnswerError,
			NotFoundError: i18n.AnswerNotFound,
		},
	}
}

func (a *AnswerRepo) GetBy2ID(teamID, flagID uint, preloadL ...string) (model.Answer, bool, string) {
	var answer model.Answer
	res := a.DB.Model(&model.Answer{}).Where("team_id = ? AND flag_id = ?", teamID, flagID)
	res = preload(res, preloadL...).Limit(1).Find(&answer)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Answer: %s", res.Error)
		return model.Answer{}, false, i18n.GetAnswerError
	}
	if res.RowsAffected == 0 {
		return model.Answer{}, false, i18n.AnswerNotFound
	}
	return answer, true, i18n.Success
}

func (a *AnswerRepo) Count(flagID uint) (int64, bool, string) {
	var count int64
	res := a.DB.Model(&model.Answer{}).Where("flag_id = ?", flagID).Count(&count)
	if res.Error != nil {
		log.Logger.Warningf("Failed to count Answers: %s", res.Error)
		return 0, false, i18n.CountModelError
	}
	return count, true, i18n.Success
}

func (a *AnswerRepo) GetAll(flagID uint, limit, offset int, preloadL ...string) ([]model.Answer, int64, bool, string) {
	var (
		answers        = make([]model.Answer, 0)
		count, ok, msg = a.Count(flagID)
	)
	if !ok {
		return answers, count, false, msg
	}
	res := a.DB.Model(&model.Answer{}).Where("flag_id = ?", flagID)
	res = preload(res, preloadL...).Limit(limit).Offset(offset).Find(&answers)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Answers: %s", res.Error)
		return answers, count, false, i18n.GetAnswerError
	}
	return answers, count, true, i18n.Success
}

func (a *AnswerRepo) Update(id uint, options UpdateAnswerOptions) (bool, string) {
	var count int
	data := utils.UpdateOptions2Map(options)
	for {
		count++
		if count > 10 {
			log.Logger.Warningf("Failed to update Answer: too many times failed due to optimistic lock")
			return false, i18n.DeadLock
		}
		answer, ok, msg := a.GetByID(id)
		if !ok {
			return ok, msg
		}
		data["version"] = answer.Version + 1
		res := a.DB.Model(&model.Answer{}).Where("id = ? AND version = ?", id, answer.Version).Updates(data)
		if res.Error != nil {
			log.Logger.Warningf("Failed to update Answer: %s", res.Error)
			return false, i18n.UpdateAnswerError
		}
		if res.RowsAffected == 0 {
			continue
		}
		break
	}
	return true, i18n.Success
}
