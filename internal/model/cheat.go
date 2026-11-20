package model

import (
	"crypto/md5"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"
)

type CheatType string

type CheatReasonTmpl string

type CheatReasonType string

const (
	CheaterType    CheatType = "cheater"
	SuspiciousType CheatType = "suspicious"

	ReqWebSameIPTmpl        CheatReasonTmpl = "%s request web with same IP"
	ReqVictimSameIPTmpl     CheatReasonTmpl = "%s request victim with same IP"
	SubmitOtherTeamFlagTmpl CheatReasonTmpl = "Team %d submitted flag of %s for %s in Contest %d"

	ReasonTypeSameWebIPType    CheatReasonType = "same_web_ip"
	ReasonTypeSameVictimIPType CheatReasonType = "same_victim_ip"
	ReasonTypeWrongFlagType    CheatReasonType = "wrong_flag"
)

// Cheat 是一条作弊线索。ReviewedBy 与 ReviewedAt 记录最后一次给出结论的管理员和时间；
// DedupKey 描述线索指向的具体对象，LegacyHash 是引入证据链之前的线索身份，两者都只参与身份计算。
type Cheat struct {
	Time       time.Time         `gorm:"default:null" json:"time"`
	Model      map[string][]uint `gorm:"serializer:json;type:jsonb;default:'{}'" json:"model"`
	IP         string            `json:"ip"`
	Reason     string            `json:"reason"`
	ReasonType CheatReasonType   `gorm:"index" json:"reason_type"`
	Type       CheatType         `json:"type"`
	Hash       string            `gorm:"type:varchar(32);uniqueIndex:idx_cheats_hash_active,where:deleted_at IS NULL;not null" json:"hash"`
	Comment    string            `json:"comment"`
	Evidence   CheatEvidence     `gorm:"serializer:json;type:jsonb;default:'{}'" json:"evidence"`
	BaseModel
	ContestID  uint      `gorm:"index" json:"contest_id"`
	Checked    bool      `gorm:"index" json:"checked"`
	ReviewedBy uint      `gorm:"index" json:"reviewed_by"`
	ReviewedAt time.Time `gorm:"default:null" json:"reviewed_at"`
	DedupKey   string    `gorm:"-" json:"-"`
	LegacyHash string    `gorm:"-" json:"-"`
}

// CheatHash 由比赛、判定类型和线索指向的对象派生线索身份。
// 身份只取决于判定对象，因此重复扫描不会为同一对象生成第二条线索。
func CheatHash(contestID uint, reasonType CheatReasonType, subject string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("v2|%d|%s|%s", contestID, reasonType, subject))))
}

// LegacyCheatHash 复现引入证据链之前的线索身份，使重新扫描把旧线索就地升级，
// 而不是为同一判定再生成一条重复线索。
func LegacyCheatHash(contestID uint, models map[string][]uint, reasonType CheatReasonType, ip, comment string) string {
	keys := make([]string, 0, len(models))
	for key := range models {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	parts := make([]string, 0, len(keys)+4)
	parts = append(parts, strconv.FormatUint(uint64(contestID), 10))
	for _, key := range keys {
		ids := slices.Clone(models[key])
		slices.Sort(ids)
		group := make([]string, 0, len(ids))
		for _, id := range ids {
			group = append(group, strconv.FormatUint(uint64(id), 10))
		}
		parts = append(parts, fmt.Sprintf("%s-%s", key, strings.Join(group, "-,")))
	}
	parts = append(parts, string(reasonType), ip, comment)
	return fmt.Sprintf("%x", md5.Sum([]byte(strings.Join(parts, "-"))))
}
