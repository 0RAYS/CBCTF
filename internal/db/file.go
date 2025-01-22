package db

import (
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"mime/multipart"
)

// RecordFile 添加文件记录
func RecordFile(ctx context.Context, path string, uploader uint, file *multipart.FileHeader, hash string, isAdmin bool, isAttachment bool) (model.File, bool, string) {
	f := model.InitFile(path, uploader, file, hash, isAdmin, isAttachment)
	res := DB.WithContext(ctx).Model(model.File{}).Create(&f)
	if res.Error != nil {
		log.Logger.Warningf("Failed to record file: %v", res.Error)
		return model.File{}, false, "CreateFileRecordError"
	}
	return f, true, "Success"
}

// GetFileByID 以 ID 获取文件记录
func GetFileByID(ctx context.Context, id string) (model.File, bool, string) {
	var file model.File
	res := DB.WithContext(ctx).Model(model.File{}).Where("id = ?", id).Find(&file).Limit(1)
	if res.RowsAffected != 1 {
		return model.File{}, false, "FileNotFound"
	}
	return file, true, "Success"
}

func GetFileByHash(ctx context.Context, hash string) (model.File, bool, string) {
	var file model.File
	res := DB.WithContext(ctx).Model(model.File{}).Where("hash = ?", hash).Find(&file).Limit(1)
	if res.RowsAffected != 1 {
		return model.File{}, false, "FileNotFound"
	}
	return file, true, "Success"
}

// DeleteFile 以 ID 删除文件记录
func DeleteFile(ctx context.Context, id string) (bool, string) {
	if err := DB.WithContext(ctx).Model(model.File{}).Where("id = ?", id).Delete(&model.File{}).Error; err != nil {
		log.Logger.Warningf("Failed to delete file: %v", id)
		return false, "DeleteFileError"
	}
	return true, "Success"
}

// GetFiles 批量获取文件记录
func GetFiles(ctx context.Context, limit int, offset int) ([]model.File, int64, bool, string) {
	if limit <= 0 {
		limit = -1
	}
	if offset <= 0 {
		offset = -1
	}
	var files []model.File
	var count int64
	res := DB.WithContext(ctx).Model(&model.File{})
	if res = res.Count(&count); res.Error != nil {
		log.Logger.Warningf("Failed to get files: %s", res.Error.Error())
		return nil, 0, false, "UnknownError"
	}
	if res = res.Limit(limit).Offset(offset).Find(&files); res.Error != nil {
		log.Logger.Warningf("Failed to get files: %s", res.Error.Error())
		return nil, 0, false, "FileNotFound"
	}
	return files, count, true, "Success"
}
