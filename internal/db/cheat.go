package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"crypto/md5"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"

	"gorm.io/gorm"
)

type CheatRepo struct {
	BaseRepo[model.Cheat]
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

func (c *CheatRepo) Create(cheat model.Cheat) (model.Cheat, model.RetVal) {
	keys := make([]string, 0)
	for k := range cheat.Model {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	hash := strconv.FormatUint(uint64(cheat.ContestID), 10) + "-"
	for _, k := range keys {
		ids := make([]uint, len(cheat.Model[k]))
		copy(ids, cheat.Model[k])
		sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
		hash += fmt.Sprintf("%s-", k)
		for _, id := range ids {
			hash += strconv.FormatUint(uint64(id), 10) + "-,"
		}
		hash = strings.TrimSuffix(hash, "-")
	}
	cheat.Hash = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s-%s-%s-%s-%s", hash, cheat.ReasonType, cheat.Magic, cheat.IP, cheat.Comment))))
	mu, _ := CheatMutex.LoadOrStore(cheat.Hash, &sync.Mutex{})
	mu.(*sync.Mutex).Lock()
	defer mu.(*sync.Mutex).Unlock()

	var existing model.Cheat
	if res := c.DB.Where("hash = ?", cheat.Hash).First(&existing); res.Error == nil {
		return existing, model.SuccessRetVal()
	} else if !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		log.Logger.Warningf("Failed to query Cheat by hash: %s", res.Error)
		return model.Cheat{}, model.RetVal{Msg: i18n.Model.Cheat.CreateError, Attr: map[string]any{"Error": res.Error.Error()}}
	}

	if res := c.DB.Create(&cheat); res.Error != nil {
		log.Logger.Warningf("Failed to create Cheat: %s", res.Error)
		return model.Cheat{}, model.RetVal{Msg: i18n.Model.Cheat.CreateError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return cheat, model.SuccessRetVal()
}

func (c *CheatRepo) DeleteByContestID(idL ...uint) model.RetVal {
	return c.DeleteByFieldID("contest_id", idL...)
}
