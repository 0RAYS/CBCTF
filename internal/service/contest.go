package service

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
)

func CreateContest(tx *gorm.DB, form dto.CreateContestForm) (model.Contest, model.RetVal) {
	var result model.Contest
	ret := db.WithTransactionDB(tx, func(tx *gorm.DB) model.RetVal {
		repo := db.InitContestRepo(tx)
		if form.Start.IsZero() {
			form.Start = time.Now()
		}
		if form.Duration == 0 {
			form.Duration = 3600 * 24 * 7
		}
		if form.Victims < 1 {
			form.Victims = 1
		}
		if len(form.Rules) == 0 {
			form.Rules = []string{
				"参赛者必须遵守比赛规则和道德准则",
				"禁止攻击比赛平台和其他参赛者",
				"禁止分享题目答案和解题思路",
				"每支队伍人数不得超过4人",
				"比赛采用动态积分机制",
				"设有First Blood奖励",
				"违规行为将导致成绩作废",
			}
		}
		if len(form.Timelines) == 0 {
			duration := time.Duration(form.Duration) * time.Second
			form.Timelines = []model.Timeline{
				{
					Date:        form.Start,
					Title:       "比赛开始",
					Description: "题目公布, 正式开始解题",
				},
				{
					Date:        form.Start.Add(duration),
					Title:       "比赛结束",
					Description: "停止计分, 公布最终排名",
				},
				{
					Date:        form.Start.Add(duration).Add(time.Hour * 2),
					Title:       "截至提交 Writeup",
					Description: "截至提交 Writeup, 成绩作废",
				},
			}
		}
		if len(form.Prizes) == 0 {
			form.Prizes = []model.Prize{
				{
					Amount:      "$0",
					Description: "",
				},
			}
		}
		contest, ret := repo.Create(model.Contest{
			Name:        form.Name,
			Description: form.Description,
			Captcha:     form.Captcha,
			Picture:     "",
			Prefix:      form.Prefix,
			Size:        form.Size,
			Start:       form.Start,
			Duration:    time.Duration(form.Duration) * time.Second,
			Blood:       form.Blood,
			Hidden:      form.Hidden,
			Victims:     form.Victims,
			Rules:       form.Rules,
			Prizes:      form.Prizes,
			Timelines:   form.Timelines,
		})
		if !ret.OK {
			result = contest
			return ret
		}
		// GORM applies default:true to zero-value bools on insert. Persist the
		// explicit form choices as a map update, just as the edit endpoint does.
		if contest.Blood != form.Blood || contest.Hidden != form.Hidden {
			if ret = repo.Update(contest.ID, db.UpdateContestOptions{Blood: &form.Blood, Hidden: &form.Hidden}); !ret.OK {
				return ret
			}
			result, ret = repo.GetByID(contest.ID)
			return ret
		}
		result = contest
		return ret
	})
	return result, ret
}

func UpdateContest(tx *gorm.DB, contest model.Contest, form dto.UpdateContestForm) model.RetVal {
	repo := db.InitContestRepo(tx)
	if form.Victims != nil && *form.Victims < 1 {
		form.Victims = new(int64(1))
	}
	var duration *time.Duration
	if form.Duration != nil {
		duration = new(time.Duration(*form.Duration) * time.Second)
	}
	return repo.Update(contest.ID, db.UpdateContestOptions{
		Name:        form.Name,
		Description: form.Description,
		Captcha:     form.Captcha,
		Prefix:      form.Prefix,
		Size:        form.Size,
		Start:       form.Start,
		Duration:    duration,
		Blood:       form.Blood,
		Hidden:      form.Hidden,
		Victims:     form.Victims,
		Rules:       form.Rules,
		Prizes:      form.Prizes,
		Timelines:   form.Timelines,
	})
}

// UpdateContestLeaderboard 开启或关闭排行榜黑灯模式。
func UpdateContestLeaderboard(tx *gorm.DB, contest model.Contest, hidden bool) (model.Contest, model.RetVal) {
	repo := db.InitContestRepo(tx)
	if ret := repo.Update(contest.ID, db.UpdateContestOptions{LeaderboardHidden: &hidden}); !ret.OK {
		return contest, ret
	}
	updated, ret := repo.GetByID(contest.ID)
	if !ret.OK {
		return contest, ret
	}
	return updated, model.SuccessRetVal()
}

// UpdateContestPause 暂停或恢复比赛。暂停期间选手无法访问题目，暂停的时间在恢复时补回比赛时长。
func UpdateContestPause(tx *gorm.DB, contest model.Contest, paused bool) (model.Contest, model.RetVal) {
	repo := db.InitContestRepo(tx)
	now := time.Now()
	if paused {
		if contest.Paused {
			return contest, model.SuccessRetVal()
		}
		if contest.Status() != model.ContestIsRunning {
			return contest, model.FailRetVal(i18n.Model.Contest.NotRunning, "Contest", fmt.Sprintf("contest_id=%d", contest.ID), nil)
		}
		if ret := repo.Update(contest.ID, db.UpdateContestOptions{Paused: &paused, PausedAt: &now}); !ret.OK {
			return contest, ret
		}
	} else {
		if !contest.Paused {
			return contest, model.SuccessRetVal()
		}
		duration := contest.Duration + now.Sub(contest.PausedAt)
		if ret := repo.Update(contest.ID, db.UpdateContestOptions{Paused: &paused, Duration: &duration}); !ret.OK {
			return contest, ret
		}
	}
	updated, ret := repo.GetByID(contest.ID)
	if !ret.OK {
		return contest, ret
	}
	return updated, model.SuccessRetVal()
}
