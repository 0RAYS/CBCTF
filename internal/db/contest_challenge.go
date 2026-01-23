package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"

	"gorm.io/gorm"
)

type ContestChallengeRepo struct {
	BaseRepo[model.ContestChallenge]
}

type CreateContestChallengeOptions struct {
	ContestID   uint
	ChallengeID uint
	Name        string
	Desc        string
	Type        string
	Category    string
	Hidden      bool
	Attempt     int64
	Hints       model.StringList
	Tags        model.StringList
}

func (c CreateContestChallengeOptions) Convert2Model() model.Model {
	return model.ContestChallenge{
		ContestID:   c.ContestID,
		ChallengeID: c.ChallengeID,
		Name:        c.Name,
		Desc:        c.Desc,
		Type:        c.Type,
		Category:    c.Category,
		Hidden:      c.Hidden,
		Attempt:     c.Attempt,
		Hints:       c.Hints,
		Tags:        c.Tags,
	}
}

type UpdateContestChallengeOptions struct {
	Name        *string
	Desc        *string
	Hidden      *bool
	Attempt     *int64
	Hints       *model.StringList
	Tags        *model.StringList
	DeletedSalt *string
}

func (u UpdateContestChallengeOptions) Convert2Map() map[string]any {
	options := make(map[string]any)
	if u.Name != nil {
		options["name"] = *u.Name
	}
	if u.Desc != nil {
		options["desc"] = *u.Desc
	}
	if u.Hidden != nil {
		options["hidden"] = *u.Hidden
	}
	if u.Attempt != nil {
		options["attempt"] = *u.Attempt
	}
	if u.Hints != nil {
		options["hints"] = *u.Hints
	}
	if u.Tags != nil {
		options["tags"] = *u.Tags
	}
	if u.DeletedSalt != nil {
		options["deleted_salt"] = *u.DeletedSalt
	}
	return options
}

func InitContestChallengeRepo(tx *gorm.DB) *ContestChallengeRepo {
	return &ContestChallengeRepo{
		BaseRepo: BaseRepo[model.ContestChallenge]{
			DB: tx,
		},
	}
}

func (c *ContestChallengeRepo) IsUniqueContestChallenge(contestID uint, challengeID uint) bool {
	_, ret := c.Get(GetOptions{
		Conditions: map[string]any{"contest_id": contestID, "challenge_id": challengeID},
		Selects:    []string{"id"},
	})
	return !ret.OK
}

func (c *ContestChallengeRepo) ListCategories(contestID uint, t string) ([]string, model.RetVal) {
	var categories []string
	tx := c.DB.Model(&model.ContestChallenge{}).Where("contest_id = ?", contestID)
	if t != "" {
		tx = tx.Where("type = ?", t)
	}
	if res := tx.Select("distinct category").Find(&categories); res.Error != nil {
		log.Logger.Warningf("Failed to list ContestChallenge categories: %s", res.Error)
		return nil, model.RetVal{Msg: i18n.Model.GetError, Attr: map[string]interface{}{"Error": res.Error.Error()}}
	}
	return categories, model.SuccessRetVal()
}

func (c *ContestChallengeRepo) Delete(idL ...uint) model.RetVal {
	contestChallengeL, _, ret := c.List(-1, -1, GetOptions{
		Conditions: map[string]any{"id": idL},
		Selects:    []string{"id"},
		Preloads: map[string]GetOptions{
			"ContestFlags": {Selects: []string{"id", "contest_challenge_id"}},
			"Submissions":  {Selects: []string{"id", "contest_challenge_id"}},
		},
	})
	if !ret.OK {
		if ret.Msg != i18n.ContestChallengeNotFound {
			return ret
		}
		return model.SuccessRetVal()
	}
	contestFlagIDL, submissionIDL := make([]uint, 0), make([]uint, 0)
	for _, contestChallenge := range contestChallengeL {
		deletedSalt := utils.UUID()
		if ret = c.Update(contestChallenge.ID, UpdateContestChallengeOptions{DeletedSalt: &deletedSalt}); !ret.OK {
			return ret
		}
		for _, contestFlag := range contestChallenge.ContestFlags {
			contestFlagIDL = append(contestFlagIDL, contestFlag.ID)
		}
		for _, submission := range contestChallenge.Submissions {
			submissionIDL = append(submissionIDL, submission.ID)
		}
	}
	if ret = InitContestFlagRepo(c.DB).Delete(contestFlagIDL...); !ret.OK {
		return ret
	}
	if ret = InitSubmissionRepo(c.DB).Delete(submissionIDL...); !ret.OK {
		return ret
	}
	if res := c.DB.Model(&model.ContestChallenge{}).Where("id IN ?", idL).Delete(&model.ContestChallenge{}); res.Error != nil {
		log.Logger.Warningf("Failed to delete ContestChallenge: %s", res.Error)
		return model.RetVal{Msg: i18n.Model.DeleteError, Attr: map[string]interface{}{"Model": model.ContestChallenge{}.GetModelName(), "Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}
