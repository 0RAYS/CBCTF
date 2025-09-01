package service

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/prometheus"
	"CBCTF/internal/utils"
	"database/sql"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"slices"
	"strings"

	"gorm.io/gorm"
)

func SaveAvatar(tx *gorm.DB, options db.CreateFileOptions, file *multipart.FileHeader) (model.File, bool, string) {
	var (
		fileRepo = db.InitFileRepo(tx)
		allowed  = []string{".png", ".jpg", ".jpeg"}
		suffix   = strings.ToLower(filepath.Ext(file.Filename))
	)
	if !slices.Contains(allowed, suffix) {
		return model.File{}, false, i18n.FileNotAllowed
	}
	size, hash, err := utils.GetFileInfoByHeader(file)
	if err != nil {
		log.Logger.Warningf("Failed to get file info: %s", err)
		return model.File{}, false, i18n.UnknownError
	}
	options.RandID = utils.UUID()
	options.Filename = file.Filename
	options.Size = size
	options.Path = fmt.Sprintf("%s/avatars/%s%s", config.Env.Path, utils.UUID(), suffix)
	options.Suffix = suffix
	options.Hash = hash
	options.Type = model.AvatarFileType
	record, ok, msg := fileRepo.Create(options)
	if ok {
		prometheus.UpdateFileUploadMetrics(record.Suffix, record.Size)
	}
	return record, ok, msg
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
	case "oath":
		ok, msg = db.InitOauthRepo(tx).Update(id, db.UpdateOauthOptions{Avatar: &path})
	default:
		ok, msg = false, i18n.UnsupportedKey
	}
	return string(path), ok, msg
}

func SaveChallengeFile(tx *gorm.DB, challenge model.Challenge, file *multipart.FileHeader, path string) (model.File, bool, string) {
	var (
		fileRepo = db.InitFileRepo(tx)
		suffix   = strings.ToLower(filepath.Ext(file.Filename))
	)
	size, hash, err := utils.GetFileInfoByHeader(file)
	if err != nil {
		log.Logger.Warningf("Failed to get file info: %s", err)
		return model.File{}, false, i18n.UnknownError
	}
	record, ok, msg := fileRepo.Create(db.CreateFileOptions{
		RandID:      utils.UUID(),
		Filename:    file.Filename,
		Size:        size,
		Path:        path,
		ChallengeID: sql.Null[uint]{V: challenge.ID, Valid: true},
		Suffix:      suffix,
		Hash:        hash,
		Type:        model.ChallengeFileType,
	})
	if ok {
		prometheus.UpdateFileUploadMetrics(record.Suffix, file.Size)
	}
	return record, ok, msg
}

func SaveWriteUp(tx *gorm.DB, user model.User, contest model.Contest, team model.Team, file *multipart.FileHeader) (model.File, bool, string) {
	var (
		allowed = []string{".pdf", ".docx", ".doc"}
		suffix  = strings.ToLower(filepath.Ext(file.Filename))
	)
	if !slices.Contains(allowed, suffix) {
		return model.File{}, false, i18n.FileNotAllowed
	}
	size, hash, err := utils.GetFileInfoByHeader(file)
	if err != nil {
		log.Logger.Warningf("Failed to get file info: %s", err)
		return model.File{}, false, i18n.UnknownError
	}
	record, ok, msg := db.InitFileRepo(tx).Create(db.CreateFileOptions{
		RandID:    utils.UUID(),
		Filename:  file.Filename,
		Size:      size,
		Path:      fmt.Sprintf("%s/writeups/contest-%d/team-%d/%s%s", config.Env.Path, contest.ID, team.ID, utils.UUID(), suffix),
		UserID:    sql.Null[uint]{V: user.ID, Valid: true},
		TeamID:    sql.Null[uint]{V: team.ID, Valid: true},
		ContestID: sql.Null[uint]{V: contest.ID, Valid: true},
		Suffix:    suffix,
		Hash:      hash,
		Type:      model.WriteUPFileType,
	})
	if ok {
		prometheus.UpdateFileUploadMetrics(record.Suffix, file.Size)
	}
	return record, ok, msg
}
