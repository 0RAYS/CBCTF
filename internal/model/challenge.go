package model

import (
	"CBCTF/internal/config"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net/netip"
	"path/filepath"
	"slices"
	"strconv"

	netv1 "k8s.io/api/networking/v1"
)

type ChallengeType string

const (
	StaticChallengeType  ChallengeType = "static"
	DynamicChallengeType ChallengeType = "dynamic"
	PodsChallengeType    ChallengeType = "pods"

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
	NetworkPolicies   NetworkPolicies    `gorm:"type:jsonb" json:"network_policies"`
	Template          ChallengeTemplate  `gorm:"type:jsonb" json:"-"`
	BaseModel
}

func (c Challenge) BasicDir() string {
	return filepath.Join(config.Env.Path, "challenges", strconv.FormatUint(uint64(c.ID), 10))
}

// StaticPath 获取静态题目文件的路径
func (c Challenge) StaticPath() string {
	return filepath.Join(c.BasicDir(), AttachmentFileName)
}

// GeneratorPath 获取动态题目生成器的路径
func (c Challenge) GeneratorPath() string {
	return filepath.Join(c.BasicDir(), GeneratorFileName)
}

// AttachmentPath 获取下载时, 题目附件的路径
func (c Challenge) AttachmentPath(teamID uint) string {
	switch c.Type {
	case DynamicChallengeType:
		return filepath.Join(c.BasicDir(), "attachments", strconv.FormatUint(uint64(teamID), 10)+".zip")
	default:
		return c.StaticPath()
	}
}

type NetworkPolicy struct {
	PodKey       string           `json:"pod_key"`
	ContainerKey string           `json:"container_key"`
	From         []*netv1.IPBlock `json:"from"`
	To           []*netv1.IPBlock `json:"to"`
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
