package db

import (
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"gorm.io/gorm"
	"mime/multipart"
)

// RecordFile 添加头像记录
func RecordFile(tx *gorm.DB, path string, uploader uint, file *multipart.FileHeader, hash string, t string) (model.File, bool, string) {
	f := model.InitFile(path, uploader, file, hash, t)
	res := tx.Model(&model.File{}).Create(&f)
	if res.Error != nil {
		log.Logger.Warningf("Failed to record file: %v", res.Error)

		return model.File{}, false, "CreateFileRecordError"
	}
	return f, true, "Success"
}

// GetFileByID 以 ID 获取文件记录
func GetFileByID(tx *gorm.DB, id string) (model.File, bool, string) {
	var file model.File
	res := tx.Model(&model.File{}).Where("id = ?", id).Find(&file).Limit(1)
	if res.RowsAffected != 1 {
		return model.File{}, false, "FileNotFound"
	}
	return file, true, "Success"
}

// GetFileByHash 以 Hash 获取文件记录
func GetFileByHash(tx *gorm.DB, hash string) (model.File, bool, string) {
	var file model.File
	res := tx.Model(&model.File{}).Where("hash = ?", hash).Find(&file).Limit(1)
	if res.RowsAffected != 1 {
		return model.File{}, false, "FileNotFound"
	}
	return file, true, "Success"
}

// DeleteFile 以 ID 删除文件记录
func DeleteFile(tx *gorm.DB, id string) (bool, string) {
	if err := tx.Model(&model.File{}).Where("id = ?", id).Delete(&model.File{}).Error; err != nil {
		log.Logger.Warningf("Failed to delete file: %v", id)
		return false, "DeleteFileError"
	}
	return true, "Success"
}

// GetAvatars 批量获取文件记录
func GetAvatars(tx *gorm.DB, limit int, offset int) ([]model.File, int64, bool, string) {
	var files []model.File
	var count int64
	res := tx.Model(&model.File{}).Where("type = ?", model.Avatar)
	if res = res.Count(&count); res.Error != nil {
		log.Logger.Warningf("Failed to get files: %s", res.Error)
		return make([]model.File, 0), 0, false, "UnknownError"
	}
	limit, offset = utils.TidyPaginate(int(count), limit, offset)
	if res = res.Limit(limit).Offset(offset).Find(&files); res.Error != nil {
		log.Logger.Warningf("Failed to get files: %s", res.Error)
		return make([]model.File, 0), 0, false, "FileNotFound"
	}
	return files, count, true, "Success"
}
