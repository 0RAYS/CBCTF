package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"crypto/md5"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type CheatRepo struct {
	BaseRepo[model.Cheat]
}

type CreateCheatOptions struct {
	ContestID  uint
	Model      model.CheatRefModel
	Magic      string
	IP         string
	Reason     string
	ReasonType string
	Type       string
	Checked    bool
	Comment    string
	Time       time.Time
}

func (c CreateCheatOptions) Convert2Model() model.Model {
	keys := make([]string, 0)
	for k := range c.Model {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	hash := fmt.Sprintf("%d-", c.ContestID)
	for _, k := range keys {
		ids := make([]uint, len(c.Model[k]))
		copy(ids, c.Model[k])
		sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
		hash += fmt.Sprintf("%s-", k)
		for _, id := range ids {
			hash += fmt.Sprintf("%d-,", id)
		}
		hash = strings.TrimSuffix(hash, "-")
	}
	hash = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s-%s-%s-%s-%s", hash, c.ReasonType, c.Magic, c.IP, c.Comment))))
	return model.Cheat{
		ContestID:  c.ContestID,
		Model:      c.Model,
		Magic:      c.Magic,
		IP:         c.IP,
		Reason:     c.Reason,
		ReasonType: c.ReasonType,
		Type:       c.Type,
		Checked:    c.Checked,
		Comment:    c.Comment,
		Time:       c.Time,
		Hash:       hash,
	}
}

type UpdateCheatRepo struct {
	Reason  *string
	Type    *string
	Checked *bool
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
		var mysqlErr *mysql.MySQLError
		if errors.As(res.Error, &mysqlErr) && mysqlErr.Number == 1062 {
			// Duplicate key: fallback to query existing record
			var existing model.Cheat
			if err := c.DB.Where("hash = ?", m.Hash).First(&existing).Error; err == nil {
				return existing, model.SuccessRetVal()
			}
		}
		log.Logger.Warningf("Failed to create Cheat: %s", res.Error)
		return model.Cheat{}, model.RetVal{Msg: i18n.Model.CreateError, Attr: map[string]any{"Model": model.Cheat{}.ModelName(), "Error": res.Error.Error()}}
	}
	return m, model.SuccessRetVal()
}
