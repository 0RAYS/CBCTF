package db

import (
	"io"
	"os"
	"path/filepath"

	"gorm.io/gorm"

	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
)

// 比赛题目深拷贝迁移。
//
// 旧的比赛题目只复制了展示字段，运行期仍回查题库行，题库一改比赛就受影响。
// 迁移把比赛题目补齐为独立副本：
//   - contest_challenges 追加 rand_id / template / generator_image / network_policies
//   - contest_flags 追加 binding / generation / active
//   - 旧行按题库题目回填副本字段、公开 id 与 flag 注入绑定
//   - 移除「同一比赛不可重复添加同一题库题目」的唯一索引
//
// 迁移以 rand_id 是否为空作为标记：已经有 rand_id 的比赛题目视为已迁移，不会被再次覆盖，
// 因此迁移可以重复执行。
func MigrateContestChallengeCopies(tx *gorm.DB) model.RetVal {
	exists, ret := tableExists(tx, "contest_challenges")
	if !ret.OK {
		return ret
	}
	if !exists {
		// 全新安装：表结构由 AutoMigrate 一次性建立，无需回填
		return model.SuccessRetVal()
	}
	if ret = addContestChallengeCopyColumns(tx); !ret.OK {
		return ret
	}
	if ret = dropLegacyContestChallengeUniqueIndex(tx); !ret.OK {
		return ret
	}
	return backfillContestChallengeCopies(tx)
}

// addContestChallengeCopyColumns 追加副本相关列。
// 唯一索引依赖 rand_id 有值，因此列必须先于 AutoMigrate 建立。
func addContestChallengeCopyColumns(tx *gorm.DB) model.RetVal {
	statements := []string{
		`ALTER TABLE contest_challenges ADD COLUMN IF NOT EXISTS rand_id varchar(36)`,
		`ALTER TABLE contest_challenges ADD COLUMN IF NOT EXISTS template jsonb DEFAULT '{}'`,
		`ALTER TABLE contest_challenges ADD COLUMN IF NOT EXISTS generator_image text`,
		`ALTER TABLE contest_challenges ADD COLUMN IF NOT EXISTS network_policies jsonb`,
		`ALTER TABLE contest_flags ADD COLUMN IF NOT EXISTS binding jsonb DEFAULT '{}'`,
		`ALTER TABLE contest_flags ADD COLUMN IF NOT EXISTS generation integer DEFAULT 1`,
		`ALTER TABLE contest_flags ADD COLUMN IF NOT EXISTS active boolean DEFAULT true`,
	}
	for _, statement := range statements {
		if res := tx.Exec(statement); res.Error != nil {
			log.Logger.Warningf("Failed to add contest challenge copy column: %s error=%v", statement, res.Error)
			return model.RetVal{Msg: i18n.Model.ContestChallenge.GetError, Attr: map[string]any{"Error": res.Error.Error()}}
		}
	}
	return model.SuccessRetVal()
}

