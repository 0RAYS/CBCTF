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

func InitFileRepo(tx *gorm.DB) *FileRepo {
	return &FileRepo{
		BaseRepo: BaseRepo[model.File]{
			DB: tx,
		},
	}
}

func (f *FileRepo) Create(file model.File) (model.File, model.RetVal) {
	records, ret := f.Get(GetOptions{Conditions: map[string]any{"hash": file.Hash}})
	if ret.OK {
		file.Path = records.Path
	}
	if res := f.DB.Model(&model.File{}).Create(&file); res.Error != nil {
		log.Logger.Warningf("Failed to create File: %s", res.Error)
		return model.File{}, model.RetVal{Msg: i18n.Model.File.CreateError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return file, model.SuccessRetVal()
}

func (f *FileRepo) GetByRandID(randID string, optionsL ...GetOptions) (model.File, model.RetVal) {
	return f.GetByUniqueField("rand_id", randID, optionsL...)
}

func (f *FileRepo) DeleteByRandID(randIDL ...string) model.RetVal {
	if len(randIDL) == 0 {
		return model.SuccessRetVal()
	}
	var fileIDL []uint
	if res := f.DB.Model(&model.File{}).Where("rand_id IN ?", randIDL).Pluck("id", &fileIDL); res.Error != nil {
		log.Logger.Warningf("Failed to get Files by rand IDs %v: %s", randIDL, res.Error)
		return model.RetVal{Msg: i18n.Model.File.DeleteError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return f.Delete(fileIDL...)
}
