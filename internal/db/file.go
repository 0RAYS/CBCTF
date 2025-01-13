package db

import (
	"RayWar/internal/log"
	"RayWar/internal/model"
	"mime/multipart"
)

// RecordFile 添加文件记录
func RecordFile(ownerID uint, random string, path string, fileHeader *multipart.FileHeader) (model.File, bool, string) {
	file := model.InitFile(ownerID, random, path, fileHeader)
	res := DB.Model(model.File{}).Create(&file)
	if res.Error != nil {
		log.Logger.Warningf("Failed to record file: %v", res.Error)
		return model.File{}, false, "CreateFileRecordError"
	}
	return file, true, "Success"
}

// GetFile 以 ID 获取文件记录
func GetFile(fileID string) (model.File, bool, string) {
	var file model.File
	res := DB.Model(model.File{}).Where("id = ?", fileID).Find(&file).Limit(1)
	if res.RowsAffected != 1 {
		return model.File{}, false, "FileNotFound"
	}
	return file, true, "Success"
}

// DeleteFile 以 ID 删除文件记录
func DeleteFile(fileID string) (bool, string) {
	if err := DB.Model(model.File{}).Where("id = ?", fileID).Delete(&model.File{}).Error; err != nil {
		log.Logger.Warningf("Failed to delete file: %v", fileID)
		return false, "DeleteFileError"
	}
	return true, "Success"
}

// GetFiles 批量获取文件记录
func GetFiles(limit int, offset int) ([]model.File, int64, bool, string) {
	if limit <= 0 {
		limit = -1
	}
	if offset <= 0 {
		offset = -1
	}
	var files []model.File
	var total int64
	if res := DB.Model(&model.File{}).Count(&total); res.Error != nil {
		log.Logger.Warningf("Failed to get files: %s", res.Error.Error())
		return nil, 0, false, "UnknownError"
	}
	if res := DB.Model(&model.File{}).Limit(limit).Offset(offset).Find(&files); res.Error != nil {
		log.Logger.Warningf("Failed to get files: %s", res.Error.Error())
		return nil, 0, false, "FileNotFound"
	}
	return files, total, true, "Success"
}
