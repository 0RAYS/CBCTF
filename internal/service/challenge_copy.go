package service

import (
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

// CopyChallengeToContest 把题库题目深拷贝成比赛题目副本。
// 副本包含展示字段、靶机规格、生成器镜像、网络策略、flag（值与注入绑定）以及附件文件，
// 比赛运行期只依赖副本，ChallengeID 仅用于溯源。
func CopyChallengeToContest(tx *gorm.DB, contest model.Contest, challenge model.Challenge) (model.ContestChallenge, model.RetVal) {
	contestChallenge, ret := db.InitContestChallengeRepo(tx).Create(model.ContestChallenge{
		ContestID:       contest.ID,
		ChallengeID:     challenge.ID,
		RandID:          utils.UUID(),
		Name:            challenge.Name,
		Description:     challenge.Description,
		Type:            challenge.Type,
		Category:        challenge.Category,
		Hidden:          true,
		Template:        challenge.Template.DeepCopy(),
		GeneratorImage:  challenge.GeneratorImage,
		NetworkPolicies: challenge.NetworkPolicies.DeepCopy(),
	})
	if !ret.OK {
		return model.ContestChallenge{}, ret
	}
	contestFlagRepo := db.InitContestFlagRepo(tx)
	for _, flag := range challenge.ChallengeFlags {
		if _, ret = contestFlagRepo.Create(model.ContestFlag{
			ContestID:          contest.ID,
			ContestChallengeID: contestChallenge.ID,
			ChallengeFlagID:    flag.ID,
			Value:              flag.Value,
			Binding:            flag.Binding,
			Generation:         1,
			Active:             true,
			Score:              1000,
			CurrentScore:       1000,
			Decay:              100,
			MinScore:           100,
			ScoreType:          0,
			Solvers:            0,
			Last:               time.Now(),
		}); !ret.OK {
			return model.ContestChallenge{}, ret
		}
	}
	if ret = copyChallengeFiles(tx, challenge, contestChallenge); !ret.OK {
		return model.ContestChallenge{}, ret
	}
	return contestChallenge, model.SuccessRetVal()
}

// copyChallengeFiles 把题库题目的附件 / 生成器复制到比赛题目目录，并登记 files 记录。
// 源文件缺失时只记录告警，不阻断加入比赛。
func copyChallengeFiles(tx *gorm.DB, challenge model.Challenge, contestChallenge model.ContestChallenge) model.RetVal {
	var source, target string
	switch challenge.Type {
	case model.DynamicChallengeType:
		source, target = challenge.GeneratorPath(), contestChallenge.GeneratorPath()
	default:
		source, target = challenge.StaticPath(), contestChallenge.StaticPath()
	}
	sourceRecord, sourceRet := sourceFileRecord(tx, challenge)
	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		log.Logger.Warningf("Failed to create contest challenge dir: contest_challenge_id=%d error=%v", contestChallenge.ID, err)
		return model.SuccessRetVal()
	}
	if err := copyFileContent(source, target); err != nil {
		if !os.IsNotExist(err) {
			log.Logger.Warningf("Failed to copy challenge attachment: challenge_id=%d contest_challenge_id=%d error=%v", challenge.ID, contestChallenge.ID, err)
		}
		return model.SuccessRetVal()
	}
	record := model.File{
		RandID:   utils.UUID(),
		Filename: filepath.Base(target),
		Size:     0,
		Path:     model.FilePath(target),
		Model:    model.Name(contestChallenge),
		ModelID:  contestChallenge.ID,
		Suffix:   filepath.Ext(target),
		Type:     model.ChallengeFileType,
	}
	if sourceRet.OK {
		record.Filename = sourceRecord.Filename
		record.Size = sourceRecord.Size
		record.Hash = sourceRecord.Hash
		record.Suffix = sourceRecord.Suffix
	}
	if info, err := os.Stat(target); err == nil {
		record.Size = info.Size()
	}
	if _, ret := db.InitFileRepo(tx).Create(record); !ret.OK {
		log.Logger.Warningf("Failed to record contest challenge attachment: contest_challenge_id=%d reason=%v", contestChallenge.ID, ret)
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
