package model

import (
	"CBCTF/internal/config"
	"CBCTF/internal/utils"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/plugin/optimisticlock"
	"time"
)

const (
	StaticScore      uint = 0
	LinearScore      uint = 1
	LogarithmicScore uint = 2
)

type Usage struct {
	ID             uint                   `json:"id" gorm:"primaryKey"`
	ContestID      uint                   `json:"contest_id"`
	ChallengeID    string                 `json:"challenge_id"`
	Name           string                 `json:"name" gorm:"not null"`
	Desc           string                 `json:"desc"`
	Flag           string                 `json:"flag"`
	Category       string                 `json:"category"`
	GeneratorImage string                 `json:"generator" gorm:"column:generator"`
	DockerImage    string                 `json:"docker" gorm:"column:docker"`
	Port           int32                  `json:"port" gorm:"default:8080"`
	Type           string                 `json:"type" gorm:"default:'static'"`
	Hidden         bool                   `json:"hidden" default:"true"`
	Score          float64                `json:"score" gorm:"default:1000"`
	CurrentScore   float64                `json:"current_score" gorm:"default:1000"`
	ScoreType      uint                   `json:"score_type" gorm:"default:0"`
	MinScore       float64                `json:"min_score" gorm:"default:100"`
	Decay          float64                `json:"decay" gorm:"default:100"`
	Attempt        int64                  `json:"attempt" gorm:"default:0"`
	Solvers        int64                  `json:"solvers" gorm:"default:0"`
	Hints          utils.Strings          `json:"hints"`
	Tags           utils.Strings          `json:"tags"`
	First          uint                   `json:"first" gorm:"default:0"`
	Second         uint                   `json:"second" gorm:"default:0"`
	Third          uint                   `json:"third" gorm:"default:0"`
	Last           time.Time              `json:"last"`
	NetworkPolicy  utils.NetworkPolicy    `json:"network_policy"`
	CreatedAt      time.Time              `json:"-"`
	UpdatedAt      time.Time              `json:"-"`
	DeletedAt      gorm.DeletedAt         `json:"-" gorm:"index"`
	Version        optimisticlock.Version `json:"-" gorm:"default:1"`
}

// BasicDir 获取题目相关文件的目录
func (u *Usage) BasicDir() string {
	return fmt.Sprintf("%s/challenges/%s", config.Env.Path, u.ChallengeID)
}

// StaticPath 获取静态题目文件的路径
func (u *Usage) StaticPath() string {
	return fmt.Sprintf("%s/%s", u.BasicDir(), StaticFile)
}

// GeneratorPath 获取动态题目生成器的路径
func (u *Usage) GeneratorPath() string {
	return fmt.Sprintf("%s/%s", u.BasicDir(), DynamicFile)
}

// AttachmentPath 获取下载时, 题目附件的路径
func (u *Usage) AttachmentPath(teamID uint) string {
	switch u.Type {
	case Dynamic:
		return fmt.Sprintf("%s/attachment/%s/%d.zip", config.Env.Path, u.ChallengeID, teamID)
	default:
		return u.StaticPath()
	}
}

// CalcScore 依据 Solver ScoreType 计算当前分数
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

func (u *Usage) CalcBlood(teamID uint) (float64, string) {
	mapping := []struct {
		value float64
		name  string
		id    uint
	}{
		{0.05, "first", u.First},
		{0.03, "second", u.Second},
		{0.01, "third", u.Third},
	}
	for _, m := range mapping {
		if m.id == 0 || m.id == teamID {
			return m.value, m.name
		}
	}
	return 0, ""
}

func InitUsage(challenge Challenge, contestID uint) Usage {
	usage := Usage{
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
	if challenge.Type == Container {
		// defaultPolicy 允许外部访问, 不允许访问内网, 允许访问外网
		defaultPolicy := utils.NetworkPolicy{
			From: []utils.IPBlock{},
			To: []utils.IPBlock{
				{
					CIDR: "0.0.0.0/0",
					Except: []string{
						"10.0.0.0/8",
						"172.16.0.0/12",
						"192.168.0.0/16",
						"100.64.0.0/10",
					},
				},
			},
		}
		usage.NetworkPolicy = defaultPolicy
	}
	return usage
}
