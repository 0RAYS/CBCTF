package service

import (
	"os"

	"gorm.io/gorm"

	"CBCTF/internal/db"
	"CBCTF/internal/model"
)

// AttachmentPath 返回题目附件的路径。
// dynamic 题按队伍在该题目上的 flag 列表算缓存 key。
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

// buildSourceFileName 返回附件的展示文件名，题目各自登记自己的附件记录。
func buildSourceFileName(tx *gorm.DB, source model.GeneratorSource, teamID uint) string {
	path, ret := AttachmentPath(tx, source, teamID)
	if !ret.OK {
		return ""
	}
	record, _ := sourceFileRecord(tx, source)
	filename := model.AttachmentFileName
	if string(record.Path) == path && record.Filename != "" {
		filename = record.Filename
	}
	if _, err := os.Stat(path); err != nil {
		return ""
	}
	return filename
}

// teamScopedFlagValues 返回队伍持有的 flag 值，scope 限定了队伍 flag 指向的
// flag，idColumn 指明队伍 flag 关联的 flag 表字段。
func teamScopedFlagValues(teamID uint, scope *gorm.DB, idColumn string) ([]string, model.RetVal) {
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
// 队伍为 0 时取题目自身生效中的 flag（测试模式），否则取队伍实际持有的答案。
func sourceFlagValues(tx *gorm.DB, source model.GeneratorSource, teamID uint) ([]string, model.RetVal) {
	flags, _, ret := db.InitChallengeFlagRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"challenge_id": source.SourceID(), "active": true},
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
	ids := make([]uint, 0, len(flags))
	for _, flag := range flags {
		ids = append(ids, flag.ID)
	}
	if len(ids) == 0 {
		return nil, model.SuccessRetVal()
	}
	return teamScopedFlagValues(teamID, tx.Where("challenge_flag_id IN (?)", ids), "challenge_flag_id")
}
