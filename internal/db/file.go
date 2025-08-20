package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FileRepo struct {
	BasicRepo[model.File]
}

type CreateFileOptions struct {
	RandID    string
	Filename  string
	Size      int64
	Path      string
	AdminID   *uint
	UserID    *uint
	TeamID    *uint
	ContestID *uint
	OauthID   *uint
	Suffix    string
	Hash      string
	Type      string
}

func (c CreateFileOptions) Convert2Model() model.Model {
	return model.File{
		RandID:    c.RandID,
		Filename:  c.Filename,
		Size:      c.Size,
		Path:      c.Path,
		AdminID:   c.AdminID,
		UserID:    c.UserID,
		TeamID:    c.TeamID,
		ContestID: c.ContestID,
		OauthID:   c.OauthID,
		Suffix:    c.Suffix,
		Hash:      c.Hash,
		Type:      c.Type,
	}
}

type UpdateFileOptions struct{}

func (u UpdateFileOptions) Convert2Map() map[string]any {
	return make(map[string]any)
}

type DiffUpdateFileOptions struct{}

func (d DiffUpdateFileOptions) Convert2Expr() map[string]clause.Expr {
	return nil
}

func InitFileRepo(tx *gorm.DB) *FileRepo {
	return &FileRepo{
		BasicRepo: BasicRepo[model.File]{
			DB: tx,
		},
	}
}

func (f *FileRepo) GetByRandID(randID string, optionsL ...GetOptions) (model.File, bool, string) {
	return f.GetByUniqueKey("rand_id", randID, optionsL...)
}

func (f *FileRepo) GetByHash(hash string, optionsL ...GetOptions) (model.File, bool, string) {
	return f.GetByUniqueKey("hash", hash, optionsL...)
}

func (f *FileRepo) DeleteByRandID(randIDL ...string) (bool, string) {
	if res := f.DB.Model(&model.File{}).Where("rand_id IN ?", randIDL).Delete(&model.File{}); res.Error != nil {
		log.Logger.Warningf("Failed to delete File: %s", res.Error)
		return false, i18n.DeleteFileError
	}
	return true, i18n.Success
}
