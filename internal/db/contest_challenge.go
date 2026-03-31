package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"sort"

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

func (c *ContestChallengeRepo) ListContestImages(contestID uint) ([]string, bool, model.RetVal) {
	imageSet := make(map[string]struct{})
	images := make([]string, 0)
	addImage := func(image string) {
		if image == "" {
			return
		}
		if _, ok := imageSet[image]; ok {
			return
		}
		imageSet[image] = struct{}{}
		images = append(images, image)
	}

	dynamicChallenges, _, ret := c.List(-1, -1, GetOptions{
		Conditions: map[string]any{"contest_id": contestID, "type": model.DynamicChallengeType},
		Preloads:   map[string]GetOptions{"Challenge": {}},
	})
	if !ret.OK && ret.Msg != i18n.Model.NotFound {
		return nil, false, ret
	}
	for _, contestChallenge := range dynamicChallenges {
		addImage(contestChallenge.Challenge.GeneratorImage)
	}

	podChallenges, _, ret := c.List(-1, -1, GetOptions{
		Conditions: map[string]any{"contest_id": contestID, "type": model.PodsChallengeType},
		Preloads:   map[string]GetOptions{"Challenge": {}},
	})
	if !ret.OK && ret.Msg != i18n.Model.NotFound {
		return nil, false, ret
	}
	for _, contestChallenge := range podChallenges {
		for _, pod := range contestChallenge.Challenge.Template.Pods {
			for _, container := range pod.Containers {
				addImage(container.Image)
			}
		}
	}

	sort.Strings(images)
	return images, len(podChallenges) > 0, model.SuccessRetVal()
}

func (c *ContestChallengeRepo) ListUnsolvedID(teamID, contestID uint, category string, limit, offset int) ([]uint, int64, model.RetVal) {
	base := c.DB.Model(&model.ContestChallenge{}).
		Where("contest_challenges.contest_id = ?", contestID).
		Where("contest_challenges.hidden = ?", false).
		Where("EXISTS (SELECT 1 FROM contest_flags WHERE contest_flags.contest_challenge_id = contest_challenges.id AND contest_flags.deleted_at IS NULL)").
		Where(`
			EXISTS (
				SELECT 1
				FROM contest_flags
				WHERE contest_flags.contest_challenge_id = contest_challenges.id
					AND contest_flags.deleted_at IS NULL
					AND NOT EXISTS (
						SELECT 1
						FROM team_flags
						WHERE team_flags.team_id = ?
							AND team_flags.contest_flag_id = contest_flags.id
							AND team_flags.solved = true
							AND team_flags.deleted_at IS NULL
					)
			)
		`, teamID)
	if category != "" {
		base = base.Where("contest_challenges.category = ?", category)
	}

	var count int64
	if res := base.Count(&count); res.Error != nil {
		log.Logger.Warningf("Failed to list ContestChallenge: %s", res.Error)
		return nil, 0, model.RetVal{Msg: i18n.Model.ContestChallenge.GetError, Attr: map[string]any{"Error": res.Error.Error()}}
	}

	var ids []uint
	res := base.Select("contest_challenges.id").
		Order("contest_challenges.id ASC").
		Limit(limit).Offset(offset).
		Scan(&ids)
	if res.Error != nil {
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
