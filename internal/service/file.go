package service

import (
	"CBCTF/internal/config"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	db "CBCTF/internal/repo"
	"CBCTF/internal/utils"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"gorm.io/gorm"
	"io"
	"mime/multipart"
	"path/filepath"
	"slices"
	"strings"
)

func SaveAvatar(tx *gorm.DB, options db.CreateFileOptions, file *multipart.FileHeader) (model.File, bool, string) {
	src, err := file.Open()
	if err != nil {
		log.Logger.Warningf("Failed to open file: %v", err)
		return model.File{}, false, i18n.BadRequest
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
		return model.File{}, false, i18n.UnknownError
	}
	var (
		fileRepo      = db.InitFileRepo(tx)
		hash          = hex.EncodeToString(sha256Sum.Sum(nil))
		record, ok, _ = fileRepo.GetByHash(hash)
		path          string
		allowed       = []string{".png", ".jpg", ".jpeg"}
		suffix        = strings.ToLower(filepath.Ext(file.Filename))
	)
	if !slices.Contains(allowed, suffix) {
		return model.File{}, false, i18n.FileNotAllowed
	}
	if !ok {
		path = fmt.Sprintf("%s/avatars/%s%s", config.Env.Path, utils.UUID(), suffix)
	} else {
		path = record.Path
	}
	options.RandID = utils.UUID()
	options.Filename = file.Filename
	options.Size = file.Size
	options.Path = path
	options.Suffix = suffix
	options.Hash = hash
	options.Type = model.AvatarFile
	return fileRepo.Create(options)
}

func UpdateAvatar(tx *gorm.DB, v string, id uint, record model.File) (string, bool, string) {
	var (
		ok  bool
		msg string
	)
	path := model.AvatarURL(fmt.Sprintf("/avatars/%s", record.RandID))
	switch v {
	case "admin":
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
		ok, msg = false, i18n.UnsupportedKey
	}
	return string(path), ok, msg
}

func SaveWriteUp(tx *gorm.DB, user model.User, contest model.Contest, team model.Team, file *multipart.FileHeader) (model.File, bool, string) {
	src, err := file.Open()
	if err != nil {
		log.Logger.Warningf("Failed to open file: %v", err)
		return model.File{}, false, i18n.BadRequest
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
		return model.File{}, false, i18n.UnknownError
	}
	var (
		fileRepo      = db.InitFileRepo(tx)
		hash          = hex.EncodeToString(sha256Sum.Sum(nil))
		record, ok, _ = fileRepo.GetByHash(hash)
		path          string
		allowed       = []string{".pdf", ".docx", ".doc"}
		suffix        = strings.ToLower(filepath.Ext(file.Filename))
	)
	if !slices.Contains(allowed, suffix) {
		return model.File{}, false, i18n.FileNotAllowed
	}
	if !ok {
		path = fmt.Sprintf("%s/writeups/contest-%d/team-%d", config.Env.Path, contest.ID, team.ID)
		path += fmt.Sprintf("/%s%s", utils.UUID(), suffix)
	} else {
		path = record.Path
	}
	return fileRepo.Create(db.CreateFileOptions{
		RandID:    utils.UUID(),
		Filename:  file.Filename,
		Size:      file.Size,
		Path:      path,
		UserID:    &user.ID,
		TeamID:    &team.ID,
		ContestID: &contest.ID,
		Suffix:    suffix,
		Hash:      hash,
		Type:      model.WriteUPFile,
	})
}
