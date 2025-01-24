package db

import (
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
	"context"
	"errors"
	"fmt"
	"mime/multipart"
)

// RecordFile 添加文件记录
func RecordFile(ctx context.Context, path string, uploader uint, file *multipart.FileHeader, hash string) (model.Avatar, bool, string) {
	f := model.InitFile(path, uploader, file, hash)
	res := DB.WithContext(ctx).Model(model.Avatar{}).Create(&f)
	if res.Error != nil {
		log.Logger.Warningf("Failed to record file: %v", res.Error)
		return model.Avatar{}, false, "CreateFileRecordError"
	}
	go func() {
		if err := redis.DelFilesCache(); err != nil && !errors.Is(err, context.DeadlineExceeded) {
			log.Logger.Warningf("Failed to delete files cache: %v", err)
		}
	}()
	return f, true, "Success"
}

// GetFileByID 以 ID 获取文件记录
func GetFileByID(ctx context.Context, id string) (model.Avatar, bool, string) {
	cacheKey := fmt.Sprintf("file:%s", id)
	if file, ok := redis.GetFileCache(cacheKey); ok {
		return file, true, "Success"
	}
	var file model.Avatar
	res := DB.WithContext(ctx).Model(model.Avatar{}).Where("id = ?", id).Find(&file).Limit(1)
	if res.RowsAffected != 1 {
		return model.Avatar{}, false, "FileNotFound"
	}
	go func() {
		if err := redis.SetFileCache(cacheKey, file); err != nil && !errors.Is(err, context.DeadlineExceeded) {
			log.Logger.Errorf("Failed to delete file cache: %v", err)
		}
	}()
	return file, true, "Success"
}

func GetFileByHash(ctx context.Context, hash string) (model.Avatar, bool, string) {
	cacheKey := fmt.Sprintf("file:hash:%s", hash)
	if file, ok := redis.GetFileCache(cacheKey); ok {
		return file, true, "Success"
	}
	var file model.Avatar
	res := DB.WithContext(ctx).Model(model.Avatar{}).Where("hash = ?", hash).Find(&file).Limit(1)
	if res.RowsAffected != 1 {
		return model.Avatar{}, false, "FileNotFound"
	}
	go func() {
		if err := redis.SetFileCache(cacheKey, file); err != nil && !errors.Is(err, context.DeadlineExceeded) {
			log.Logger.Errorf("Failed to delete file cache: %v", err)
		}
	}()
	return file, true, "Success"
}

// DeleteFile 以 ID 删除文件记录
func DeleteFile(ctx context.Context, id string) (bool, string) {
	if err := DB.WithContext(ctx).Model(model.Avatar{}).Where("id = ?", id).Delete(&model.Avatar{}).Error; err != nil {
		log.Logger.Warningf("Failed to delete file: %v", id)
		return false, "DeleteFileError"
	}
	go func() {
		if err := redis.DelFileCache(id); err != nil && !errors.Is(err, context.DeadlineExceeded) {
			log.Logger.Warningf("Failed to delete file cache: %v", err)
		}
		if err := redis.DelFilesCache(); err != nil && !errors.Is(err, context.DeadlineExceeded) {
			log.Logger.Warningf("Failed to delete files cache: %v", err)
		}
	}()
	return true, "Success"
}

// GetFiles 批量获取文件记录
func GetFiles(ctx context.Context, limit int, offset int) ([]model.Avatar, int64, bool, string) {
	if limit <= 0 {
		limit = -1
	}
	if offset <= 0 {
		offset = -1
	}
	var files []model.Avatar
	var count int64
	res := DB.WithContext(ctx).Model(&model.Avatar{})
	if res = res.Count(&count); res.Error != nil {
		log.Logger.Warningf("Failed to get files: %s", res.Error.Error())
		return nil, 0, false, "UnknownError"
	}
	cacheKey := fmt.Sprintf("files:%d:%d", limit, offset)
	if files, ok := redis.GetFilesCache(cacheKey); ok {
		return files, int64(len(files)), true, "Success"
	}
	if res = res.Limit(limit).Offset(offset).Find(&files); res.Error != nil {
		log.Logger.Warningf("Failed to get files: %s", res.Error.Error())
		return nil, 0, false, "FileNotFound"
	}
	go func() {
		if err := redis.SetFilesCache(cacheKey, files); err != nil && !errors.Is(err, context.DeadlineExceeded) {
			log.Logger.Errorf("Failed to delete file cache: %v", err)
		}
	}()
	return files, count, true, "Success"
}
