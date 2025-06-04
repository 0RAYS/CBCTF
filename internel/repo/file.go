package repo

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/model"
	"gorm.io/gorm"
)

type FileRepo struct {
	Basic[model.File]
}

type CreateFileOptions struct {
	RandID    string
	Filename  string
	Size      int64
	Path      string
	AdminID   uint
	UserID    uint
	TeamID    uint
	ContestID uint
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
		Suffix:    c.Suffix,
		Hash:      c.Hash,
		Type:      c.Type,
	}
}

type UpdateFileOptions struct {
}

func (u UpdateFileOptions) Convert2Map() map[string]any {
	return make(map[string]any)
}

func InitFileRepo(tx *gorm.DB) *FileRepo {
	return &FileRepo{
		Basic: Basic[model.File]{
			DB: tx,
		},
	}
}

func (f *FileRepo) GetByRandID(randID string) (model.File, bool, string) {
	return f.getUniqueByKey("rand_id", randID)
}

func (f *FileRepo) GetByHash(hash string) (model.File, bool, string) {
	return f.getUniqueByKey("hash", hash)
}

func (f *FileRepo) DeleteByRandID(randIDL ...string) (bool, string) {
	idL := make([]uint, 0)
	for _, randID := range randIDL {
		file, ok, _ := f.GetByRandID(randID)
		if !ok {
			continue
		}
		idL = append(idL, file.ID)
	}
	if ok, msg := f.Delete(idL...); !ok {
		return false, msg
	}
	return true, i18n.Success
}
