package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"database/sql"

	"gorm.io/gorm"
)

type FileRepo struct {
	BasicRepo[model.File]
}

type CreateFileOptions struct {
	RandID      string
	Filename    string
	Size        int64
	Path        string
	AdminID     sql.Null[uint]
	UserID      sql.Null[uint]
	TeamID      sql.Null[uint]
	ContestID   sql.Null[uint]
	OauthID     sql.Null[uint]
	ChallengeID sql.Null[uint]
	Suffix      string
	Hash        string
	Type        string
}

func (c CreateFileOptions) Convert2Model() model.Model {
	return model.File{
		RandID:      c.RandID,
		Filename:    c.Filename,
		Size:        c.Size,
		Path:        c.Path,
		AdminID:     c.AdminID,
		UserID:      c.UserID,
		TeamID:      c.TeamID,
		ContestID:   c.ContestID,
		OauthID:     c.OauthID,
		ChallengeID: c.ChallengeID,
		Suffix:      c.Suffix,
		Hash:        c.Hash,
		Type:        c.Type,
	}
}

func InitFileRepo(tx *gorm.DB) *FileRepo {
	return &FileRepo{
		BasicRepo: BasicRepo[model.File]{
			DB: tx,
		},
	}
}

func (f *FileRepo) Create(options CreateFileOptions) (model.File, bool, string) {
	records, ok, _ := f.Get(GetOptions{Conditions: map[string]any{"hash": options.Hash}})
	if ok {
		options.Path = records.Path
	}
	m := options.Convert2Model().(model.File)
	if res := f.DB.Model(&model.File{}).Create(&m); res.Error != nil {
		log.Logger.Warningf("Failed to create File: %s", res.Error)
		return model.File{}, false, i18n.CreateFileError
	}
	return m, true, i18n.Success
}

func (f *FileRepo) GetByRandID(randID string, optionsL ...GetOptions) (model.File, bool, string) {
	return f.GetByUniqueKey("rand_id", randID, optionsL...)
}

func (f *FileRepo) DeleteByRandID(randIDL ...string) (bool, string) {
	if res := f.DB.Model(&model.File{}).Where("rand_id IN ?", randIDL).Delete(&model.File{}); res.Error != nil {
		log.Logger.Warningf("Failed to delete File: %s", res.Error)
		return false, i18n.DeleteFileError
	}
	return true, i18n.Success
}
