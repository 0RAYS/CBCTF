package model

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"path/filepath"

	"CBCTF/internal/worker"
)

// GeneratorSource 生成器与附件所需的题目来源。
// 题库题目（Challenge）与比赛题目副本（ContestChallenge）都实现该接口，
// 生成器与附件相关的代码只依赖它。
type GeneratorSource interface {
	// SourceID 题库题目为题库 id，比赛题目副本为比赛题目 id
	SourceID() uint
	SourceName() string
	SourceType() ChallengeType
	SourceImage() string
	// SourceDir 题目基础目录，附件与生成器都放在这里
	SourceDir() string
	SourceGeneratorPath() string
	SourceAttachmentPath(teamID uint) string
	SourceAttachmentCachePath(teamID uint, flags []string) string
	SourceAttachmentCachePathForRevision(teamID uint, flags []string, revision string) string
	// SourceChallengeID 溯源用的题库题目 id
	SourceChallengeID() uint
	// SourceContestChallengeID 比赛题目 id，为 0 表示不是比赛题目
	SourceContestChallengeID() uint
	// SourceContestID 比赛 id，为 0 表示不属于某场比赛
	SourceContestID() uint
	SourceTarget() GeneratorTarget
}

func (c Challenge) SourceID() uint { return c.ID }

func (c Challenge) SourceName() string { return c.Name }

func (c Challenge) SourceType() ChallengeType { return c.Type }

func (c Challenge) SourceImage() string { return c.GeneratorImage }

func (c Challenge) SourceDir() string { return c.BasicDir() }

func (c Challenge) SourceGeneratorPath() string { return c.GeneratorPath() }

func (c Challenge) SourceAttachmentPath(teamID uint) string { return c.AttachmentPath(teamID) }

func (c Challenge) SourceAttachmentCachePath(teamID uint, flags []string) string {
	return c.AttachmentCachePath(teamID, flags)
}

func (c Challenge) SourceAttachmentCachePathForRevision(teamID uint, flags []string, revision string) string {
	return c.AttachmentCachePathForRevision(teamID, flags, revision)
}

func (c Challenge) SourceChallengeID() uint { return c.ID }

func (c Challenge) SourceContestChallengeID() uint { return 0 }

func (c Challenge) SourceContestID() uint { return 0 }

func (c Challenge) SourceTarget() GeneratorTarget { return NewGeneratorTarget(c.ID) }

func (c ContestChallenge) SourceID() uint { return c.ID }

func (c ContestChallenge) SourceName() string { return c.Name }

func (c ContestChallenge) SourceType() ChallengeType { return c.Type }

func (c ContestChallenge) SourceImage() string { return c.GeneratorImage }

func (c ContestChallenge) SourceDir() string { return c.BasicDir() }

func (c ContestChallenge) SourceGeneratorPath() string { return c.GeneratorPath() }

func (c ContestChallenge) SourceAttachmentPath(teamID uint) string { return c.AttachmentPath(teamID) }

func (c ContestChallenge) SourceAttachmentCachePath(teamID uint, flags []string) string {
	return c.AttachmentCachePath(teamID, flags)
}

func (c ContestChallenge) SourceAttachmentCachePathForRevision(teamID uint, flags []string, revision string) string {
	return c.AttachmentCachePathForRevision(teamID, flags, revision)
}

func (c ContestChallenge) SourceChallengeID() uint { return c.ChallengeID }

func (c ContestChallenge) SourceContestChallengeID() uint { return c.ID }

func (c ContestChallenge) SourceContestID() uint { return c.ContestID }

func (c ContestChallenge) SourceTarget() GeneratorTarget {
	return NewContestGeneratorTarget(c.ContestID, c.ID)
}

// GeneratorTarget 生成器池的目标标识。生成器池按比赛题目划分，否则同一题库题目
// 在多场比赛加入时会互相顶掉。
type GeneratorTarget struct {
	ContestID          uint
	ChallengeID        uint
	ContestChallengeID uint
}

func NewGeneratorTarget(challengeID uint) GeneratorTarget {
	return GeneratorTarget{ChallengeID: challengeID}
}

