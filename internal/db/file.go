package db

import (
	"fmt"

	"gorm.io/gorm"

	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
)

type FileRepo struct {
	BaseRepo[model.File]
}

func InitFileRepo(tx *gorm.DB) *FileRepo {
	return &FileRepo{
		DB: tx,
	}
}

func (f *FileRepo) Create(file model.File) (model.File, model.RetVal) {
	records, ret := f.Get(GetOptions{Conditions: map[string]any{"hash": file.Hash}})
	if ret.OK {
		file.Path = records.Path
	}
	if res := f.DB.Model(&model.File{}).Create(&file); res.Error != nil {
		target := fmt.Sprintf("hash=%s", file.Hash)
		log.Logger.Warningf("Failed to create File: %s (target=%s)", res.Error, target)
		return model.File{}, model.FailRetVal(i18n.Model.File.CreateError, "File", target, res.Error)
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
		target := fmt.Sprintf("rand_id=%v", randIDL)
		log.Logger.Warningf("Failed to get Files by rand IDs %v: %s (target=%s)", randIDL, res.Error, target)
		return model.FailRetVal(i18n.Model.File.DeleteError, "File", target, res.Error)
	}
	return f.Delete(fileIDL...)
}
