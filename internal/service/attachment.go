package service

import (
	"gorm.io/gorm"

	"CBCTF/internal/db"
	"CBCTF/internal/model"
)

// AttachmentPath 返回题目附件的路径。
// 比赛内使用比赛题目副本，dynamic 题按队伍在该副本上的 flag 列表算缓存 key。
func AttachmentPath(tx *gorm.DB, source model.GeneratorSource, teamID uint) (string, model.RetVal) {
	if source.SourceType() != model.DynamicChallengeType {
		return source.SourceAttachmentPath(teamID), model.SuccessRetVal()
	}
	values, ret := sourceFlagValues(tx, source, teamID)
	if !ret.OK {
		return "", ret
	}
	return source.SourceAttachmentCachePath(teamID, values), model.SuccessRetVal()
}

// sourceFileRecord 返回题目来源登记的附件文件记录。
func sourceFileRecord(tx *gorm.DB, source model.GeneratorSource) (model.File, model.RetVal) {
	return db.InitFileRepo(tx).Get(db.GetOptions{
		Conditions: map[string]any{
			"model":    model.Name(source),
			"model_id": source.SourceID(),
			"type":     model.ChallengeFileType,
		},
	})
}

// teamScopedFlagValues 返回队伍持有的 flag 值，scope 限定了队伍 flag 指向的
// flag，idColumn 指明队伍 flag 关联的 flag 表字段。
func teamScopedFlagValues(tx *gorm.DB, teamID uint, scope *gorm.DB, idColumn string) ([]string, model.RetVal) {
	teamFlags, _, ret := db.InitTeamFlagRepo(scope).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"team_id": teamID},
		Sort:       []string{idColumn + " ASC"},
	})
	if !ret.OK {
		return nil, ret
	}
	values := make([]string, 0, len(teamFlags))
	for _, flag := range teamFlags {
		values = append(values, flag.Value)
	}
	return values, model.SuccessRetVal()
}

// sourceFlagValues 返回附件生成所需的 flag 列表。
// 比赛题目副本取该副本生效中的 flag，题库题目取题库 flag（仅测试模式）。
func sourceFlagValues(tx *gorm.DB, source model.GeneratorSource, teamID uint) ([]string, model.RetVal) {
	if source.SourceContestChallengeID() > 0 {
		contestFlags, _, ret := db.InitContestFlagRepo(tx).List(-1, -1, db.GetOptions{
			Conditions: map[string]any{
				"contest_challenge_id": source.SourceContestChallengeID(),
				"active":               true,
			},
			Sort: []string{"id ASC"},
		})
		if !ret.OK {
			return nil, ret
		}
		if teamID == 0 {
			values := make([]string, 0, len(contestFlags))
			for _, flag := range contestFlags {
				values = append(values, flag.Value)
			}
			return values, model.SuccessRetVal()
		}
		ids := make([]uint, 0, len(contestFlags))
		for _, flag := range contestFlags {
			ids = append(ids, flag.ID)
		}
		if len(ids) == 0 {
			return nil, model.SuccessRetVal()
		}
		return teamScopedFlagValues(tx, teamID, tx.Where("contest_flag_id IN (?)", ids), "contest_flag_id")
	}
	flags, _, ret := db.InitChallengeFlagRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"challenge_id": source.SourceChallengeID()},
		Sort:       []string{"id ASC"},
	})
	if !ret.OK {
		return nil, ret
	}
	if teamID == 0 {
		values := make([]string, 0, len(flags))
		for _, flag := range flags {
			values = append(values, flag.Value)
		}
		return values, model.SuccessRetVal()
	}
	ids := tx.Model(&model.ChallengeFlag{}).Select("id").Where("challenge_id = ?", source.SourceChallengeID())
	return teamScopedFlagValues(tx, teamID, tx.Where("challenge_flag_id IN (?)", ids), "challenge_flag_id")
}
