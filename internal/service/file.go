package service

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/prometheus"
	"CBCTF/internal/utils"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"slices"
	"strings"

	"gorm.io/gorm"
)

func SaveAvatar(tx *gorm.DB, options db.CreateFileOptions, file *multipart.FileHeader) (model.File, bool, string) {
	src, err := file.Open()
	if err != nil {
		log.Logger.Warningf("Failed to open file: %s", err)
		return model.File{}, false, i18n.BadRequest
	}
	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {
			log.Logger.Warningf("Failed to close file: %s", err)
		}
	}(src)
	sha256Sum := sha256.New()
	if _, err := io.Copy(sha256Sum, src); err != nil {
		log.Logger.Warningf("Failed to hash file: %s", err)
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
	f, ok, msg := fileRepo.Create(options)
	if ok {
		prometheus.UpdateFileUploadMetrics(record.Suffix, record.Size)
	}
	return f, ok, msg
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
	src, err := file.Open()
	if err != nil {
		log.Logger.Warningf("Failed to open file: %s", err)
		return model.File{}, false, i18n.BadRequest
	}
	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {
			log.Logger.Warningf("Failed to close file: %s", err)
		}
	}(src)
	sha256Sum := sha256.New()
	if _, err := io.Copy(sha256Sum, src); err != nil {
		log.Logger.Warningf("Failed to hash file: %s", err)
		return model.File{}, false, i18n.UnknownError
	}
	var (
		fileRepo = db.InitFileRepo(tx)
		hash     = hex.EncodeToString(sha256Sum.Sum(nil))
		suffix   = strings.ToLower(filepath.Ext(file.Filename))
	)
	records, _, ok, msg := fileRepo.List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"path": path}, Selects: []string{"id"},
	})
	if ok {
		idL := make([]uint, 0)
		for _, record := range records {
			idL = append(idL, record.ID)
		}
		fileRepo.Delete(idL...)
	}
	record, ok, msg := fileRepo.Create(db.CreateFileOptions{
		RandID:      utils.UUID(),
		Filename:    file.Filename,
		Size:        file.Size,
		Path:        path,
		ChallengeID: sql.Null[uint]{V: challenge.ID, Valid: true},
		Suffix:      suffix,
		Hash:        hash,
		Type:        model.ChallengeFile,
	})
	if ok {
		prometheus.UpdateFileUploadMetrics(record.Suffix, file.Size)
	}
	return record, ok, msg
}

func SaveWriteUp(tx *gorm.DB, user model.User, contest model.Contest, team model.Team, file *multipart.FileHeader) (model.File, bool, string) {
	src, err := file.Open()
	if err != nil {
		log.Logger.Warningf("Failed to open file: %s", err)
		return model.File{}, false, i18n.BadRequest
	}
	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {
			log.Logger.Warningf("Failed to close file: %s", err)
		}
	}(src)
	sha256Sum := sha256.New()
	if _, err := io.Copy(sha256Sum, src); err != nil {
		log.Logger.Warningf("Failed to hash file: %s", err)
		return model.File{}, false, i18n.UnknownError
	}
	var (
		fileRepo      = db.InitFileRepo(tx)
		hash          = hex.EncodeToString(sha256Sum.Sum(nil))
		record, ok, _ = fileRepo.GetByHash(hash)
		path          string
		allowed       = []string{".pdf", ".docx", ".doc"}
		suffix        = strings.ToLower(filepath.Ext(file.Filename))
		msg           string
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
	record, ok, msg = fileRepo.Create(db.CreateFileOptions{
		RandID:    utils.UUID(),
		Filename:  file.Filename,
		Size:      file.Size,
		Path:      path,
		UserID:    sql.Null[uint]{V: user.ID, Valid: true},
		TeamID:    sql.Null[uint]{V: team.ID, Valid: true},
		ContestID: sql.Null[uint]{V: contest.ID, Valid: true},
		Suffix:    suffix,
		Hash:      hash,
		Type:      model.WriteUPFile,
	})
	if ok {
		prometheus.UpdateFileUploadMetrics(record.Suffix, file.Size)
	}
	return record, ok, msg
}
