package db

import (
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"gorm.io/gorm"
	"mime/multipart"
)

// RecordAvatar 添加头像记录
func RecordAvatar(tx *gorm.DB, path string, uploader uint, file *multipart.FileHeader, hash string) (model.Avatar, bool, string) {
	f := model.InitAvatar(path, uploader, file, hash)
	res := tx.Model(model.Avatar{}).Create(&f)
	if res.Error != nil {
		log.Logger.Warningf("Failed to record file: %v", res.Error)

		return model.Avatar{}, false, "CreateFileRecordError"
	}
	//go func() {
	//	if err := redis.DelFilesCache(); err != nil && !errors.Is(err, context.DeadlineExceeded) {
	//		log.Logger.Warningf("Failed to delete files cache: %v", err)
	//	}
	//}()
	return f, true, "Success"
}

// GetAvatarByID 以 ID 获取文件记录
func GetAvatarByID(tx *gorm.DB, id string) (model.Avatar, bool, string) {
	//if file, ok := redis.GetFileCache(id); ok {
	//	return file, true, "Success"
	//}
	var file model.Avatar
	res := tx.Model(model.Avatar{}).Where("id = ?", id).Find(&file).Limit(1)
	if res.RowsAffected != 1 {
		return model.Avatar{}, false, "FileNotFound"
	}
	//go func() {
	//	if err := redis.SetFileCache(file); err != nil && !errors.Is(err, context.DeadlineExceeded) {
	//		log.Logger.Errorf("Failed to delete file cache: %v", err)
	//	}
	//}()
	return file, true, "Success"
}

// GetAvatarByHash 以 Hash 获取文件记录
func GetAvatarByHash(tx *gorm.DB, hash string) (model.Avatar, bool, string) {
	var file model.Avatar
	res := tx.Model(model.Avatar{}).Where("hash = ?", hash).Find(&file).Limit(1)
	if res.RowsAffected != 1 {
		return model.Avatar{}, false, "FileNotFound"
	}
	return file, true, "Success"
}

// DeleteAvatar 以 ID 删除文件记录
func DeleteAvatar(tx *gorm.DB, id string) (bool, string) {
	if err := tx.Model(model.Avatar{}).Where("id = ?", id).Delete(&model.Avatar{}).Error; err != nil {
		log.Logger.Warningf("Failed to delete file: %v", id)
		return false, "DeleteFileError"
	}
	//go func() {
	//	if err := redis.DelFileCache(id); err != nil && !errors.Is(err, context.DeadlineExceeded) {
	//		log.Logger.Warningf("Failed to delete file cache: %v", err)
	//	}
	//	if err := redis.DelFilesCache(); err != nil && !errors.Is(err, context.DeadlineExceeded) {
	//		log.Logger.Warningf("Failed to delete files cache: %v", err)
	//	}
	//}()
	return true, "Success"
}

// GetAvatars 批量获取文件记录
func GetAvatars(tx *gorm.DB, limit int, offset int) ([]model.Avatar, int64, bool, string) {
	var files []model.Avatar
	var count int64
	res := tx.Model(&model.Avatar{})
	if res = res.Count(&count); res.Error != nil {
		log.Logger.Warningf("Failed to get files: %s", res.Error)
		return nil, 0, false, "UnknownError"
	}
	//if files, ok := redis.GetFilesCache(); ok {
	//	limit, offset = utils.TidyPaginate(len(files), limit, offset)
	//	return files[offset:limit], int64(len(files)), true, "Success"
	//}
	if res = res.Find(&files); res.Error != nil {
		log.Logger.Warningf("Failed to get files: %s", res.Error)
		return nil, 0, false, "FileNotFound"
	}
	//go func() {
	//	if err := redis.SetFilesCache(files); err != nil && !errors.Is(err, context.DeadlineExceeded) {
	//		log.Logger.Errorf("Failed to delete file cache: %v", err)
	//	}
	//}()
	limit, offset = utils.TidyPaginate(int(count), limit, offset)
	return files[offset:limit], count, true, "Success"
}
