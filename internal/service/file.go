package service

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/prometheus"
	"CBCTF/internal/utils"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"slices"
	"strings"

	"gorm.io/gorm"
)

func SavePicture(tx *gorm.DB, options db.CreateFileOptions, file *multipart.FileHeader) (model.File, model.RetVal) {
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
		return model.File{}, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}}
	}
	options.RandID = utils.UUID()
	options.Filename = file.Filename
	options.Size = size
	options.Path = fmt.Sprintf("%s/pictures/%s%s", config.Env.Path, utils.UUID(), suffix)
	options.Suffix = suffix
	options.Hash = hash
	options.Type = model.PictureFileType
	record, ret := fileRepo.Create(options)
	if ret.OK {
		prometheus.RecordFileUpload(record.Suffix, record.Size)
	}
	return record, ret
}

func UpdatePicture(tx *gorm.DB, v string, id uint, record model.File) (string, model.RetVal) {
	var ret model.RetVal
	path := model.FileURL(fmt.Sprintf("/pictures/%s", record.RandID))
	switch v {
	case "self":
		ret = db.InitUserRepo(tx).Update(id, db.UpdateUserOptions{Picture: &path})
	case "user":
		ret = db.InitUserRepo(tx).Update(id, db.UpdateUserOptions{Picture: &path})
	case "contest":
		ret = db.InitContestRepo(tx).Update(id, db.UpdateContestOptions{Picture: &path})
	case "team":
		ret = db.InitTeamRepo(tx).Update(id, db.UpdateTeamOptions{Picture: &path})
	case "oauth":
		ret = db.InitOauthRepo(tx).Update(id, db.UpdateOauthOptions{Picture: &path})
	default:
		ret = model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": "Invalid Picture"}}
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
		return model.File{}, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}}
	}
	record, ret := fileRepo.Get(db.GetOptions{
		Conditions: map[string]any{"model": challenge.ModelName(), "model_id": challenge.ID, "type": model.ChallengeFileType},
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
		RandID:   utils.UUID(),
		Filename: file.Filename,
		Size:     size,
		Path:     path,
		Model:    challenge.ModelName(),
		ModelID:  challenge.ID,
		Suffix:   suffix,
		Hash:     hash,
		Type:     model.ChallengeFileType,
	})
	if ret.OK {
		prometheus.RecordFileUpload(record.Suffix, file.Size)
	}
	return record, ret
}

func SaveWriteUp(tx *gorm.DB, contest model.Contest, team model.Team, file *multipart.FileHeader) (model.File, model.RetVal) {
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
		return model.File{}, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}}
	}
	record, ret := db.InitFileRepo(tx).Create(db.CreateFileOptions{
		RandID:   utils.UUID(),
		Filename: file.Filename,
		Size:     size,
		Path:     fmt.Sprintf("%s/writeups/contest-%d/team-%d/%s%s", config.Env.Path, contest.ID, team.ID, utils.UUID(), suffix),
		Model:    team.ModelName(),
		ModelID:  team.ID,
		Suffix:   suffix,
		Hash:     hash,
		Type:     model.WriteupFileType,
	})
	if ret.OK {
		prometheus.RecordFileUpload(record.Suffix, file.Size)
	}
	return record, ret
}
