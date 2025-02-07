package model

import (
	"gorm.io/gorm"
	"time"
)

var StaticScore uint = 0
var LinearScore uint = 1
var LogarithmicScore uint = 2

type Usage struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	ContestID    uint           `json:"contest_id"`
	ChallengeID  string         `json:"challenge_id"`
	Hidden       bool           `json:"hidden" default:"true"`
	Score        int64          `json:"score" gorm:"default:100"`
	CurrentScore int64          `json:"current_score" gorm:"default:100"`
	ScoreType    uint           `json:"score_type" gorm:"default:0"`
	MinScore     int64          `json:"min_score" gorm:"default:10"`
	Decay        int64          `json:"decay" gorm:"default:0"`
	Flag         string         `json:"flag"`
	Attempt      int64          `json:"attempt" gorm:"default:0"`
	Solvers      int64          `json:"solvers" gorm:"default:0"`
	Hints        string         `json:"hints"`
	Tags         string         `json:"tags"`
	First        uint           `json:"first" gorm:"default:0"`
	Second       uint           `json:"second" gorm:"default:0"`
	Third        uint           `json:"third" gorm:"default:0"`
	Last         time.Time      `json:"last"`
	CreatedAt    time.Time      `json:"-"`
	UpdatedAt    time.Time      `json:"-"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

func (u *Usage) CalcScore(solvers int64) int64 {
	var calc int64 = 0
	switch u.ScoreType {
	case StaticScore:
		calc = u.Score
	case LinearScore:
		calc = u.Score - solvers*u.Decay
	case LogarithmicScore:
		calc = (((u.MinScore - u.Score) / (u.Decay * u.Decay)) * (solvers * solvers)) + u.Score
	default:
		calc = u.Score
	}
	if calc < u.MinScore {
		calc = u.MinScore
	}
	return calc
}

func InitUsage(challengeID string, contestID uint, flag string) Usage {
	return Usage{
		ContestID:   contestID,
		ChallengeID: challengeID,
		Flag:        flag,
	}
}