// dropLegacyContestChallengeUniqueIndex 移除旧的唯一索引。
// 比赛题目改为独立副本后，同一题库题目允许在同场比赛多次添加。
func dropLegacyContestChallengeUniqueIndex(tx *gorm.DB) model.RetVal {
	if res := tx.Exec(`DROP INDEX IF EXISTS idx_contest_challenges_unique_active`); res.Error != nil {
		log.Logger.Warningf("Failed to drop legacy contest challenge unique index: %v", res.Error)
		return model.RetVal{Msg: i18n.Model.ContestChallenge.GetError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}

// backfillContestChallengeCopies 为旧的比赛题目补齐副本字段与公开 id。
func backfillContestChallengeCopies(tx *gorm.DB) model.RetVal {
	type pendingCopy struct {
		ID          uint
		ContestID   uint
		ChallengeID uint
		Type        model.ChallengeType
	}
	var pending []pendingCopy
	if res := tx.Model(&model.ContestChallenge{}).
		Select("id", "contest_id", "challenge_id", "type").
		Where("rand_id IS NULL OR rand_id = ''").
		Find(&pending); res.Error != nil {
		log.Logger.Warningf("Failed to list contest challenges pending migration: %v", res.Error)
		return model.RetVal{Msg: i18n.Model.ContestChallenge.GetError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	if len(pending) == 0 {
		return model.SuccessRetVal()
	}
	migrated := make([]uint, 0, len(pending))
	for _, item := range pending {
		randID := utils.UUID()
		challenge, ret := InitChallengeRepo(tx).GetByID(item.ChallengeID, GetOptions{})
		if !ret.OK {
			// 题库题目已被删除：只补公开 id，副本字段留空
			log.Logger.Warningf("Contest challenge %d references missing challenge %d", item.ID, item.ChallengeID)
			if res := tx.Model(&model.ContestChallenge{}).Where("id = ?", item.ID).Update("rand_id", randID); res.Error != nil {
				log.Logger.Warningf("Failed to migrate contest challenge %d: %v", item.ID, res.Error)
				return model.RetVal{Msg: i18n.Model.ContestChallenge.GetError, Attr: map[string]any{"Error": res.Error.Error()}}
			}
			migrated = append(migrated, item.ID)
			continue
		}
		if res := tx.Model(&model.ContestChallenge{}).Where("id = ?", item.ID).
			Select("rand_id", "template", "generator_image", "network_policies").
			Updates(&model.ContestChallenge{
				RandID:          randID,
				Template:        challenge.Template,
				GeneratorImage:  challenge.GeneratorImage,
				NetworkPolicies: challenge.NetworkPolicies,
			}); res.Error != nil {
			log.Logger.Warningf("Failed to migrate contest challenge %d: %v", item.ID, res.Error)
			return model.RetVal{Msg: i18n.Model.ContestChallenge.GetError, Attr: map[string]any{"Error": res.Error.Error()}}
		}
		migrated = append(migrated, item.ID)
		copyContestChallengeAttachment(tx, challenge, item.ID)
	}
	// 回填比赛 flag 的注入绑定，绑定只在题目加入比赛时复制一次
	if len(migrated) > 0 {
		if res := tx.Exec(
			`UPDATE contest_flags SET binding = challenge_flags.binding FROM challenge_flags WHERE challenge_flags.id = contest_flags.challenge_flag_id AND contest_flags.deleted_at IS NULL AND contest_flags.contest_challenge_id IN ?`,
			migrated,
		); res.Error != nil {
			log.Logger.Warningf("Failed to migrate contest flag bindings: %v", res.Error)
			return model.RetVal{Msg: i18n.Model.ContestFlag.GetError, Attr: map[string]any{"Error": res.Error.Error()}}
		}
	}
	return model.SuccessRetVal()
}

// copyContestChallengeAttachment 把题库附件复制到比赛题目目录，并登记文件记录。
// 源文件缺失时只记录告警，不阻断迁移。
func copyContestChallengeAttachment(tx *gorm.DB, challenge model.Challenge, contestChallengeID uint) {
	var source, target string
	contestChallenge := model.ContestChallenge{
		BaseModel: model.BaseModel{ID: contestChallengeID},
		Type:      challenge.Type,
	}
	switch challenge.Type {
	case model.DynamicChallengeType:
		source, target = challenge.GeneratorPath(), contestChallenge.GeneratorPath()
	default:
		source, target = challenge.StaticPath(), contestChallenge.StaticPath()
	}
	sourceRecord, sourceRet := InitFileRepo(tx).Get(GetOptions{
		Conditions: map[string]any{
			"model":    model.Name(challenge),
			"model_id": challenge.ID,
			"type":     model.ChallengeFileType,
		},
	})
	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		log.Logger.Warningf("Failed to create contest challenge dir: contest_challenge_id=%d error=%v", contestChallengeID, err)
		return
	}
	if err := copyFileContent(source, target); err != nil {
		if !os.IsNotExist(err) {
			log.Logger.Warningf("Failed to copy challenge attachment: challenge_id=%d contest_challenge_id=%d error=%v", challenge.ID, contestChallengeID, err)
		}
		return
	}
	record := model.File{
		RandID:   utils.UUID(),
		Filename: filepath.Base(target),
		Path:     model.FilePath(target),
		Model:    model.Name(model.ContestChallenge{}),
		ModelID:  contestChallengeID,
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
	if _, ret := InitFileRepo(tx).Create(record); !ret.OK {
		log.Logger.Warningf("Failed to record contest challenge attachment: contest_challenge_id=%d", contestChallengeID)
	}
}

// copyFileContent 复制文件内容，源文件不存在时返回 os.ErrNotExist。
func copyFileContent(source, target string) error {
	src, err := os.Open(source)
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.Create(target)
	if err != nil {
		return err
	}
	defer dst.Close()
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	return dst.Close()
}

// tableExists 判断表是否已经建立。
func tableExists(tx *gorm.DB, name string) (bool, model.RetVal) {
	var exists bool
	if res := tx.Raw(`SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = current_schema() AND table_name = ?)`, name).Scan(&exists); res.Error != nil {
		log.Logger.Warningf("Failed to check table %s: %v", name, res.Error)
		return false, model.RetVal{Msg: i18n.DB.Unavailable, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return exists, model.SuccessRetVal()
}
