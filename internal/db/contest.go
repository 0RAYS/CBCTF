package db

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
)

type ContestRepo struct {
	BaseRepo[model.Contest]
}

type UpdateContestOptions struct {
	Name              *string
	Description       *string
	Captcha           *string
	Picture           *model.FileURL
	Prefix            *string
	Size              *int
	Start             *time.Time
	Duration          *time.Duration
	Blood             *bool
	Hidden            *bool
	Paused            *bool
	PausedAt          *time.Time
	LeaderboardHidden *bool
	Victims           *int64
	Rules             *[]string
	Prizes            *[]model.Prize
	Timelines         *[]model.Timeline
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
	if u.Paused != nil {
		options["paused"] = *u.Paused
	}
	if u.PausedAt != nil {
		options["paused_at"] = *u.PausedAt
	}
	if u.LeaderboardHidden != nil {
		options["leaderboard_hidden"] = *u.LeaderboardHidden
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
		DB: tx,
	}
}

// ListUnfinished returns the contests that have not ended yet, including the
// upcoming ones. Use Go duration arithmetic, matching Contest.IsOver exactly.
func (c *ContestRepo) ListUnfinished(now time.Time) ([]model.Contest, model.RetVal) {
	contests, ret := InitContestRepo(c.DB.Select("id", "start", "duration")).FindAll()
	if !ret.OK {
		return nil, ret
	}
	unfinished := make([]model.Contest, 0, len(contests))
	for _, contest := range contests {
		if !now.After(contest.End()) {
			unfinished = append(unfinished, contest)
		}
	}
	return unfinished, model.SuccessRetVal()
}

func (c *ContestRepo) CountUsers(contestID uint) (int64, model.RetVal) {
	var count int64
	res := c.DB.Model(&model.UserContest{}).Where("contest_id = ?", contestID).Count(&count)
	if res.Error != nil {
		target := fmt.Sprintf("contest_id=%d", contestID)
		log.Logger.Warningf("Failed to count contest users: %s (target=%s)", res.Error, target)
		return 0, model.FailRetVal(i18n.Model.Contest.GetError, "Contest", target, res.Error)
	}
	return count, model.SuccessRetVal()
}

func (c *ContestRepo) CountTeams(contestID uint) (int64, model.RetVal) {
	var count int64
	res := c.DB.Model(&model.Team{}).Where("contest_id = ?", contestID).Count(&count)
	if res.Error != nil {
		target := fmt.Sprintf("contest_id=%d", contestID)
		log.Logger.Warningf("Failed to count contest teams: %s (target=%s)", res.Error, target)
		return 0, model.FailRetVal(i18n.Model.Contest.GetError, "Contest", target, res.Error)
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
		target := fmt.Sprintf("contest_id=%v", contestIDL)
		log.Logger.Warningf("Failed to count contest users: %s (target=%s)", res.Error, target)
		return nil, model.FailRetVal(i18n.Model.Contest.GetError, "Contest", target, res.Error)
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
		target := fmt.Sprintf("contest_id=%v", contestIDL)
		log.Logger.Warningf("Failed to count contest teams: %s (target=%s)", res.Error, target)
		return nil, model.FailRetVal(i18n.Model.Contest.GetError, "Contest", target, res.Error)
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
		target := fmt.Sprintf("contest_id=%d", contestID)
		log.Logger.Warningf("Failed to count contest notices: %s (target=%s)", res.Error, target)
		return 0, model.FailRetVal(i18n.Model.Contest.GetError, "Contest", target, res.Error)
	}
	return count, model.SuccessRetVal()
}

func (c *ContestRepo) Delete(idL ...uint) model.RetVal {
	if len(idL) == 0 {
		return model.SuccessRetVal()
	}
	var ret model.RetVal
	if ret = InitContestChallengeRepo(c.DB).DeleteByContestID(idL...); !ret.OK {
		return ret
	}
	if ret = InitTeamRepo(c.DB).DeleteByContestID(idL...); !ret.OK {
		return ret
	}
	if ret = InitNoticeRepo(c.DB).DeleteByContestID(idL...); !ret.OK {
		return ret
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
		target := fmt.Sprintf("id=%v", idL)
		log.Logger.Warningf("Failed to delete Contest: %s (target=%s)", res.Error, target)
		return model.FailRetVal(i18n.Model.Contest.DeleteError, "Contest", target, res.Error)
	}
	return model.SuccessRetVal()
}
