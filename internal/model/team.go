package model

import (
	"CBCTF/internal/config"
	"CBCTF/internal/form"
	"CBCTF/internal/utils"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/plugin/optimisticlock"
	"math"
	"strings"
	"time"
)

type Team struct {
	ID        uint                   `gorm:"primarykey" json:"id"`
	Name      string                 `gorm:"index:idx_name_contest_id_deleted,unique;not null" json:"name"`
	ContestID uint                   `gorm:"index:idx_name_contest_id_deleted,unique;not null" json:"contest_id"`
	Desc      string                 `json:"desc"`
	Captcha   string                 `json:"-"`
	Avatar    string                 `json:"avatar"`
	Score     float64                `json:"score" gorm:"default:0"`
	Last      time.Time              `json:"last"`
	Banned    bool                   `gorm:"default:false" json:"banned"`
	Hidden    bool                   `gorm:"default:false" json:"hidden"`
	CaptainID uint                   `json:"captain_id"`
	Rank      int                    `json:"rank" gorm:"default:-1"`
	Users     []*User                `gorm:"many2many:user_teams;" json:"-"`
	CreatedAt time.Time              `json:"-"`
	UpdatedAt time.Time              `json:"-"`
	DeletedAt gorm.DeletedAt         `gorm:"index;index:idx_name_contest_id_deleted,unique" json:"-"`
	Version   optimisticlock.Version `json:"-" gorm:"default:1"`
}

// MarshalJSON 重写 MarshalJSON 方法, 使 Avatar 返回完整的 URL, 转换 Users 为数量
func (t *Team) MarshalJSON() ([]byte, error) {
	type Tmp Team // 定义一个别名以避免递归调用
	return json.Marshal(&struct {
		*Tmp
		Users  int     `json:"users"`
		Avatar string  `json:"avatar"`
		Score  float64 `json:"score"`
	}{
		Tmp:    (*Tmp)(t),
		Users:  len(t.Users),
		Avatar: fmt.Sprintf("%s/%s", config.Env.Backend, strings.TrimPrefix(t.Avatar, "/")),
		Score:  math.Trunc(t.Score*100) / 100,
	})
}

func InitTeam(form form.CreateTeamForm, captain User, contestID uint) Team {
	captcha := utils.UUID()
	if form.Captcha != "" {
		captcha = form.Captcha
	}
	return Team{
		Name:      form.Name,
		Desc:      form.Desc,
		Captcha:   captcha,
		Avatar:    "",
		Banned:    false,
		Hidden:    false,
		CaptainID: captain.ID,
		ContestID: contestID,
		Last:      time.Now(),
		Users:     []*User{&captain},
	}
}
