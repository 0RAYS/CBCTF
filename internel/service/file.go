package service

import (
	"CBCTF/internel/config"
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"CBCTF/internel/utils"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"gorm.io/gorm"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
)

func SaveAvatar(tx *gorm.DB, uploaderID uint, file *multipart.FileHeader) (model.File, bool, string) {
	src, err := file.Open()
	if err != nil {
		log.Logger.Warningf("Failed to open file: %v", err)
		return model.File{}, false, "BadRequest"
	}
	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {
			log.Logger.Warningf("Failed to close file: %v", err)
		}
	}(src)
	sha256Sum := sha256.New()
	if _, err := io.Copy(sha256Sum, src); err != nil {
		log.Logger.Warningf("Failed to hash file: %v", err)
		return model.File{}, false, "UnknownError"
	}
	var (
		fileRepo      = db.InitFileRepo(tx)
		hash          = hex.EncodeToString(sha256Sum.Sum(nil))
		record, ok, _ = fileRepo.GetByHash(hash)
		path          string
		allowed       = []string{".png", ".jpg", ".jpeg"}
		suffix        = strings.ToLower(filepath.Ext(file.Filename))
	)
	if !utils.In(suffix, allowed) {
		return model.File{}, false, "FileNotAllowed"
	}
	if !ok {
		basePath := fmt.Sprintf("%s/avatars", config.Env.Path)
		path = fmt.Sprintf("%s/%s%s", basePath, utils.UUID(), suffix)
	} else {
		path = record.Path
	}
	return fileRepo.Create(db.CreateFileOptions{
		ID:       utils.UUID(),
		Filename: file.Filename,
		Size:     file.Size,
		Path:     path,
		Uploader: uploaderID,
		Suffix:   suffix,
		Hash:     hash,
		Type:     model.Avatar,
	})
}

func UpdateAvatar(tx *gorm.DB, v string, id uint, record model.File) (string, bool, string) {
	var (
		ok  bool
		msg string
	)
	path := fmt.Sprintf("/avatars/%s", record.ID)
	switch v {
	case "self-admin":
		ok, msg = db.InitAdminRepo(tx).Update(id, db.UpdateAdminOptions{Avatar: &path})
	case "self-user":
		ok, msg = db.InitUserRepo(tx).Update(id, db.UpdateUserOptions{Avatar: &path})
	case "user":
		ok, msg = db.InitUserRepo(tx).Update(id, db.UpdateUserOptions{Avatar: &path})
	case "contest":
		ok, msg = db.InitContestRepo(tx).Update(id, db.UpdateContestOptions{Avatar: &path})
	case "team":
		ok, msg = db.InitTeamRepo(tx).Update(id, db.UpdateTeamOptions{Avatar: &path})
	default:
		ok, msg = false, "UnsupportedKey"
	}
	return path, ok, msg
}

func SaveWriteUp(tx *gorm.DB, contestID, teamID uint, file *multipart.FileHeader) (model.File, bool, string) {
	src, err := file.Open()
	if err != nil {
		log.Logger.Warningf("Failed to open file: %v", err)
		return model.File{}, false, "BadRequest"
	}
	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {
			log.Logger.Warningf("Failed to close file: %v", err)
		}
	}(src)
	sha256Sum := sha256.New()
	if _, err := io.Copy(sha256Sum, src); err != nil {
		log.Logger.Warningf("Failed to hash file: %v", err)
		return model.File{}, false, "UnknownError"
	}
	var (
		fileRepo      = db.InitFileRepo(tx)
		hash          = hex.EncodeToString(sha256Sum.Sum(nil))
		record, ok, _ = fileRepo.GetByHash(hash)
		path          string
		allowed       = []string{".pdf", ".docx", ".doc"}
		suffix        = strings.ToLower(filepath.Ext(file.Filename))
	)
	if !utils.In(suffix, allowed) {
		return model.File{}, false, "FileNotAllowed"
	}
	if !ok {
		basePath := fmt.Sprintf("%s/writeups/%d/%d", config.Env.Path, contestID, teamID)
		path = fmt.Sprintf("%s/%s%s", basePath, utils.UUID(), suffix)
	} else {
		path = record.Path
	}
	return fileRepo.Create(db.CreateFileOptions{
		ID:       utils.UUID(),
		Filename: file.Filename,
		Size:     file.Size,
		Path:     path,
		Uploader: teamID,
		Suffix:   suffix,
		Hash:     hash,
		Type:     model.WriteUP,
	})
}
