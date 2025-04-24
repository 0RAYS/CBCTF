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
	ID        string
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

func InitFileRepo(tx *gorm.DB) *FileRepo {
	return &FileRepo{Repo: Repo[model.File]{DB: tx, Model: "File"}}
}

func (f *FileRepo) getByUniqueKey(key string, value interface{}) (model.File, bool, string) {
	switch key {
	// 虽然 hash 并不是唯一的，但是并不影响功能
	case "id", "hash":
		value = value.(string)
	default:
		return model.File{}, false, "UnsupportedKey"
	}
	var file model.File
	res := f.DB.Model(&model.File{}).Where(key+" = ?", value).Limit(1).Find(&file)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get File: %s", res.Error)
		return model.File{}, false, "GetFileError"
	}
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

func (f *FileRepo) CountByKeyID(t string, key string, id uint) (int64, bool, string) {
	var count int64
	res := f.DB.Model(&model.File{}).Where("type = ? AND "+key+" = ?", t, id).Count(&count)
	if res.Error != nil {
		log.Logger.Warningf("Failed to count File: %s", res.Error)
		return 0, false, "CountModelError"
	}
	return count, true, "Success"
}

func (f *FileRepo) GetByKeyID(t string, key string, id uint, limit, offset int) ([]model.File, int64, bool, string) {
	var (
		files          = make([]model.File, 0)
		count, ok, msg = f.CountByKeyID(t, key, id)
	)
	if !ok {
		return files, count, false, msg
	}
	res := f.DB.Model(&model.File{}).Where("type = ? AND "+key+" = ?", t, id).Limit(limit).Offset(offset).Find(&files)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get File: %s", res.Error)
		return files, 0, false, "GetFileError"
	}
	return files, count, true, "Success"
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
	res := f.DB.Model(&model.File{}).Where("type = ?", t).Limit(limit).Offset(offset).Find(&files)
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
