package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type ContestRepo struct {
	BaseRepo[model.Contest]
}

type CreateContestOptions struct {
	Name        string
	Description string
	Captcha     string
	Picture     model.FileURL
	Prefix      string
	Size        int
	Start       time.Time
	Duration    time.Duration
	Blood       bool
	Hidden      bool
	Victims     int64
	Rules       model.StringList
	Prizes      model.Prizes
	Timelines   model.Timelines
}

func (c CreateContestOptions) Convert2Model() model.Model {
	return model.Contest{
		Name:        c.Name,
		Description: c.Description,
		Captcha:     c.Captcha,
		Picture:     c.Picture,
		Prefix:      c.Prefix,
		Size:        c.Size,
		Start:       c.Start,
		Duration:    c.Duration,
		Blood:       c.Blood,
		Hidden:      c.Hidden,
		Victims:     c.Victims,
		Rules:       c.Rules,
		Prizes:      c.Prizes,
		Timelines:   c.Timelines,
	}
}

type UpdateContestOptions struct {
	Name        *string
	Description *string
	Captcha     *string
	Picture     *model.FileURL
	Prefix      *string
	Size        *int
	Start       *time.Time
	Duration    *time.Duration
	Blood       *bool
	Hidden      *bool
	UserCount   *int64
	TeamCount   *int64
	NoticeCount *int64
	Victims     *int64
	Rules       *model.StringList
	Prizes      *model.Prizes
	Timelines   *model.Timelines
}

func (u UpdateContestOptions) Convert2Map() map[string]any {
	options := make(map[string]any)
	if u.Name != nil {
		options["name"] = *u.Name
	}
	if u.Description != nil {
		options["description"] = *u.Description
	}
	if u.Captcha != nil {
		options["captcha"] = *u.Captcha
	}
	if u.Picture != nil {
		options["picture"] = *u.Picture
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
	if u.UserCount != nil {
		options["user_count"] = *u.UserCount
	}
	if u.TeamCount != nil {
		options["team_count"] = *u.TeamCount
	}
	if u.NoticeCount != nil {
		options["notice_count"] = *u.NoticeCount
	}
	if u.Victims != nil {
		options["victims"] = *u.Victims
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
		BaseRepo: BaseRepo[model.Contest]{
			DB: tx,
		},
	}
}

func (c *ContestRepo) GetByUserID(userID uint) ([]model.Contest, model.RetVal) {
	var contests []model.Contest
	res := c.DB.Raw(`
		SlECT contests.* FROM contests
		INNER JOIN user_contests ON user_contests.contest_id = contests.id
		WHERE user_contests.user_id = ? AND contests.deleted_at IS NULL
	`, userID).Scan(&contests)
	if res.Error != nil {
		log.Logger.Fatalf("Failed to get Contests: %s", res.Error)
		return nil, model.RetVal{Msg: i18n.Model.GetError, Attr: map[string]any{"Model": model.Contest{}.ModelName(), "Error": res.Error.Error()}}
	}
	return contests, model.RetVal{}
}

func (c *ContestRepo) GetIDByUserID(userID uint) ([]uint, model.RetVal) {
	contests, ret := c.GetByUserID(userID)
	if !ret.OK {
		return nil, ret
	}
	var contestIDL []uint
	for _, contest := range contests {
		contestIDL = append(contestIDL, contest.ID)
	}
	return contestIDL, ret
}

func (c *ContestRepo) Delete(idL ...uint) model.RetVal {
	contestL, _, ret := c.List(-1, -1, GetOptions{
		Conditions: map[string]any{"id": idL},
		Preloads: map[string]GetOptions{
			"Teams":             {},
			"Notices":           {},
			"ContestChallenges": {},
			"ContestFlags":      {},
			"Submissions":       {},
		},
	})
	if !ret.OK {
		if ret.Msg != i18n.Model.NotFound {
			return ret
		}
		return model.SuccessRetVal()
	}
	teamIDL, noticeIDL, contestChallengeIDL, contestFlagIDL, submissionIDL := make([]uint, 0), make([]uint, 0), make([]uint, 0), make([]uint, 0), make([]uint, 0)
	for _, contest := range contestL {
		deletedName := fmt.Sprintf("%s_deleted_%s", contest.Name, utils.RandStr(6))
		if ret = c.Update(contest.ID, UpdateContestOptions{
			Name: &deletedName,
		}); !ret.OK {
			return ret
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
	if ret = InitTeamRepo(c.DB).Delete(teamIDL...); !ret.OK {
		return ret
	}
	if ret = InitNoticeRepo(c.DB).Delete(noticeIDL...); !ret.OK {
		return ret
	}
	if ret = InitContestChallengeRepo(c.DB).Delete(contestChallengeIDL...); !ret.OK {
		return ret
	}
	if ret = InitContestFlagRepo(c.DB).Delete(contestFlagIDL...); !ret.OK {
		return ret
	}
	if ret = InitSubmissionRepo(c.DB).Delete(submissionIDL...); !ret.OK {
		return ret
	}
	if res := c.DB.Model(&model.Contest{}).Where("id IN ?", idL).Delete(&model.Contest{}); res.Error != nil {
		log.Logger.Warningf("Failed to delete Contest: %s", res.Error)
		return model.RetVal{Msg: i18n.Model.DeleteError, Attr: map[string]any{"Model": model.Contest{}.ModelName(), "Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}
