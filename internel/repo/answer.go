package repo

import (
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
	return &AnswerRepo{Repo: Repo[model.Answer]{DB: tx, Model: "Answer"}}
}

//func (a *AnswerRepo) Create(options CreateAnswerOptions) (model.Answer, bool, string) {
//	answer, err := utils.S2S[model.Answer](options)
//	if err != nil {
//		log.Logger.Warningf("Failed to convert options to model.Answer: %s", err)
//		return model.Answer{}, false, "Options2ModelError"
//	}
//	res := a.DB.Model(&model.Answer{}).Create(&answer)
//	if res.Error != nil {
//		log.Logger.Warningf("Failed to create Answer: %s", res.Error)
//		return model.Answer{}, false, "CreateAnswerError"
//	}
//	return answer, true, "Success"
//}

//func (a *AnswerRepo) getByUniqueKey(key string, value interface{}, preload bool, depth int) (model.Answer, bool, string) {
//	switch key {
//	case "id":
//		value = value.(uint)
//	default:
//		return model.Answer{}, false, "UnsupportedKey"
//	}
//	var answer model.Answer
//	res := a.DB.Model(&model.Answer{}).Where(key+" = ?", value)
//	res = model.GetPreload(res, a.Model, preload, depth).Find(&answer).Limit(1)
//	if res.RowsAffected == 0 {
//		return model.Answer{}, false, "AnswerNotFound"
//	}
//	return answer, true, "Success"
//}

//func (a *AnswerRepo) GetByID(id uint, preload bool, depth int) (model.Answer, bool, string) {
//	return a.getByUniqueKey("id", id, preload, depth)
//}

func (a *AnswerRepo) GetBy2ID(teamID, flagID uint, preload bool, depth int) (model.Answer, bool, string) {
	var answer model.Answer
	res := a.DB.Model(&model.Answer{}).Where("team_id = ? AND flag_id = ?", teamID, flagID)
	res = model.GetPreload(res, a.Model, preload, depth).Find(&answer).Limit(1)
	if res.RowsAffected == 0 {
		return model.Answer{}, false, "AnswerNotFound"
	}
	return answer, true, "Success"
}

func (a *AnswerRepo) Count(flagID uint) (int64, bool, string) {
	var count int64
	res := a.DB.Model(&model.Answer{}).Where("flag_id = ?", flagID).Count(&count)
	if res.Error != nil {
		log.Logger.Warningf("Failed to count Answers: %s", res.Error)
		return 0, false, "CountModelError"
	}
	return count, true, "Success"
}

func (a *AnswerRepo) GetAll(flagID uint, limit, offset int, preload bool, depth int) ([]model.Answer, int64, bool, string) {
	var (
		answers        = make([]model.Answer, 0)
		count, ok, msg = a.Count(flagID)
	)
	if !ok {
		return answers, count, false, msg
	}
	res := a.DB.Model(&model.Answer{}).Where("flag_id = ?", flagID)
	res = model.GetPreload(res, a.Model, preload, depth).Find(&answers).Limit(limit).Offset(offset)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Answers: %s", res.Error)
		return answers, count, false, "GetAnswerError"
	}
	return answers, count, true, "Success"
}

func (a *AnswerRepo) Update(id uint, options UpdateAnswerOptions) (bool, string) {
	var count int
	data := utils.UpdateOptions2Map(options)
	for {
		count++
		if count > 10 {
			log.Logger.Warningf("Failed to update Answer: too many times failed due to optimistic lock")
			return false, "DeadLock"
		}
		answer, ok, msg := a.GetByID(id, false, 0)
		if !ok {
			return ok, msg
		}
		data["version"] = answer.Version + 1
		res := a.DB.Model(&model.Answer{}).Omit("id", "created_at", "updated_at", "deleted_at").
			Where("id = ? AND version = ?", id, answer.Version).Updates(data)
		if res.Error != nil {
			log.Logger.Warningf("Failed to update Answer: %s", res.Error)
			return false, "UpdateAnswerError"
		}
		if res.RowsAffected == 0 {
			continue
		}
		break
	}
	return true, "Success"
}

//func (a *AnswerRepo) Delete(idL ...uint) (bool, string) {
//	res := a.DB.Model(&model.Answer{}).Where("id IN ?", idL).Delete(&model.Answer{})
//	if res.Error != nil {
//		log.Logger.Warningf("Failed to delete Answer: %s", res.Error)
//		return false, "DeleteAnswerError"
//	}
//	return true, "Success"
//}
