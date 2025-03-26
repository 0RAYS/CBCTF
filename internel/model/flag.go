package model

import (
	"gorm.io/gorm"
	"math"
	"time"
)

const (
	StaticScore      uint = 0
	LinearScore      uint = 1
	LogarithmicScore uint = 2
)

type Flag struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	ContestID    uint           `json:"contest_id"`
	Contest      Contest        `json:"-"`
	UsageID      string         `json:"usage_id"`
	Usage        Usage          `json:"-"`
	Value        string         `json:"value"`
	Score        float64        `gorm:"default:1000" json:"score"`
	CurrentScore float64        `gorm:"default:1000" json:"current_score"`
	Decay        float64        `gorm:"default:50" json:"decay"`
	MinScore     float64        `gorm:"default:100" json:"min_score"`
	ScoreType    uint           `gorm:"default:0" json:"score_type"`
	Solvers      int64          `json:"solvers"`
	Attempt      int64          `json:"attempt"`
	Blood        Uints          `gorm:"type:json" json:"blood"`
	Answers      []Answer       `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Submissions  []Submission   `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Last         time.Time      `json:"last"`
	CreatedAt    time.Time      `json:"-"`
	UpdatedAt    time.Time      `json:"-"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Version      uint           `gorm:"default:1" json:"-"`
}

func (f *Flag) CalcNewScore(solvers int64) float64 {
	var calc float64 = 0
	switch f.ScoreType {
	case StaticScore:
		calc = f.Score
	case LinearScore:
		calc = f.Score - float64(solvers)*f.Decay
	case LogarithmicScore:
		calc = (f.MinScore-f.Score)/(f.Decay*f.Decay)*float64(solvers*solvers) + f.Score
	default:
		calc = f.Score
	}
	if calc < f.MinScore {
		calc = f.MinScore
	}
	calc = math.Trunc(calc*100) / 100
	return calc
}

func (f *Flag) CalcBlood(teamID uint) (float64, string) {
	mapping := []struct {
		value float64
		name  string
	}{
		{0.05, "first"},
		{0.03, "second"},
		{0.01, "third"},
	}
	for i := 0; i < len(f.Blood) && i < 3; i++ {
		if f.Blood[i] == 0 || f.Blood[i] == teamID {
			return mapping[i].value, mapping[i].name
		}
	}
	return 0, ""
}
