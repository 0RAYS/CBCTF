package repo

import (
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"gorm.io/gorm"
)

type FileRepo struct {
	Repo[model.File]
}

type CreateFileOptions struct {
	ID       string
	Filename string
	Size     int64
	Path     string
	Uploader uint
	Suffix   string
	Hash     string
	Type     string
}

func InitFileRepo(tx *gorm.DB) *FileRepo {
	return &FileRepo{Repo: Repo[model.File]{DB: tx, Model: "File"}}
}

//func (f *FileRepo) Create(options CreateFileOptions) (model.File, bool, string) {
//	file, err := utils.S2S[model.File](options)
//	if err != nil {
//		return model.File{}, false, "Options2ModelError"
//	}
//	res := f.DB.Model(&model.File{}).Create(&file)
//	if res.Error != nil {
//		log.Logger.Warningf("Failed to create File: %s", res.Error)
//		return model.File{}, false, "CreateFileError"
//	}
//	return file, true, "Success"
//}

func (f *FileRepo) getByUniqueKey(key string, value interface{}) (model.File, bool, string) {
	switch key {
	// 虽然 hash 并不是唯一的，但是并不影响功能
	case "id", "hash":
		value = value.(string)
	default:
		return model.File{}, false, "UnsupportedKey"
	}
	var file model.File
	res := f.DB.Model(&model.File{}).Where(key+" = ?", value).Find(&file).Limit(1)
	if res.RowsAffected == 0 {
		return model.File{}, false, "FileNotFound"
	}
	return file, true, "Success"
}

func (f *FileRepo) GetByID(id string) (model.File, bool, string) {
	return f.getByUniqueKey("id", id)
}

func (f *FileRepo) GetByHash(hash string) (model.File, bool, string) {
	return f.getByUniqueKey("hash", hash)
}

func (f *FileRepo) Count(t string) (int64, bool, string) {
	var count int64
	res := f.DB.Model(&model.File{}).Where("type = ?", t).Count(&count)
	if res.Error != nil {
		log.Logger.Warningf("Failed to count File: %s", res.Error)
		return 0, false, "CountModelError"
	}
	return count, true, "Success"
}

func (f *FileRepo) GetAll(t string, limit, offset int) ([]model.File, int64, bool, string) {
	var (
		files          = make([]model.File, 0)
		count, ok, msg = f.Count(t)
	)
	if !ok {
		return files, count, false, msg
	}
	res := f.DB.Model(&model.File{}).Where("type = ?", t).Find(&files).Limit(limit).Offset(offset)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get File: %s", res.Error)
		return files, 0, false, "GetFileError"
	}
	return files, count, true, "Success"
}

func (f *FileRepo) Delete(idL ...string) (bool, string) {
	res := f.DB.Model(&model.File{}).Where("id IN ?", idL).Delete(&model.File{})
	if res.Error != nil {
		log.Logger.Warningf("Failed to delete File: %s", res.Error)
		return false, "DeleteFileError"
	}
	return true, "Success"
}
