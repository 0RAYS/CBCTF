package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"

	"gorm.io/gorm"
)

type FileRepo struct {
	BaseRepo[model.File]
}

type CreateFileOptions struct {
	RandID   string
	Filename string
	Size     int64
	Path     string
	Model    string
	ModelID  uint
	Suffix   string
	Hash     string
	Type     model.FileType
}

func (c CreateFileOptions) Convert2Model() model.Model {
	return model.File{
		RandID:   c.RandID,
		Filename: c.Filename,
		Size:     c.Size,
		Path:     c.Path,
		Model:    c.Model,
		ModelID:  c.ModelID,
		Suffix:   c.Suffix,
		Hash:     c.Hash,
		Type:     c.Type,
	}
}

func InitFileRepo(tx *gorm.DB) *FileRepo {
	return &FileRepo{
		BaseRepo: BaseRepo[model.File]{
			DB: tx,
		},
	}
}

func (f *FileRepo) Create(options CreateFileOptions) (model.File, model.RetVal) {
	records, ret := f.Get(GetOptions{Conditions: map[string]any{"hash": options.Hash}})
	if ret.OK {
		options.Path = records.Path
	}
	m := options.Convert2Model().(model.File)
	if res := f.DB.Model(&model.File{}).Create(&m); res.Error != nil {
		log.Logger.Warningf("Failed to create File: %s", res.Error)
		return model.File{}, model.RetVal{Msg: i18n.Model.File.CreateError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return m, model.SuccessRetVal()
}

func (f *FileRepo) GetByRandID(randID string, optionsL ...GetOptions) (model.File, model.RetVal) {
	return f.GetByUniqueKey("rand_id", randID, optionsL...)
}

func (f *FileRepo) DeleteByRandID(randIDL ...string) model.RetVal {
	if res := f.DB.Model(&model.File{}).Where("rand_id IN ?", randIDL).Delete(&model.File{}); res.Error != nil {
		log.Logger.Warningf("Failed to delete File: %s", res.Error)
		return model.RetVal{Msg: i18n.Model.File.DeleteError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}
