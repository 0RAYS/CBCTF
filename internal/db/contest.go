package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
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

func (c *ContestRepo) GetIDByUserID(userID uint) ([]uint, model.RetVal) {
	var contestIDL []uint
	res := c.DB.Model(&model.UserContest{}).
		Where("user_id = ?", userID).
		Pluck("contest_id", &contestIDL)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Contest IDs by user: %s", res.Error)
		return nil, model.RetVal{Msg: i18n.Model.Contest.GetError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return contestIDL, model.SuccessRetVal()
}

func (c *ContestRepo) CountUsers(contestID uint) (int64, model.RetVal) {
	var count int64
	res := c.DB.Model(&model.UserContest{}).Where("contest_id = ?", contestID).Count(&count)
	if res.Error != nil {
		log.Logger.Warningf("Failed to count contest users: %s", res.Error)
		return 0, model.RetVal{Msg: i18n.Model.Contest.GetError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return count, model.SuccessRetVal()
}

func (c *ContestRepo) CountTeams(contestID uint) (int64, model.RetVal) {
	var count int64
	res := c.DB.Model(&model.Team{}).Where("contest_id = ?", contestID).Count(&count)
	if res.Error != nil {
		log.Logger.Warningf("Failed to count contest teams: %s", res.Error)
		return 0, model.RetVal{Msg: i18n.Model.Contest.GetError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return count, model.SuccessRetVal()
}

func (c *ContestRepo) CountUsersMap(contestIDL ...uint) (map[uint]int64, model.RetVal) {
	result := make(map[uint]int64)
	if len(contestIDL) == 0 {
		return result, model.SuccessRetVal()
	}

	type row struct {
		ContestID uint
		Count     int64
	}

	rows := make([]row, 0)
	res := c.DB.Model(&model.UserContest{}).
		Select("contest_id, COUNT(*) AS count").
		Where("contest_id IN ?", contestIDL).
		Group("contest_id").
		Scan(&rows)
	if res.Error != nil {
		log.Logger.Warningf("Failed to count contest users: %s", res.Error)
		return nil, model.RetVal{Msg: i18n.Model.Contest.GetError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	for _, item := range rows {
		result[item.ContestID] = item.Count
	}
	return result, model.SuccessRetVal()
}

func (c *ContestRepo) CountTeamsMap(contestIDL ...uint) (map[uint]int64, model.RetVal) {
	result := make(map[uint]int64)
	if len(contestIDL) == 0 {
		return result, model.SuccessRetVal()
	}

	type row struct {
		ContestID uint
		Count     int64
	}

	rows := make([]row, 0)
	res := c.DB.Model(&model.Team{}).
		Select("contest_id, COUNT(*) AS count").
		Where("contest_id IN ?", contestIDL).
		Group("contest_id").
		Scan(&rows)
	if res.Error != nil {
		log.Logger.Warningf("Failed to count contest teams: %s", res.Error)
		return nil, model.RetVal{Msg: i18n.Model.Contest.GetError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	for _, item := range rows {
		result[item.ContestID] = item.Count
	}
	return result, model.SuccessRetVal()
}

func (c *ContestRepo) CountNotices(contestID uint) (int64, model.RetVal) {
	var count int64
	res := c.DB.Model(&model.Notice{}).Where("contest_id = ?", contestID).Count(&count)
	if res.Error != nil {
		log.Logger.Warningf("Failed to count contest notices: %s", res.Error)
		return 0, model.RetVal{Msg: i18n.Model.Contest.GetError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return count, model.SuccessRetVal()
}

func (c *ContestRepo) Delete(idL ...uint) model.RetVal {
	var teamIDL []uint
	if res := c.DB.Model(&model.Team{}).Where("contest_id IN ?", idL).Pluck("id", &teamIDL); res.Error != nil {
		log.Logger.Warningf("Failed to get Teams for contests %v: %s", idL, res.Error)
		return model.RetVal{Msg: i18n.Model.Team.GetError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	var contestChallengeIDL []uint
	if res := c.DB.Model(&model.ContestChallenge{}).Where("contest_id IN ?", idL).Pluck("id", &contestChallengeIDL); res.Error != nil {
		log.Logger.Warningf("Failed to get ContestChallenges for contests %v: %s", idL, res.Error)
		return model.RetVal{Msg: i18n.Model.ContestChallenge.GetError, Attr: map[string]any{"Error": res.Error.Error()}}
	}

	var ret model.RetVal
	if ret = InitTeamRepo(c.DB).Delete(teamIDL...); !ret.OK {
		return ret
	}
	if res := c.DB.Where("contest_id IN ?", idL).Delete(&model.Notice{}); res.Error != nil {
		log.Logger.Warningf("Failed to delete Notices for contests %v: %s", idL, res.Error)
		return model.RetVal{Msg: i18n.Model.Notice.DeleteError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	if ret = InitContestChallengeRepo(c.DB).Delete(contestChallengeIDL...); !ret.OK {
		return ret
	}
	var contestFlagIDL []uint
	if res := c.DB.Model(&model.ContestFlag{}).Where("contest_id IN ?", idL).Pluck("id", &contestFlagIDL); res.Error != nil {
		log.Logger.Warningf("Failed to get ContestFlags for contests %v: %s", idL, res.Error)
		return model.RetVal{Msg: i18n.Model.ContestFlag.GetError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	if ret = InitContestFlagRepo(c.DB).Delete(contestFlagIDL...); !ret.OK {
		return ret
	}
	if res := c.DB.Where("contest_id IN ?", idL).Delete(&model.Submission{}); res.Error != nil {
		log.Logger.Warningf("Failed to delete Submissions for contests %v: %s", idL, res.Error)
		return model.RetVal{Msg: i18n.Model.Submission.DeleteError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	if ret = InitGeneratorRepo(c.DB).DeleteByContestID(idL...); !ret.OK {
		return ret
	}
	if ret = InitVictimRepo(c.DB).DeleteByContestID(idL...); !ret.OK {
		return ret
	}
	if ret = InitCheatRepo(c.DB).DeleteByContestID(idL...); !ret.OK {
		return ret
	}
	if res := c.DB.Model(&model.Contest{}).Where("id IN ?", idL).Delete(&model.Contest{}); res.Error != nil {
		log.Logger.Warningf("Failed to delete Contest: %s", res.Error)
		return model.RetVal{Msg: i18n.Model.Contest.DeleteError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}
