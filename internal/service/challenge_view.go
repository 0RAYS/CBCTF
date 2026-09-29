package service

import (
	"gorm.io/gorm"

	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/view"
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

	file, ret := db.InitFileRepo(tx).Get(db.GetOptions{
		Conditions: map[string]any{
			"model":    model.Name(challenge),
			"model_id": challenge.ID,
			"type":     model.ChallengeFileType,
		},
	})
	result.FileName = file.Filename
	if !ret.OK && ret.Msg != i18n.Model.NotFound {
		result.Unavailable = []string{"file"}
	}
	return result
}

func ListChallengeViews(tx *gorm.DB, form dto.GetChallengesForm) ([]view.ChallengeView, int64, model.RetVal) {
	options := db.GetOptions{
		Conditions: make(map[string]any),
		Search:     make(map[string]string),
		Preloads:   map[string]db.GetOptions{"ChallengeFlags": {}},
	}
	if form.Type != "" {
		options.Conditions["type"] = form.Type
	}
	if form.Category != "" {
		options.Conditions["category"] = form.Category
	}
	if form.Name != "" {
		options.Search["name"] = form.Name
	}
	if form.Description != "" {
		options.Search["description"] = form.Description
	}
	challenges, count, ret := db.InitChallengeRepo(tx).List(form.Limit, form.Offset, options)
	if !ret.OK {
		return nil, 0, ret
	}
	views := make([]view.ChallengeView, 0, len(challenges))
	for _, challenge := range challenges {
		views = append(views, BuildChallengeView(tx, challenge))
	}
	return views, count, model.SuccessRetVal()
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

func GetChallengeWithFlags(tx *gorm.DB, challenge model.Challenge) (model.Challenge, model.RetVal) {
	return db.InitChallengeRepo(tx).GetByID(challenge.ID, db.GetOptions{
		Preloads: map[string]db.GetOptions{"ChallengeFlags": {}},
	})
}

func DeleteChallenge(tx *gorm.DB, challenge model.Challenge) model.RetVal {
	return db.WithTransactionDB(tx, func(tx *gorm.DB) model.RetVal {
		if ret := stopGeneratorResources(tx, db.GetOptions{Conditions: map[string]any{"challenge_id": challenge.ID}}); !ret.OK {
			return ret
		}
		return db.InitChallengeRepo(tx).Delete(challenge.RandID)
	})
}
