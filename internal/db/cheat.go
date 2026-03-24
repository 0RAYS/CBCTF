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
	"sync"
	"time"

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
	ReasonType model.CheatReasonType
	Type       model.CheatType
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
	Type    *model.CheatType
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

var CheatMutex sync.Map

func (c *CheatRepo) Create(options CreateCheatOptions) (model.Cheat, model.RetVal) {
	m := options.Convert2Model().(model.Cheat)

	mu, _ := CheatMutex.LoadOrStore(m.Hash, &sync.Mutex{})
	mu.(*sync.Mutex).Lock()
	defer mu.(*sync.Mutex).Unlock()

	var existing model.Cheat
	if res := c.DB.Where("hash = ?", m.Hash).First(&existing); res.Error == nil {
		return existing, model.SuccessRetVal()
	} else if !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		log.Logger.Warningf("Failed to query Cheat by hash: %s", res.Error)
		return model.Cheat{}, model.RetVal{Msg: i18n.Model.Cheat.CreateError, Attr: map[string]any{"Error": res.Error.Error()}}
	}

	if res := c.DB.Create(&m); res.Error != nil {
		log.Logger.Warningf("Failed to create Cheat: %s", res.Error)
		return model.Cheat{}, model.RetVal{Msg: i18n.Model.Cheat.CreateError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return m, model.SuccessRetVal()
}

func (c *CheatRepo) DeleteByContestID(contestID uint) model.RetVal {
	res := c.DB.Model(&model.Cheat{}).Where("contest_id = ?", contestID).Delete(&model.Cheat{})
	if res.Error != nil {
		log.Logger.Warningf("Failed to delete Cheat: %s", res.Error)
		return model.RetVal{Msg: i18n.Model.Cheat.DeleteError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}
