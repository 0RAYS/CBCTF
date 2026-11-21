package model

import (
	"crypto/sha256"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net/netip"
	"path/filepath"
	"slices"
	"strconv"

	netv1 "k8s.io/api/networking/v1"

	"CBCTF/internal/config"
	"CBCTF/internal/worker"
)

type ChallengeType string

const (
	StaticChallengeType  ChallengeType = "static"
	DynamicChallengeType ChallengeType = "dynamic"
	PodsChallengeType    ChallengeType = "pods"

	AttachmentFileName = "attachment.zip"
	GeneratorFileName  = "generator.zip"
)

// Challenge 题目。
// 题目在比赛中创建，ContestID 为空表示不归属任何比赛。
// 复制出的题目彼此独立，没有关联关系。
type Challenge struct {
	ChallengeFlags  []ChallengeFlag   `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Submissions     []Submission      `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	RandID          string            `gorm:"type:varchar(36);uniqueIndex:idx_challenges_rand_id_active,where:deleted_at IS NULL;not null" json:"rand_id"`
	ContestID       sql.Null[uint]    `gorm:"index" json:"contest_id"`
	Contest         Contest           `json:"-"`
	Name            string            `gorm:"index" json:"name"`
	Description     string            `json:"description"`
	Category        string            `gorm:"index" json:"category"`
	Type            ChallengeType     `gorm:"index" json:"type"`
	Hidden          bool              `gorm:"index" json:"hidden"`
	Attempt         int64             `json:"attempt"`
	Hints           []string          `gorm:"serializer:json;type:jsonb;default:'[]'" json:"hints"`
	Tags            []string          `gorm:"serializer:json;type:jsonb;default:'[]'" json:"tags"`
	GeneratorImage  string            `json:"generator_image"`
	NetworkPolicies NetworkPolicies   `gorm:"type:jsonb" json:"network_policies"`
	Template        ChallengeTemplate `gorm:"serializer:json;type:jsonb;default:'{}'" json:"-"`
	BaseModel
}

// AttachmentCachePath Cache generations include the exact ordered team flags and source revision.
// Old and new flags never publish to the same path, even during a reset race.
func (c Challenge) AttachmentCachePath(teamID uint, flags []string) string {
	return c.AttachmentCachePathForRevision(teamID, flags, worker.SourceRevision(c.GeneratorPath()))
}

func (c Challenge) AttachmentCachePathForRevision(teamID uint, flags []string, revision string) string {
	data, _ := json.Marshal([]any{c.ID, teamID, c.UpdatedAt.UnixNano(), c.GeneratorImage, revision, flags})
	hash := sha256.Sum256(data)
	return filepath.Join(c.BasicDir(), "cache", fmt.Sprintf("%d-%x.zip", teamID, hash))
}

// DataPath 返回数据根目录。配置尚未初始化时返回空字符串。
func DataPath() string {
	if config.Env == nil {
		return ""
	}
	return config.Env.Path
}

func (c Challenge) BasicDir() string {
	return filepath.Join(DataPath(), "challenges", strconv.FormatUint(uint64(c.ID), 10))
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

// NetworkPolicy 描述单个容器的网络策略。
// AllowInternet 为 nil 表示条目没有声明这一项，沿用旧的 CIDR 规则语义。
type NetworkPolicy struct {
	Service       string           `json:"service"`
	AllowInternet *bool            `json:"allow_internet,omitempty"`
	EgressCIDRs   []string         `json:"egress_cidrs,omitempty"`
	Peers         []string         `json:"peers,omitempty"`
	From          []*netv1.IPBlock `json:"from"`
	To            []*netv1.IPBlock `json:"to"`
}

type NetworkPolicies []NetworkPolicy

// DeepCopy 返回网络策略的深拷贝。
// NetworkPolicy 内含 []*netv1.IPBlock 指针，Value() 序列化时会就地裁剪这些块，
// 共享指针会造成跨副本的意外修改。
func (n NetworkPolicies) DeepCopy() NetworkPolicies {
	if n == nil {
		return nil
	}
	result := make(NetworkPolicies, 0, len(n))
	for _, policy := range n {
		result = append(result, NetworkPolicy{
			Service:       policy.Service,
			AllowInternet: copyBool(policy.AllowInternet),
			EgressCIDRs:   append([]string(nil), policy.EgressCIDRs...),
			Peers:         append([]string(nil), policy.Peers...),
			From:          copyIPBlocks(policy.From),
			To:            copyIPBlocks(policy.To),
		})
	}
	return result
}

func copyBool(value *bool) *bool {
	if value == nil {
		return nil
	}
	copied := *value
	return &copied
}

func copyIPBlocks(blocks []*netv1.IPBlock) []*netv1.IPBlock {
	if blocks == nil {
		return nil
	}
	result := make([]*netv1.IPBlock, 0, len(blocks))
	for _, block := range blocks {
		if block == nil {
			continue
		}
		copied := *block
		copied.Except = append([]string(nil), block.Except...)
		result = append(result, &copied)
	}
	return result
}

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
		n[i].EgressCIDRs = slices.DeleteFunc(policy.EgressCIDRs, func(cidr string) bool {
			_, err := netip.ParsePrefix(cidr)
			return err != nil
		})
		n[i].Peers = slices.DeleteFunc(policy.Peers, func(peer string) bool {
			return peer == ""
		})
	}
	return json.Marshal(n)
}

func (n *NetworkPolicies) Scan(value any) error {
	if err := scanJSON(value, n); err != nil {
		return fmt.Errorf("failed to scan NetworkPolicy value")
	}
	return nil
}
