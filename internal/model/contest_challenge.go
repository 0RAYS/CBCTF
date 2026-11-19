package model

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"

	"CBCTF/internal/worker"
)

// ContestChallenge 比赛中的题目副本。
// 题目加入比赛时会把题库题目深拷贝一份到这里，此后题库的改动不再影响已加入的比赛。
// ChallengeID 仅用于溯源。
type ContestChallenge struct {
	ContestID    uint          `gorm:"index" json:"contest_id"`
	Contest      Contest       `json:"-"`
	ChallengeID  uint          `gorm:"index" json:"challenge_id"`
	Challenge    Challenge     `json:"-"`
	RandID       string        `gorm:"type:varchar(36);uniqueIndex:idx_contest_challenges_rand_id_active,where:deleted_at IS NULL" json:"rand_id"`
	ContestFlags []ContestFlag `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Submissions  []Submission  `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Type         ChallengeType `gorm:"index" json:"type"`
	Category     string        `gorm:"index" json:"category"`
	Hidden       bool          `gorm:"index" json:"hidden"`
	Attempt      int64         `json:"attempt"`
	Hints        []string      `gorm:"serializer:json;type:jsonb;default:'[]'" json:"hints"`
	Tags         []string      `gorm:"serializer:json;type:jsonb;default:'[]'" json:"tags"`
	// 以下字段为题库题目的副本，题库改动不会同步到这里
	Template        ChallengeTemplate `gorm:"serializer:json;type:jsonb;default:'{}'" json:"-"`
	GeneratorImage  string            `json:"generator_image"`
	NetworkPolicies NetworkPolicies   `gorm:"type:jsonb" json:"network_policies"`
	BaseModel
}

// BasicDir 比赛题目副本的存储目录
func (c ContestChallenge) BasicDir() string {
	return filepath.Join(DataPath(), "contest_challenges", strconv.FormatUint(uint64(c.ID), 10))
}

// StaticPath 静态附件的路径
func (c ContestChallenge) StaticPath() string {
	return filepath.Join(c.BasicDir(), AttachmentFileName)
}

// GeneratorPath 生成器的路径
func (c ContestChallenge) GeneratorPath() string {
	return filepath.Join(c.BasicDir(), GeneratorFileName)
}

// AttachmentPath 下载时题目附件的路径
func (c ContestChallenge) AttachmentPath(teamID uint) string {
	switch c.Type {
	case DynamicChallengeType:
		return filepath.Join(c.BasicDir(), "attachments", strconv.FormatUint(uint64(teamID), 10)+".zip")
	default:
		return c.StaticPath()
	}
}

// AttachmentCachePath Cache generations include the exact ordered team flags and source revision.
// Old and new flags never publish to the same path, even during a reset race.
func (c ContestChallenge) AttachmentCachePath(teamID uint, flags []string) string {
	return c.AttachmentCachePathForRevision(teamID, flags, worker.SourceRevision(c.GeneratorPath()))
}

func (c ContestChallenge) AttachmentCachePathForRevision(teamID uint, flags []string, revision string) string {
	data, _ := json.Marshal([]any{c.ID, teamID, c.UpdatedAt.UnixNano(), c.GeneratorImage, revision, flags})
	hash := sha256.Sum256(data)
	return filepath.Join(c.BasicDir(), "cache", fmt.Sprintf("%d-%x.zip", teamID, hash))
}
