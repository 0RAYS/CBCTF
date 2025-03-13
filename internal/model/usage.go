package model

import (
	"CBCTF/internal/utils"
	"gorm.io/gorm"
	"time"
)

const (
	StaticScore      uint = 0
	LinearScore      uint = 1
	LogarithmicScore uint = 2
)

type Usage struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	ContestID      uint           `json:"contest_id"`
	ChallengeID    string         `json:"challenge_id"`
	Name           string         `json:"name" gorm:"not null"`
	Desc           string         `json:"desc"`
	Flag           string         `json:"flag"`
	Category       string         `json:"category"`
	GeneratorImage string         `json:"generator" gorm:"column:generator"`
	DockerImage    string         `json:"docker" gorm:"column:docker"`
	Port           int32          `json:"port" gorm:"default:8080"`
	Type           string         `json:"type" gorm:"default:'static'"`
	Hidden         bool           `json:"hidden" default:"true"`
	Score          float64        `json:"score" gorm:"default:1000"`
	CurrentScore   float64        `json:"current_score" gorm:"default:1000"`
	ScoreType      uint           `json:"score_type" gorm:"default:0"`
	MinScore       float64        `json:"min_score" gorm:"default:100"`
	Decay          float64        `json:"decay" gorm:"default:100"`
	Attempt        int64          `json:"attempt" gorm:"default:0"`
	Solvers        int64          `json:"solvers" gorm:"default:0"`
	Hints          utils.Strings  `json:"hints"`
	Tags           utils.Strings  `json:"tags"`
	First          uint           `json:"first" gorm:"default:0"`
	Second         uint           `json:"second" gorm:"default:0"`
	Third          uint           `json:"third" gorm:"default:0"`
	Last           time.Time      `json:"last"`
	CreatedAt      time.Time      `json:"-"`
	UpdatedAt      time.Time      `json:"-"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

func (u *Usage) CalcScore(solvers int64) float64 {
	var calc float64 = 0
	switch u.ScoreType {
	case StaticScore:
		calc = u.CurrentScore
	case LinearScore:
		calc = u.CurrentScore - float64(solvers)*u.Decay
	case LogarithmicScore:
		calc = (((u.MinScore - u.CurrentScore) / (u.Decay * u.Decay)) * float64(solvers*solvers)) + u.CurrentScore
	default:
		calc = u.CurrentScore
	}
	if calc < u.MinScore {
		calc = u.MinScore
	}
	return calc
}

func InitUsage(challenge Challenge, contestID uint) Usage {
	return Usage{
		ContestID:      contestID,
		ChallengeID:    challenge.ID,
		Name:           challenge.Name,
		Desc:           challenge.Desc,
		Flag:           challenge.Flag,
		Category:       challenge.Category,
		Type:           challenge.Type,
		GeneratorImage: challenge.GeneratorImage,
		DockerImage:    challenge.DockerImage,
		Port:           challenge.Port,
		Last:           time.Now(),
	}
}
