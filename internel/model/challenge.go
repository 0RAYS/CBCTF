package model

import (
	"CBCTF/internel/config"
	"CBCTF/internel/i18n"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net"
	"slices"
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
	NetworkPolicies   NetworkPolicies    `gorm:"type:json" json:"network_policies"`
	BasicModel
}

func (c Challenge) GetModelName() string {
	return "Challenge"
}

func (c Challenge) GetVersion() uint {
	return c.Version
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

func (c Challenge) GetForeignKeys() []string {
	return []string{"id"}
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

// AttachmentPath 获取下载时, 题目附件的路径
func (c Challenge) AttachmentPath(teamID uint) string {
	switch c.Type {
	case DynamicChallengeType:
		return fmt.Sprintf("%s/attachments/team-%d/%d.zip", config.Env.Path, teamID, c.ID)
	default:
		return c.StaticPath()
	}
}

type Target struct {
	CIDR   string   `json:"cidr"`
	Except []string `json:"except"`
}

func (t Target) isValidIPBlock() bool {
	if t.CIDR == "" {
		return false
	}
	_, ipNet, err := net.ParseCIDR(t.CIDR)
	if err != nil {
		return false
	}
	for _, ex := range t.Except {
		_, exNet, err := net.ParseCIDR(ex)
		if err != nil {
			return false
		}
		if !ipNet.Contains(exNet.IP) {
			return false
		}
	}
	return true
}

type NetworkPolicy struct {
	From []Target `json:"from"`
	To   []Target `json:"to"`
}

var DefaultNetworkPolicy = NetworkPolicy{
	To: []Target{
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
	for _, p := range n {
		p.From = slices.DeleteFunc(p.From, func(t Target) bool {
			return !t.isValidIPBlock()
		})
		p.To = slices.DeleteFunc(p.From, func(t Target) bool {
			return !t.isValidIPBlock()
		})
	}
	if len(n) == 0 {
		return nil, nil
	}
	return json.Marshal(n)
}

func (n *NetworkPolicies) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan NetworkPolicy value")
	}
	if len(bytes) == 0 {
		*n = nil
		return nil
	}
	return json.Unmarshal(bytes, n)
}
