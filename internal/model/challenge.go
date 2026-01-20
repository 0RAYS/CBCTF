package model

import (
	"CBCTF/internal/config"
	"CBCTF/internal/i18n"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net"
	"slices"

	netv1 "k8s.io/api/networking/v1"
)

const (
	StaticChallengeType   = "static"
	QuestionChallengeType = "question"
	DynamicChallengeType  = "dynamic"
	PodsChallengeType     = "pods"

	AttachmentFile = "attachment.zip"
	GeneratorFile  = "generator.zip"
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
	Desc              string             `json:"desc"`
	Category          string             `json:"category"`
	Type              string             `json:"type"`
	GeneratorImage    string             `json:"generator_image"`
	Options           Options            `gorm:"type:json" json:"options"`
	NetworkPolicies   NetworkPolicies    `gorm:"type:json" json:"network_policies"`
	BaseModel
}

func (c Challenge) GetModelName() string {
	return "Challenge"
}

func (c Challenge) GetBaseModel() BaseModel {
	return c.BaseModel
}

func (c Challenge) CreateErrorString() string {
	return i18n.CreateChallengeError
}

func (c Challenge) DeleteErrorString() string {
	return i18n.DeleteChallengeError
}

func (c Challenge) GetErrorString() string {
	return i18n.GetChallengeError
}

func (c Challenge) NotFoundErrorString() string {
	return i18n.ChallengeNotFound
}

func (c Challenge) UpdateErrorString() string {
	return i18n.UpdateChallengeError
}

func (c Challenge) GetUniqueKey() []string {
	return []string{"id", "rand_id"}
}

func (c Challenge) BasicDir() string {
	return fmt.Sprintf("%s/challenges/%d", config.Env.Path, c.ID)
}

// StaticPath 获取静态题目文件的路径
func (c Challenge) StaticPath() string {
	return fmt.Sprintf("%s/%s", c.BasicDir(), AttachmentFile)
}

// GeneratorPath 获取动态题目生成器的路径
func (c Challenge) GeneratorPath() string {
	return fmt.Sprintf("%s/%s", c.BasicDir(), GeneratorFile)
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
			_, cidr, err := net.ParseCIDR(b.CIDR)
			if err != nil {
				return true
			}
			b.Except = slices.DeleteFunc(b.Except, func(except string) bool {
				_, exceptCidr, err := net.ParseCIDR(except)
				if err != nil {
					return true
				}
				return !cidr.Contains(exceptCidr.IP)
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
