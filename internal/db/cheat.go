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
	BaseRepo[model.Cheat]
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
	Time               time.Time
}

func (c CreateCheatOptions) Convert2Model() model.Model {
	hash := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf(
		"%s-%d-%d-%d-%d-%d-%s-%s-%s",
		c.Time.Format("2006-01-02"), c.UserID.V, c.TeamID.V, c.ContestID.V, c.ContestChallengeID.V, c.ContestFlagID.V, c.Magic, c.IP, c.Comment,
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
		Time:               c.Time,
		Hash:               hash,
	}
}

type UpdateCheatRepo struct {
	Reason  *string
	Type    *string
	Checked *bool
	Hash    *string
	Comment *string
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
	return options
}

func InitCheatRepo(tx *gorm.DB) *CheatRepo {
	return &CheatRepo{
		BaseRepo: BaseRepo[model.Cheat]{
			DB: tx,
		},
	}
}

func (c *CheatRepo) Create(options CreateCheatOptions) (model.Cheat, model.RetVal) {
	m := options.Convert2Model().(model.Cheat)
	if res := c.DB.Model(&model.Cheat{}).Attrs(m).FirstOrCreate(&m, model.Cheat{Hash: m.Hash}); res.Error != nil {
		log.Logger.Warningf("Failed to create Cheat: %s", res.Error)
		return model.Cheat{}, model.RetVal{Msg: i18n.Model.CreateError, Attr: map[string]any{"Model": model.Cheat{}.GetModelName(), "Error": res.Error.Error()}}
	}
	return m, model.SuccessRetVal()
}
