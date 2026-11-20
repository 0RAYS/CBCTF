package db

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
)

type CheatRepo struct {
	BaseRepo[model.Cheat]
}

// UpdateCheatRepo 是一次线索更新需要写入的字段。
// ReviewedBy 与 ReviewedAt 记录最后一次给出结论的管理员和时间。
type UpdateCheatRepo struct {
	Reason     *string
	Type       *model.CheatType
	Checked    *bool
	Comment    *string
	ReviewedBy *uint
	ReviewedAt *time.Time
}

func (u UpdateCheatRepo) Convert2Map() map[string]any {
	options := make(map[string]any)
	if u.Reason != nil {
		options["reason"] = *u.Reason
	}
	if u.Type != nil {
		options["type"] = *u.Type
	}
	if u.Checked != nil {
		options["checked"] = *u.Checked
	}
	if u.Comment != nil {
		options["comment"] = *u.Comment
	}
	if u.ReviewedBy != nil {
		options["reviewed_by"] = *u.ReviewedBy
	}
	if u.ReviewedAt != nil {
		options["reviewed_at"] = *u.ReviewedAt
	}
	return options
}

// RefreshCheatRepo 更新线索的证据链与身份哈希。
// 重复扫描时同一条线索就地刷新，不会因为证据增长而另开一条记录。
type RefreshCheatRepo struct {
	Hash     string
	Evidence model.CheatEvidence
	Time     time.Time
}

func (r RefreshCheatRepo) Convert2Map() map[string]any {
	return map[string]any{
		"hash":     r.Hash,
		"evidence": r.Evidence,
		"time":     r.Time,
	}
}

func InitCheatRepo(tx *gorm.DB) *CheatRepo {
	return &CheatRepo{
		DB: tx,
	}
}

var CheatMutex sync.Map

func (c *CheatRepo) Create(cheat model.Cheat) (model.Cheat, bool, model.RetVal) {
	if cheat.Hash == "" {
		cheat.Hash = model.CheatHash(cheat.ContestID, cheat.ReasonType, cheat.DedupKey)
	}
	target := fmt.Sprintf("hash=%s", cheat.Hash)
	mu, _ := CheatMutex.LoadOrStore(cheat.Hash, &sync.Mutex{})
	mu.(*sync.Mutex).Lock()
	defer mu.(*sync.Mutex).Unlock()

	// The legacy hash recognises records created before detections carried evidence,
	// so a rescan upgrades them in place instead of duplicating them.
	var existing model.Cheat
	for _, hash := range []string{cheat.Hash, cheat.LegacyHash} {
		if hash == "" {
			continue
		}
		res := c.DB.Where("hash = ?", hash).First(&existing)
		if res.Error == nil {
			return existing, false, model.SuccessRetVal()
		} else if !errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.Logger.Warningf("Failed to query Cheat by hash: %s (target=%s)", res.Error, target)
			return model.Cheat{}, false, model.FailRetVal(i18n.Model.Cheat.CreateError, "Cheat", target, res.Error)
		}
	}

	res := c.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&cheat)
	if res.Error != nil {
		log.Logger.Warningf("Failed to create Cheat: %s (target=%s)", res.Error, target)
		return model.Cheat{}, false, model.FailRetVal(i18n.Model.Cheat.CreateError, "Cheat", target, res.Error)
	}
	if res.RowsAffected == 0 {
		if err := c.DB.Where("hash = ?", cheat.Hash).First(&existing).Error; err != nil {
			return model.Cheat{}, false, model.FailRetVal(i18n.Model.Cheat.CreateError, "Cheat", target, err)
		}
		return existing, false, model.SuccessRetVal()
	}
	return cheat, true, model.SuccessRetVal()
}

func (c *CheatRepo) DeleteByContestID(idL ...uint) model.RetVal {
	return c.DeleteByFieldID("contest_id", idL...)
}
