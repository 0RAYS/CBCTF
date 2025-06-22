package repo

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"CBCTF/internel/utils"
	"gorm.io/gorm"
)

type ContestChallengeRepo struct {
	BasicRepo[model.ContestChallenge]
}

type CreateContestChallengeOptions struct {
	ContestID   uint
	ChallengeID uint
	Type        string
	Name        string
	Desc        string
	Hidden      bool
	Attempt     int64
	Hints       model.StringList
	Tags        model.StringList
}

func (c CreateContestChallengeOptions) Convert2Model() model.Model {
	return model.ContestChallenge{
		ContestID:   c.ContestID,
		ChallengeID: c.ChallengeID,
		Type:        c.Type,
		Name:        c.Name,
		Desc:        c.Desc,
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
		BasicRepo: BasicRepo[model.ContestChallenge]{
			DB: tx,
		},
	}
}

func (c *ContestChallengeRepo) IsUniqueContestChallenge(contestID uint, challengeID uint) bool {
	_, ok, _ := c.Get(GetOptions{
		Conditions: map[string]any{
			"contest_id":   contestID,
			"challenge_id": challengeID,
		},
		Selects: []string{"id"},
	})
	return !ok
}

func (c *ContestChallengeRepo) Delete(idL ...uint) (bool, string) {
	contestChallengeL, _, ok, msg := c.List(-1, -1, GetOptions{
		Conditions: map[string]any{"id": idL},
		Selects:    []string{"id"},
		Preloads: map[string]GetOptions{
			"ContestFlags": {Selects: []string{"id"}},
			"Submissions":  {Selects: []string{"id"}},
		},
	})
	if !ok && msg != i18n.ContestChallengeNotFound {
		return false, msg
	}
	contestFlagIDL, submissionIDL := make([]uint, 0), make([]uint, 0)
	for _, contestChallenge := range contestChallengeL {
		deletedSalt := utils.UUID()
		if ok, msg = c.Update(contestChallenge.ID, UpdateContestChallengeOptions{DeletedSalt: &deletedSalt}); !ok {
			return false, msg
		}
		for _, contestFlag := range contestChallenge.ContestFlags {
			contestFlagIDL = append(contestFlagIDL, contestFlag.ID)
		}
		for _, submission := range contestChallenge.Submissions {
			submissionIDL = append(submissionIDL, submission.ID)
		}
	}
	if ok, msg = InitContestFlagRepo(c.DB).Delete(contestFlagIDL...); !ok {
		return false, msg
	}
	if ok, msg = InitSubmissionRepo(c.DB).Delete(submissionIDL...); !ok {
		return false, msg
	}
	if res := c.DB.Model(&model.ContestChallenge{}).Where("id IN ?", idL).Delete(&model.ContestChallenge{}); res.Error != nil {
		log.Logger.Warningf("Failed to delete ContestChallenge: %v", res.Error)
		return false, i18n.DeleteContestChallengeError
	}
	return true, i18n.Success
}
