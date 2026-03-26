package model

import (
	"CBCTF/internal/config"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net/netip"
	"slices"

	netv1 "k8s.io/api/networking/v1"
)

type ChallengeType string

const (
	StaticChallengeType   ChallengeType = "static"
	QuestionChallengeType ChallengeType = "question"
	DynamicChallengeType  ChallengeType = "dynamic"
	PodsChallengeType     ChallengeType = "pods"

	AttachmentFileName = "attachment.zip"
	GeneratorFileName  = "generator.zip"
)

// Challenge 题库中的挑战
// HasMany ChallengeFlag
// HasMany ContestChallenge
// HasMany Submission
type Challenge struct {
	ChallengeFlags    []ChallengeFlag    `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	ContestChallenges []ContestChallenge `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Submissions       []Submission       `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	RandID            string             `gorm:"type:varchar(36);uniqueIndex:idx_challenges_rand_id_active,where:deleted_at IS NULL;not null" json:"rand_id"`
	Name              string             `gorm:"index" json:"name"`
	Description       string             `json:"description"`
	Category          string             `gorm:"index" json:"category"`
	Type              ChallengeType      `gorm:"index" json:"type"`
	GeneratorImage    string             `json:"generator_image"`
	Options           Options            `gorm:"type:jsonb" json:"options"`
	NetworkPolicies   NetworkPolicies    `gorm:"type:jsonb" json:"network_policies"`
	Template          ChallengeTemplate  `gorm:"type:jsonb" json:"-"`
	TemplateVersion   int                `gorm:"default:1" json:"-"`
	BaseModel
}

func (c Challenge) BasicDir() string {
	return fmt.Sprintf("%s/challenges/%d", config.Env.Path, c.ID)
}

// StaticPath 获取静态题目文件的路径
func (c Challenge) StaticPath() string {
	return fmt.Sprintf("%s/%s", c.BasicDir(), AttachmentFileName)
}

// GeneratorPath 获取动态题目生成器的路径
func (c Challenge) GeneratorPath() string {
	return fmt.Sprintf("%s/%s", c.BasicDir(), GeneratorFileName)
}

// AttachmentPath 获取下载时, 题目附件的路径
func (c Challenge) AttachmentPath(teamID uint) string {
	switch c.Type {
	case DynamicChallengeType:
		return fmt.Sprintf("%s/attachments/%d.zip", c.BasicDir(), teamID)
	default:
		return c.StaticPath()
	}
}

type Option struct {
	RandID  string `json:"rand_id"`
	Content string `json:"content"`
	Correct bool   `json:"correct"`
}

type Options []Option

func (o Options) Value() (driver.Value, error) {
	return json.Marshal(o)
}

func (o *Options) Scan(value any) error {
	if err := scanJSON(value, o); err != nil {
		return fmt.Errorf("failed to scan Options value")
	}
	return nil
}

type NetworkPolicy struct {
	From []*netv1.IPBlock `json:"from"`
	To   []*netv1.IPBlock `json:"to"`
}

type NetworkPolicies []NetworkPolicy

func (n NetworkPolicies) Value() (driver.Value, error) {
	tidy := func(blocks []*netv1.IPBlock) []*netv1.IPBlock {
		return slices.DeleteFunc(blocks, func(b *netv1.IPBlock) bool {
			if b.CIDR == "" {
				return true
			}
			cidr, err := netip.ParsePrefix(b.CIDR)
			if err != nil {
				return true
			}
			b.Except = slices.DeleteFunc(b.Except, func(except string) bool {
				exceptPrefix, err := netip.ParsePrefix(except)
				if err != nil {
					return true
				}
				return !cidr.Contains(exceptPrefix.Addr())
			})
			return false
		})
	}

	for i, policy := range n {
		n[i].From = tidy(policy.From)
		n[i].To = tidy(policy.To)
	}
	return json.Marshal(n)
}

func (n *NetworkPolicies) Scan(value any) error {
	if err := scanJSON(value, n); err != nil {
		return fmt.Errorf("failed to scan NetworkPolicy value")
	}
	return nil
}
