package service

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/prometheus"
	"CBCTF/internal/utils"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

func SavePicture(tx *gorm.DB, modelName string, modelID uint, file *multipart.FileHeader) (model.File, model.RetVal) {
	var (
		fileRepo = db.InitFileRepo(tx)
		allowed  = []string{".png", ".jpg", ".jpeg", ".gif"}
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
	options := db.CreateFileOptions{
		RandID:   utils.UUID(),
		Filename: file.Filename,
		Size:     size,
		Path:     model.FilePath(filepath.Join(config.Env.Path, "pictures", utils.UUID()+suffix)),
		Model:    modelName,
		ModelID:  modelID,
		Suffix:   suffix,
		Hash:     hash,
		Type:     model.PictureFileType,
	}
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
	case "branding":
		ret = db.InitBrandingRepo(tx).Update(id, db.UpdateBrandingOptions{HomeLogo: &path})
	default:
		ret = model.RetVal{Msg: i18n.Model.File.InvalidPicture}
	}
	return string(path), ret
}

func SaveChallenge(tx *gorm.DB, challenge model.Challenge, file *multipart.FileHeader, path string) (model.File, model.RetVal) {
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
		Conditions: map[string]any{"model": model.Name(challenge), "model_id": challenge.ID, "type": model.ChallengeFileType},
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
		Path:     model.FilePath(path),
		Model:    model.Name(challenge),
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

func SaveWriteup(tx *gorm.DB, contest model.Contest, team model.Team, file *multipart.FileHeader) (model.File, model.RetVal) {
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
		Path: model.FilePath(filepath.Join(
			config.Env.Path,
			"writeups",
			"contest-"+strconv.FormatUint(uint64(contest.ID), 10),
			"team-"+strconv.FormatUint(uint64(team.ID), 10),
			time.Now().Format("2006-01-02 15:04:05")+suffix,
		)),
		Model:   model.Name(team),
		ModelID: team.ID,
		Suffix:  suffix,
		Hash:    hash,
		Type:    model.WriteupFileType,
	})
	if ret.OK {
		prometheus.RecordFileUpload(record.Suffix, file.Size)
	}
	return record, ret
}

func ListFiles(tx *gorm.DB, form dto.GetFilesForm) ([]model.File, int64, model.RetVal) {
	options := db.GetOptions{Sort: []string{"id DESC"}}
	if form.Type != "" {
		options.Conditions = map[string]any{"type": form.Type}
	}
	return db.InitFileRepo(tx).List(form.Limit, form.Offset, options)
}

func ListWriteups(tx *gorm.DB, team model.Team, form dto.ListModelsForm) ([]model.File, int64, model.RetVal) {
	return db.InitFileRepo(tx).List(form.Limit, form.Offset, db.GetOptions{
		Conditions: map[string]any{"model": model.Name(team), "model_id": team.ID, "type": model.WriteupFileType},
		Sort:       []string{"id DESC"},
	})
}

func DeleteFiles(tx *gorm.DB, form dto.DeleteFileForm) model.RetVal {
	return db.InitFileRepo(tx).DeleteByRandID(form.FileIDs...)
}