func NewContestGeneratorTarget(contestID, contestChallengeID uint) GeneratorTarget {
	return GeneratorTarget{ContestID: contestID, ContestChallengeID: contestChallengeID}
}

// GeneratorSourceSnapshot 生成器与附件任务的题目快照。
// 任务在队列中异步执行，无法再访问数据库，因此快照携带生成器所需的全部字段。
type GeneratorSourceSnapshot struct {
	ID                 uint          `msgpack:"id"`
	ChallengeID        uint          `msgpack:"challenge_id"`
	ContestChallengeID uint          `msgpack:"contest_challenge_id"`
	ContestID          uint          `msgpack:"contest_id"`
	Name               string        `msgpack:"name"`
	Type               ChallengeType `msgpack:"type"`
	Image              string        `msgpack:"image"`
	Dir                string        `msgpack:"dir"`
	GeneratorFile      string        `msgpack:"generator_file"`
	UpdatedAt          int64         `msgpack:"updated_at"`
}

func NewChallengeSource(c Challenge) GeneratorSourceSnapshot {
	return GeneratorSourceSnapshot{
		ID:            c.ID,
		ChallengeID:   c.ID,
		Name:          c.Name,
		Type:          c.Type,
		Image:         c.GeneratorImage,
		Dir:           c.BasicDir(),
		GeneratorFile: c.GeneratorPath(),
		UpdatedAt:     c.UpdatedAt.UnixNano(),
	}
}

func NewContestChallengeSource(c ContestChallenge) GeneratorSourceSnapshot {
	return GeneratorSourceSnapshot{
		ID:                 c.ID,
		ChallengeID:        c.ChallengeID,
		ContestChallengeID: c.ID,
		ContestID:          c.ContestID,
		Name:               c.Name,
		Type:               c.Type,
		Image:              c.GeneratorImage,
		Dir:                c.BasicDir(),
		GeneratorFile:      c.GeneratorPath(),
		UpdatedAt:          c.UpdatedAt.UnixNano(),
	}
}

func (s GeneratorSourceSnapshot) SourceID() uint { return s.ID }

func (s GeneratorSourceSnapshot) SourceName() string { return s.Name }

func (s GeneratorSourceSnapshot) SourceType() ChallengeType { return s.Type }

func (s GeneratorSourceSnapshot) SourceImage() string { return s.Image }

func (s GeneratorSourceSnapshot) SourceDir() string { return s.Dir }

func (s GeneratorSourceSnapshot) SourceGeneratorPath() string { return s.GeneratorFile }

func (s GeneratorSourceSnapshot) SourceAttachmentPath(teamID uint) string {
	if s.Type == DynamicChallengeType {
		return filepath.Join(s.Dir, "attachments", fmt.Sprint(teamID)+".zip")
	}
	return filepath.Join(s.Dir, AttachmentFileName)
}

func (s GeneratorSourceSnapshot) SourceAttachmentCachePath(teamID uint, flags []string) string {
	return s.SourceAttachmentCachePathForRevision(teamID, flags, worker.SourceRevision(s.GeneratorFile))
}

func (s GeneratorSourceSnapshot) SourceAttachmentCachePathForRevision(teamID uint, flags []string, revision string) string {
	data, _ := json.Marshal([]any{s.ID, teamID, s.UpdatedAt, s.Image, revision, flags})
	hash := sha256.Sum256(data)
	return filepath.Join(s.Dir, "cache", fmt.Sprintf("%d-%x.zip", teamID, hash))
}

func (s GeneratorSourceSnapshot) SourceChallengeID() uint { return s.ChallengeID }

func (s GeneratorSourceSnapshot) SourceContestChallengeID() uint { return s.ContestChallengeID }

func (s GeneratorSourceSnapshot) SourceContestID() uint { return s.ContestID }

func (s GeneratorSourceSnapshot) SourceTarget() GeneratorTarget {
	if s.ContestChallengeID > 0 {
		return NewContestGeneratorTarget(s.ContestID, s.ContestChallengeID)
	}
	return NewGeneratorTarget(s.ChallengeID)
}
