package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"

	"gorm.io/gorm"
)

type ContestChallengeRepo struct {
	BaseRepo[model.ContestChallenge]
}

type CreateContestChallengeOptions struct {
	ContestID   uint
	ChallengeID uint
	Name        string
	Description string
	Type        model.ChallengeType
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
		Description: c.Description,
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
	Description *string
	Hidden      *bool
	Attempt     *int64
	Hints       *model.StringList
	Tags        *model.StringList
}

func (u UpdateContestChallengeOptions) Convert2Map() map[string]any {
	options := make(map[string]any)
	if u.Name != nil {
		options["name"] = *u.Name
	}
	if u.Description != nil {
		options["description"] = *u.Description
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
	})
	return !ret.OK
}

func (c *ContestChallengeRepo) ListCategories(contestID uint, t model.ChallengeType) ([]string, model.RetVal) {
	var categories []string
	tx := c.DB.Model(&model.ContestChallenge{}).Where("contest_id = ?", contestID)
	if t != "" {
		tx = tx.Where("type = ?", t)
	}
	if res := tx.Distinct().Order("category ASC").Pluck("category", &categories); res.Error != nil {
		log.Logger.Warningf("Failed to list ContestChallenge categories: %s", res.Error)
		return nil, model.RetVal{Msg: i18n.Model.ContestChallenge.GetError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return categories, model.SuccessRetVal()
}

func (c *ContestChallengeRepo) ListUnsolvedID(teamID, contestID uint, category string, limit, offset int) ([]uint, int64, model.RetVal) {
	subSolved := c.DB.Table("submissions").
		Select("COUNT(*)").
		Where("submissions.contest_challenge_id = contest_challenges.id").
		Where("submissions.team_id = ?", teamID).
		Where("submissions.solved = ?", true).
		Where("submissions.deleted_at IS NULL")

	subFlags := c.DB.Table("contest_flags").
		Select("COUNT(*)").
		Where("contest_flags.contest_challenge_id = contest_challenges.id").
		Where("contest_flags.deleted_at IS NULL")

	var count int64
	res := c.DB.Table("contest_challenges").
		Select("COUNT(*)").
		Where("contest_challenges.contest_id = ?", contestID).
		Where("contest_challenges.hidden = ?", false).
		Where("contest_challenges.deleted_at IS NULL").
		Where("(?) < (?)", subSolved, subFlags)
	if category != "" {
		res = res.Where("contest_challenges.category = ?", category)
	}
	if res = res.Scan(&count); res.Error != nil {
		log.Logger.Warningf("Failed to list ContestChallenge: %s", res.Error)
		return nil, 0, model.RetVal{Msg: i18n.Model.ContestChallenge.GetError, Attr: map[string]any{"Error": res.Error.Error()}}
	}

	var ids []uint
	res = c.DB.Table("contest_challenges").
		Select("contest_challenges.id").
		Where("contest_challenges.contest_id = ?", contestID).
		Where("contest_challenges.hidden = ?", false).
		Where("contest_challenges.deleted_at IS NULL").
		Where("(?) < (?)", subSolved, subFlags)
	if category != "" {
		res = res.Where("contest_challenges.category = ?", category)
	}
	if res = res.Limit(limit).Offset(offset).Scan(&ids); res.Error != nil {
		log.Logger.Warningf("Failed to list ContestChallenge: %s", res.Error)
		return nil, 0, model.RetVal{Msg: i18n.Model.ContestChallenge.GetError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return ids, count, model.SuccessRetVal()
}

func (c *ContestChallengeRepo) Delete(idL ...uint) model.RetVal {
	contestChallengeL, _, ret := c.List(-1, -1, GetOptions{
		Conditions: map[string]any{"id": idL},
		Preloads:   map[string]GetOptions{"ContestFlags": {}, "Submissions": {}},
	})
	if !ret.OK {
		if ret.Msg != i18n.Model.NotFound {
			return ret
		}
		return model.SuccessRetVal()
	}
	contestFlagIDL, submissionIDL := make([]uint, 0), make([]uint, 0)
	for _, contestChallenge := range contestChallengeL {
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
		return model.RetVal{Msg: i18n.Model.ContestChallenge.DeleteError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}
