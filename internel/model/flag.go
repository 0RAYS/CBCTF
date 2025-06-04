package model

import (
	"math"
	"regexp"
	"time"
)

const (
	StaticScore      uint = 0
	LinearScore      uint = 1
	LogarithmicScore uint = 2
)

var (
	StaticFlag  = regexp.MustCompile(`static\{(.*?)\}`)
	UUIDFlag    = regexp.MustCompile(`uuid\{(.*?)\}`)
	DynamicFlag = regexp.MustCompile(`dynamic\{(.*?)\}`)
)

type Flag struct {
	ContestID    uint         `json:"contest_id"`
	Contest      Contest      `json:"-"`
	UsageID      uint         `json:"usage_id"`
	Usage        Usage        `json:"-"`
	Answers      []Answer     `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Submissions  []Submission `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Value        string       `json:"value"`
	Score        float64      `gorm:"default:1000" json:"score"`
	CurrentScore float64      `gorm:"default:1000" json:"current_score"`
	Decay        float64      `gorm:"default:50" json:"decay"`
	MinScore     float64      `gorm:"default:100" json:"min_score"`
	ScoreType    uint         `gorm:"default:0" json:"score_type"`
	Solvers      int64        `json:"solvers"`
	Last         time.Time    `gorm:"default:null" json:"last"`
	Blood        UintList     `gorm:"type:json" json:"blood"`
	BaseModel
}

func (f *Flag) CalcCurrentScore(solvers int64) float64 {
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

func (f *Flag) CalcBlood(teamID uint) (float64, int) {
	mapping := []struct {
		value float64
		name  int
	}{
		{0.05, 0},
		{0.03, 1},
		{0.01, 2},
	}
	for i := 0; i < len(f.Blood) && i < 3; i++ {
		if f.Blood[i] == 0 || f.Blood[i] == teamID {
			return mapping[i].value, mapping[i].name
		}
	}
	return 0, -1
}
