package model

import (
	"CBCTF/internal/config"
	f "CBCTF/internal/form"
	"CBCTF/internal/utils"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/plugin/optimisticlock"
	"strings"
	"time"
)

type Contest struct {
	ID        uint                   `gorm:"primarykey" json:"id"`
	Name      string                 `gorm:"index:idx_name_deleted,unique;not null" json:"name"`
	Desc      string                 `json:"desc"`
	Captcha   string                 `json:"-"`
	Avatar    string                 `json:"avatar"`
	Prefix    string                 `json:"prefix" gorm:"default:'CBCTF'"`
	Size      int                    `json:"size"`
	Start     time.Time              `json:"start"`
	Duration  time.Duration          `json:"-"`
	Blood     bool                   `json:"blood" gorm:"default:true"`
	Hidden    bool                   `gorm:"default:true" json:"hidden"`
	Rules     utils.Strings          `json:"rules" gorm:"type:json"`
	Prizes    utils.Prizes           `json:"prizes" gorm:"type:json"`
	Timelines utils.Timelines        `json:"timelines" gorm:"type:json"`
	Teams     []*Team                `json:"-"`
	Users     []*User                `gorm:"many2many:user_contests;" json:"-"`
	Notices   []*Notice              `json:"-"`
	CreatedAt time.Time              `json:"-"`
	UpdatedAt time.Time              `json:"-"`
	DeletedAt gorm.DeletedAt         `gorm:"index;index:idx_name_deleted,unique;" json:"-"`
	Version   optimisticlock.Version `json:"-"`
}

func (c *Contest) MarshalJSON() ([]byte, error) {
	type Tmp Contest // 定义一个别名以避免递归调用
	return json.Marshal(&struct {
		*Tmp
		Users    int    `json:"users"`
		Teams    int    `json:"teams"`
		Notices  int    `json:"notices"`
		Avatar   string `json:"avatar"`
		Duration int64  `json:"duration"`
	}{
		Tmp:      (*Tmp)(c),
		Users:    len(c.Users),
		Teams:    len(c.Teams),
		Notices:  len(c.Notices),
		Avatar:   fmt.Sprintf("%s/%s", config.Env.Backend, strings.TrimPrefix(c.Avatar, "/")),
		Duration: int64(c.Duration.Seconds()),
	})
}

func (c *Contest) IsOver() bool {
	return time.Now().After(c.Start.Add(c.Duration))
}

func (c *Contest) IsNotStart() bool {
	return time.Now().Before(c.Start)
}

func (c *Contest) IsRunning() bool {
	return (c.IsOver() || c.IsNotStart()) != true
}

func (c *Contest) Status() string {
	if c.IsOver() {
		return "ContestIsOver"
	}
	if c.IsNotStart() {
		return "ContestNotRunning"
	}
	return "ContestIsRunning"
}

func InitContest(form f.CreateContestForm) Contest {
	if len(form.Rules) == 0 {
		form.Rules = utils.Strings{
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
		form.Timelines = utils.Timelines{
			utils.Timeline{
				Date:  form.Start,
				Title: "比赛开始",
				Desc:  "题目公布，正式开始解题",
			},
			utils.Timeline{
				Date:  form.Start.Add(time.Duration(form.Duration)),
				Title: "比赛结束",
				Desc:  "停止计分，公布最终排名",
			},
		}
	}
	if len(form.Prizes) == 0 {
		form.Prizes = utils.Prizes{
			utils.Prize{
				Amount: "$0",
				Desc:   "",
			},
		}
	}
	return Contest{
		Name:      form.Name,
		Desc:      form.Desc,
		Captcha:   form.Captcha,
		Avatar:    "",
		Blood:     form.Blood,
		Size:      form.Size,
		Start:     form.Start,
		Hidden:    form.Hidden,
		Duration:  time.Duration(form.Duration),
		Rules:     form.Rules,
		Prizes:    form.Prizes,
		Timelines: form.Timelines,
	}
}
