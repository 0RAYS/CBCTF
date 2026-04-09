package service

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/model"
	"CBCTF/internal/view"

	"gorm.io/gorm"
)

func BuildChallengeView(tx *gorm.DB, challenge model.Challenge) view.ChallengeView {
	result := view.ChallengeView{
		Challenge: challenge,
		Flags:     make([]view.ChallengeFlagView, 0),
	}
	if challenge.Type != model.PodsChallengeType {
		for _, flag := range challenge.ChallengeFlags {
			result.Flags = append(result.Flags, view.ChallengeFlagView{
				ID:    flag.ID,
				Value: flag.Value,
			})
		}
	} else {
		result.DockerCompose = Template2Yaml(challenge.Template, challenge.ChallengeFlags)
	}

	file, _ := db.InitFileRepo(tx).Get(db.GetOptions{
		Conditions: map[string]any{
			"model":    model.ModelName(challenge),
			"model_id": challenge.ID,
			"type":     model.ChallengeFileType,
		},
	})
	result.FileName = file.Filename
	return result
}

func BuildChallengeViews(tx *gorm.DB, challenges []model.Challenge) []view.ChallengeView {
	views := make([]view.ChallengeView, 0, len(challenges))
	for _, challenge := range challenges {
		views = append(views, BuildChallengeView(tx, challenge))
	}
	return views
}

func GetChallengeView(tx *gorm.DB, challenge model.Challenge) view.ChallengeView {
	return BuildChallengeView(tx, challenge)
}

func GetSimpleChallengeView(challenge model.Challenge) view.SimpleChallengeView {
	return view.SimpleChallengeView{Challenge: challenge}
}

func ListChallengeViews(tx *gorm.DB, form dto.GetChallengesForm) ([]view.ChallengeView, int64, model.RetVal) {
	challenges, count, ret := GetChallenges(tx, form)
	if !ret.OK {
		return nil, 0, ret
	}
	return BuildChallengeViews(tx, challenges), count, model.SuccessRetVal()
}

func ListChallengesNotInContest(tx *gorm.DB, contest model.Contest, form dto.GetChallengesForm) ([]view.SimpleChallengeView, int64, model.RetVal) {
	challenges, count, ret := db.InitChallengeRepo(tx).ListChallengesNotInContest(
		contest.ID,
		form.Limit,
		form.Offset,
		form.Name,
		form.Description,
		form.Category,
		form.Type,
	)
	if !ret.OK {
		return nil, 0, ret
	}
	views := make([]view.SimpleChallengeView, 0, len(challenges))
	for _, challenge := range challenges {
		views = append(views, view.SimpleChallengeView{Challenge: challenge})
	}
	return views, count, model.SuccessRetVal()
}

func ListChallengeCategories(tx *gorm.DB, form dto.GetCategoriesForm) ([]string, model.RetVal) {
	return db.InitChallengeRepo(tx).ListCategories(form.Type)
}

func CreateChallengeWithTransaction(tx *gorm.DB, form dto.CreateChallengeForm) (model.Challenge, model.RetVal) {
	var challenge model.Challenge
	ret := db.WithTransactionDB(tx, func(tx2 *gorm.DB) model.RetVal {
		var createRet model.RetVal
		challenge, createRet = CreateChallenge(tx2, form)
		return createRet
	})
	return challenge, ret
}

func GetChallengeWithFlags(tx *gorm.DB, challenge model.Challenge) (model.Challenge, model.RetVal) {
	return db.InitChallengeRepo(tx).GetByID(challenge.ID, db.GetOptions{
		Preloads: map[string]db.GetOptions{"ChallengeFlags": {}},
	})
}

func UpdateChallengeWithTransaction(tx *gorm.DB, challenge model.Challenge, form dto.UpdateChallengeForm) model.RetVal {
	loaded, ret := GetChallengeWithFlags(tx, challenge)
	if !ret.OK {
		return ret
	}
	return db.WithTransactionDB(tx, func(tx2 *gorm.DB) model.RetVal {
		return UpdateChallenge(tx2, loaded, form)
	})
}

func DeleteChallenge(tx *gorm.DB, challenge model.Challenge) model.RetVal {
	return db.InitChallengeRepo(tx).Delete(challenge.RandID)
}

func DeleteChallengeWithTransaction(tx *gorm.DB, challenge model.Challenge) model.RetVal {
	return db.WithTransactionDB(tx, func(tx2 *gorm.DB) model.RetVal {
		return DeleteChallenge(tx2, challenge)
	})
}
