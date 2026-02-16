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

const (
	StaticChallengeType   = "static"
	QuestionChallengeType = "question"
	DynamicChallengeType  = "dynamic"
	PodsChallengeType     = "pods"

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
	Dockers           []Docker           `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	RandID            string             `gorm:"type:varchar(36);uniqueIndex;not null" json:"rand_id"`
	Name              string             `json:"name"`
	Description       string             `json:"description"`
	Category          string             `json:"category"`
	Type              string             `json:"type"`
	GeneratorImage    string             `json:"generator_image"`
	Options           Options            `gorm:"type:json" json:"options"`
	NetworkPolicies   NetworkPolicies    `gorm:"type:json" json:"network_policies"`
	BaseModel
}

func (c Challenge) TableName() string {
	return "challenges"
}

func (c Challenge) ModelName() string {
	return "Challenge"
}

func (c Challenge) GetBaseModel() BaseModel {
	return c.BaseModel
}

func (c Challenge) UniqueFields() []string {
	return []string{"id", "rand_id"}
}

func (c Challenge) QueryFields() []string {
	return []string{"id", "rand_id", "name", "description", "category", "type", "generator_image"}
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

func (c Challenge) NFSBasicDir() string {
	return fmt.Sprintf("%s/challenges/%d", config.Env.NFS.Path, c.ID)
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
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan Options value")
	}
	return json.Unmarshal(bytes, o)
}

type NetworkPolicy struct {
	From []*netv1.IPBlock `json:"from"`
	To   []*netv1.IPBlock `json:"to"`
}

var DefaultNetworkPolicy = NetworkPolicy{
	To: []*netv1.IPBlock{
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
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan NetworkPolicy value")
	}
	return json.Unmarshal(bytes, n)
}
