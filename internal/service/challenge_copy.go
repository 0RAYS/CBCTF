package service

import (
	"database/sql"
	"io"
	"os"
	"path/filepath"
	"time"

	"gorm.io/gorm"

	"CBCTF/internal/db"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
)

// sqlNullContestID 把比赛 id 转成可空的归属，0 表示不归属任何比赛。
func sqlNullContestID(contestID uint) sql.Null[uint] {
	return sql.Null[uint]{V: contestID, Valid: contestID > 0}
}

// CopyChallenge 把题目复制成一个新的独立题目。
// 副本包含展示字段、靶机规格、生成器镜像、网络策略、flag 模板与注入绑定以及附件文件，
// 与源题目没有任何关联：解题记录、得分与计分状态都不复制。
func CopyChallenge(tx *gorm.DB, source model.Challenge, contestID uint) (model.Challenge, model.RetVal) {
	challenge, ret := db.InitChallengeRepo(tx).Create(model.Challenge{
		RandID:          utils.UUID(),
		ContestID:       sqlNullContestID(contestID),
		Name:            source.Name,
		Description:     source.Description,
		Type:            source.Type,
		Category:        source.Category,
		Hidden:          contestID != 0,
		Attempt:         source.Attempt,
		Hints:           append([]string(nil), source.Hints...),
		Tags:            append([]string(nil), source.Tags...),
		Template:        source.Template.DeepCopy(),
		GeneratorImage:  source.GeneratorImage,
		NetworkPolicies: source.NetworkPolicies.DeepCopy(),
	})
	if !ret.OK {
		return model.Challenge{}, ret
	}
	flagRepo := db.InitChallengeFlagRepo(tx)
	for _, flag := range source.ChallengeFlags {
		if _, ret = flagRepo.Create(model.ChallengeFlag{
			ChallengeID:  challenge.ID,
			Value:        flag.Value,
			Binding:      flag.Binding,
			Generation:   1,
			Active:       true,
			Score:        1000,
			CurrentScore: 1000,
			Decay:        100,
			MinScore:     100,
			ScoreType:    model.StaticScoreType,
			Last:         time.Now(),
		}); !ret.OK {
			return model.Challenge{}, ret
		}
	}
	if ret = copyChallengeFiles(tx, source, challenge); !ret.OK {
		return model.Challenge{}, ret
	}
	return challenge, model.SuccessRetVal()
}

// copyChallengeFiles 把题目的附件 / 生成器复制到新题目目录，并登记 files 记录。
// 源文件缺失时只记录告警，不阻断复制。
func copyChallengeFiles(tx *gorm.DB, source model.Challenge, target model.Challenge) model.RetVal {
	var sourcePath, targetPath string
	switch source.Type {
	case model.DynamicChallengeType:
		sourcePath, targetPath = source.GeneratorPath(), target.GeneratorPath()
	default:
		sourcePath, targetPath = source.StaticPath(), target.StaticPath()
	}
	sourceRecord, sourceRet := sourceFileRecord(tx, source)
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		log.Logger.Warningf("Failed to create challenge dir: challenge_id=%d error=%v", target.ID, err)
		return model.SuccessRetVal()
	}
	if err := copyFileContent(sourcePath, targetPath); err != nil {
		if !os.IsNotExist(err) {
			log.Logger.Warningf("Failed to copy challenge attachment: challenge_id=%d target_id=%d error=%v", source.ID, target.ID, err)
		}
		return model.SuccessRetVal()
	}
	record := model.File{
		RandID:   utils.UUID(),
		Filename: filepath.Base(targetPath),
		Size:     0,
		Path:     model.FilePath(targetPath),
		Model:    model.Name(target),
		ModelID:  target.ID,
		Suffix:   filepath.Ext(targetPath),
		Type:     model.ChallengeFileType,
	}
	if sourceRet.OK {
		record.Filename = sourceRecord.Filename
		record.Size = sourceRecord.Size
		record.Hash = sourceRecord.Hash
		record.Suffix = sourceRecord.Suffix
	}
	if info, err := os.Stat(targetPath); err == nil {
		record.Size = info.Size()
	}
	if _, ret := db.InitFileRepo(tx).Create(record); !ret.OK {
		log.Logger.Warningf("Failed to record challenge attachment: challenge_id=%d reason=%v", target.ID, ret)
	}
	return model.SuccessRetVal()
}

// copyFileContent 复制文件内容，源文件不存在时返回 os.ErrNotExist。
func copyFileContent(source, target string) error {
	src, err := os.Open(source)
	if err != nil {
		return err
	}
	defer func() { _ = src.Close() }()
	dst, err := os.Create(target)
	if err != nil {
		return err
	}
	defer func() { _ = dst.Close() }()
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	return dst.Close()
}
