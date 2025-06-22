package repo

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"CBCTF/internel/utils"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type ContestRepo struct {
	BasicRepo[model.Contest]
}

type CreateContestOptions struct {
	Name      string
	Desc      string
	Captcha   string
	Avatar    model.AvatarURL
	Prefix    string
	Size      int
	Start     time.Time
	Duration  time.Duration
	Blood     bool
	Hidden    bool
	Rules     model.StringList
	Prizes    model.Prizes
	Timelines model.Timelines
}

func (c CreateContestOptions) Convert2Model() model.Model {
	return model.Contest{
		Name:      c.Name,
		Desc:      c.Desc,
		Captcha:   c.Captcha,
		Avatar:    c.Avatar,
		Prefix:    c.Prefix,
		Size:      c.Size,
		Start:     c.Start,
		Duration:  c.Duration,
		Blood:     c.Blood,
		Hidden:    c.Hidden,
		Rules:     c.Rules,
		Prizes:    c.Prizes,
		Timelines: c.Timelines,
	}
}

type UpdateContestOptions struct {
	Name      *string
	Desc      *string
	Captcha   *string
	Avatar    *model.AvatarURL
	Prefix    *string
	Size      *int
	Start     *time.Time
	Duration  *time.Duration
	Blood     *bool
	Hidden    *bool
	Rules     *model.StringList
	Prizes    *model.Prizes
	Timelines *model.Timelines
}

func (u UpdateContestOptions) Convert2Map() map[string]any {
	options := make(map[string]any)
	if u.Name != nil {
		options["name"] = *u.Name
	}
	if u.Desc != nil {
		options["desc"] = *u.Desc
	}
	if u.Captcha != nil {
		options["captcha"] = *u.Captcha
	}
	if u.Avatar != nil {
		options["avatar"] = *u.Avatar
	}
	if u.Prefix != nil {
		options["prefix"] = *u.Prefix
	}
	if u.Size != nil {
		options["size"] = *u.Size
	}
	if u.Start != nil {
		options["start"] = *u.Start
	}
	if u.Duration != nil {
		options["duration"] = *u.Duration
	}
	if u.Blood != nil {
		options["blood"] = *u.Blood
	}
	if u.Hidden != nil {
		options["hidden"] = *u.Hidden
	}
	if u.Rules != nil {
		options["rules"] = *u.Rules
	}
	if u.Prizes != nil {
		options["prizes"] = *u.Prizes
	}
	if u.Timelines != nil {
		options["timelines"] = *u.Timelines
	}
	return options
}

func InitContestRepo(tx *gorm.DB) *ContestRepo {
	return &ContestRepo{
		BasicRepo: BasicRepo[model.Contest]{
			DB: tx,
		},
	}
}

func (c *ContestRepo) IsUniqueName(name string) bool {
	_, ok, _ := c.GetByUniqueKey("name", name)
	return !ok
}

func (c *ContestRepo) Delete(idL ...uint) (bool, string) {
	contestL, _, ok, msg := c.List(-1, -1, GetOptions{
		Conditions: map[string]any{"id": idL},
		Selects:    []string{"id", "name"},
		Preloads: map[string]GetOptions{
			"Teams":             {Selects: []string{"id"}},
			"Notices":           {Selects: []string{"id"}},
			"ContestChallenges": {Selects: []string{"id"}},
			"ContestFlags":      {Selects: []string{"id"}},
			"Submissions":       {Selects: []string{"id"}},
		},
	})
	if !ok && msg != i18n.ContestNotFound {
		return ok, msg
	}
	teamIDL, noticeIDL, contestChallengeIDL, contestFlagIDL, submissionIDL := make([]uint, 0), make([]uint, 0), make([]uint, 0), make([]uint, 0), make([]uint, 0)
	for _, contest := range contestL {
		deletedName := fmt.Sprintf("%s_deleted_%s", contest.Name, utils.RandStr(6))
		if ok, msg = c.Update(contest.ID, UpdateContestOptions{
			Name: &deletedName,
		}); !ok {
			return false, msg
		}
		for _, team := range contest.Teams {
			teamIDL = append(teamIDL, team.ID)
		}
		for _, notice := range contest.Notices {
			noticeIDL = append(noticeIDL, notice.ID)
		}
		for _, contestChallenge := range contest.ContestChallenges {
			contestChallengeIDL = append(contestChallengeIDL, contestChallenge.ID)
		}
		for _, contestFlag := range contest.ContestFlags {
			contestFlagIDL = append(contestFlagIDL, contestFlag.ID)
		}
		for _, submission := range contest.Submissions {
			submissionIDL = append(submissionIDL, submission.ID)
		}
	}
	if ok, msg = InitTeamRepo(c.DB).Delete(teamIDL...); !ok {
		return false, msg
	}
	if ok, msg = InitNoticeRepo(c.DB).Delete(noticeIDL...); !ok {
		return false, msg
	}
	if ok, msg = InitContestChallengeRepo(c.DB).Delete(contestChallengeIDL...); !ok {
		return false, msg
	}
	if ok, msg = InitContestFlagRepo(c.DB).Delete(contestFlagIDL...); !ok {
		return false, msg
	}
	if ok, msg = InitSubmissionRepo(c.DB).Delete(submissionIDL...); !ok {
		return false, msg
	}
	if res := c.DB.Model(&model.Contest{}).Where("id IN ?", idL).Delete(&model.Contest{}); res.Error != nil {
		log.Logger.Warningf("Failed to delete Contest: %v", res.Error)
		return false, i18n.DeleteContestError
	}
	return true, i18n.Success
}
