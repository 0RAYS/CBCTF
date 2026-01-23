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

func SaveAvatar(tx *gorm.DB, options db.CreateFileOptions, file *multipart.FileHeader) (model.File, model.RetVal) {
	var (
		fileRepo = db.InitFileRepo(tx)
		allowed  = []string{".png", ".jpg", ".jpeg"}
		suffix   = strings.ToLower(filepath.Ext(file.Filename))
	)
	if !slices.Contains(allowed, suffix) {
		return model.File{}, model.RetVal{Msg: i18n.Model.File.NotAllowed}
	}
	size, hash, err := utils.GetFileInfoByHeader(file)
	if err != nil {
		log.Logger.Warningf("Failed to get file info: %s", err)
		return model.File{}, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err}}
	}
	options.RandID = utils.UUID()
	options.Filename = file.Filename
	options.Size = size
	options.Path = fmt.Sprintf("%s/avatars/%s%s", config.Env.Path, utils.UUID(), suffix)
	options.Suffix = suffix
	options.Hash = hash
	options.Type = model.AvatarFileType
	record, ret := fileRepo.Create(options)
	if ret.OK {
		prometheus.UpdateFileUploadMetrics(record.Suffix, record.Size)
	}
	return record, ret
}

func UpdateAvatar(tx *gorm.DB, v string, id uint, record model.File) (string, model.RetVal) {
	var ret model.RetVal
	path := model.AvatarURL(fmt.Sprintf("/avatars/%s", record.RandID))
	switch v {
	case "admin":
		ret = db.InitAdminRepo(tx).Update(id, db.UpdateAdminOptions{Avatar: &path})
	case "self-user":
		ret = db.InitUserRepo(tx).Update(id, db.UpdateUserOptions{Avatar: &path})
	case "user":
		ret = db.InitUserRepo(tx).Update(id, db.UpdateUserOptions{Avatar: &path})
	case "contest":
		ret = db.InitContestRepo(tx).Update(id, db.UpdateContestOptions{Avatar: &path})
	case "team":
		ret = db.InitTeamRepo(tx).Update(id, db.UpdateTeamOptions{Avatar: &path})
	case "oauth":
		ret = db.InitOauthRepo(tx).Update(id, db.UpdateOauthOptions{Avatar: &path})
	default:
		ret = model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": "Invalid Avatar"}}
	}
	return string(path), ret
}

func SaveChallengeFile(tx *gorm.DB, challenge model.Challenge, file *multipart.FileHeader, path string) (model.File, model.RetVal) {
	var (
		fileRepo = db.InitFileRepo(tx)
		suffix   = strings.ToLower(filepath.Ext(file.Filename))
	)
	size, hash, err := utils.GetFileInfoByHeader(file)
	if err != nil {
		log.Logger.Warningf("Failed to get file info: %s", err)
		return model.File{}, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err}}
	}
	record, ret := fileRepo.Get(db.GetOptions{
		Conditions: map[string]any{"challenge_id": challenge.ID, "type": model.ChallengeFileType},
	})
	if ret.OK {
		if hash == record.Hash {
			return record, model.SuccessRetVal()
		}
		if ret = fileRepo.Delete(record.ID); !ret.OK {
			return model.File{}, ret
		}
	}
	record, ret = fileRepo.Create(db.CreateFileOptions{
		RandID:      utils.UUID(),
		Filename:    file.Filename,
		Size:        size,
		Path:        path,
		ChallengeID: sql.Null[uint]{V: challenge.ID, Valid: true},
		Suffix:      suffix,
		Hash:        hash,
		Type:        model.ChallengeFileType,
	})
	if ret.OK {
		prometheus.UpdateFileUploadMetrics(record.Suffix, file.Size)
	}
	return record, ret
}

func SaveWriteUp(tx *gorm.DB, user model.User, contest model.Contest, team model.Team, file *multipart.FileHeader) (model.File, model.RetVal) {
	var (
		allowed = []string{".pdf", ".docx", ".doc"}
		suffix  = strings.ToLower(filepath.Ext(file.Filename))
	)
	if !slices.Contains(allowed, suffix) {
		return model.File{}, model.RetVal{Msg: i18n.Model.File.NotAllowed}
	}
	size, hash, err := utils.GetFileInfoByHeader(file)
	if err != nil {
		log.Logger.Warningf("Failed to get file info: %s", err)
		return model.File{}, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err}}
	}
	record, ret := db.InitFileRepo(tx).Create(db.CreateFileOptions{
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
	if ret.OK {
		prometheus.UpdateFileUploadMetrics(record.Suffix, file.Size)
	}
	return record, ret
}
