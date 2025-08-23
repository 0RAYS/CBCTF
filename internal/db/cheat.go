package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"crypto/md5"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type CheatRepo struct {
	BasicRepo[model.Cheat]
}

type CreateCheatOptions struct {
	UserID             sql.Null[uint]
	TeamID             sql.Null[uint]
	ContestID          sql.Null[uint]
	ContestChallengeID sql.Null[uint]
	ContestFlagID      sql.Null[uint]
	Magic              string
	IP                 string
	Reason             string
	Type               string
	Checked            bool
	Comment            string
	References         model.UintList
}

func (c CreateCheatOptions) Convert2Model() model.Model {
	hash := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf(
		"%s-%d-%d-%d-%d-%d-%s-%s-%s",
		time.Now().Format("2006-01-02"), c.UserID.V, c.TeamID.V, c.ContestID.V, c.ContestChallengeID.V, c.ContestFlagID.V, c.Magic, c.IP, c.Comment,
	))))
	return model.Cheat{
		UserID:             c.UserID,
		TeamID:             c.TeamID,
		ContestID:          c.ContestID,
		ContestChallengeID: c.ContestChallengeID,
		ContestFlagID:      c.ContestFlagID,
		Magic:              c.Magic,
		IP:                 c.IP,
		Reason:             c.Reason,
		Type:               c.Type,
		Checked:            c.Checked,
		Comment:            c.Comment,
		References:         c.References,
		Hash:               hash,
	}
}

type UpdateCheatRepo struct {
	Reason     *string
	Type       *string
	Checked    *bool
	Hash       *string
	Comment    *string
	References *model.UintList
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
	if u.Hash != nil {
		options["hash"] = *u.Hash
	}
	if u.Comment != nil {
		options["comment"] = *u.Comment
	}
	if u.References != nil {
		options["references"] = u.References
	}
	return options
}

func InitCheatRepo(tx *gorm.DB) *CheatRepo {
	return &CheatRepo{
		BasicRepo: BasicRepo[model.Cheat]{
			DB: tx,
		},
	}
}

func (c *CheatRepo) Create(options CreateCheatOptions) (model.Cheat, bool, string) {
	m := options.Convert2Model().(model.Cheat)
	if cheat, ok, _ := c.GetByHash(m.Hash); ok && !cheat.Checked {
		return cheat, true, i18n.Success
	}
	if res := c.DB.Model(&model.Cheat{}).Create(&m); res.Error != nil {
		log.Logger.Warningf("Failed to create Cheat: %s", res.Error)
		return model.Cheat{}, false, i18n.CreateCheatError
	}
	return m, true, i18n.Success
}

func (c *CheatRepo) GetByHash(hash string, optionsL ...GetOptions) (model.Cheat, bool, string) {
	return c.GetByUniqueKey("hash", hash, optionsL...)
}
