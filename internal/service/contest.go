package service

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"time"

	"gorm.io/gorm"
)

func CreateContest(tx *gorm.DB, form dto.CreateContestForm) (model.Contest, model.RetVal) {
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
		form.Rules = model.StringList{
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
		form.Timelines = model.Timelines{
			{
				Date:        form.Start,
				Title:       "比赛开始",
				Description: "题目公布, 正式开始解题",
			},
			{
				Date:        form.Start.Add(time.Duration(form.Duration)),
				Title:       "比赛结束",
				Description: "停止计分, 公布最终排名",
			},
			{
				Date:        form.Start.Add(time.Duration(form.Duration)).Add(time.Hour * 2),
				Title:       "截至提交 WriteUp",
				Description: "截至提交 WriteUp, 成绩作废",
			},
		}
	}
	if len(form.Prizes) == 0 {
		form.Prizes = model.Prizes{
			{
				Amount:      "$0",
				Description: "",
			},
		}
	}
	return repo.Create(db.CreateContestOptions{
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
}

func UpdateContest(tx *gorm.DB, contest model.Contest, form dto.UpdateContestForm) model.RetVal {
	repo := db.InitContestRepo(tx)
	if form.Duration != nil {
		*form.Duration = *form.Duration * 1e9
	}
	if form.Victims != nil && *form.Victims < 1 {
		form.Victims = utils.Ptr(int64(1))
	}
	return repo.Update(contest.ID, db.UpdateContestOptions{
		Name:        form.Name,
		Description: form.Description,
		Captcha:     form.Captcha,
		Prefix:      form.Prefix,
		Size:        form.Size,
		Start:       form.Start,
		Duration:    (*time.Duration)(form.Duration),
		Blood:       form.Blood,
		Hidden:      form.Hidden,
		Victims:     form.Victims,
		Rules:       form.Rules,
		Prizes:      form.Prizes,
		Timelines:   form.Timelines,
	})
}
