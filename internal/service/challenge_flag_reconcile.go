package service

import (
	"gorm.io/gorm"

	"CBCTF/internal/db"
	"CBCTF/internal/model"
)

// challengeFlagMatch 描述一条派生 flag 与题库现有 flag 的对齐结果。
type challengeFlagMatch struct {
	// ExistingID 命中的题库 flag id，为 0 表示需要新建
	ExistingID uint
	// Value 最终的 flag 值模板
	Value string
	// Binding 最终的注入绑定
	Binding model.FlagBinding
	// Changed 是否需要写回题库
	Changed bool
	// Created 是否需要新建
	Created bool
}

// matchChallengeFlags 把 compose 派生出的 flag 与题库现有 flag 对齐。
//
// 对齐规则：
//   - 值与绑定完全一致：直接复用，flag id 不变
//   - 绑定一致但值变化：原地更新值，flag id 不变（附件缓存 key 因此保持稳定）
//   - 其余：新建 flag
//
// 返回值与派生 flag 一一对应，removed 是已消失、需要软删除的题库 flag id。
func matchChallengeFlags(existing []model.ChallengeFlag, derived []model.ChallengeFlag) ([]challengeFlagMatch, []uint) {
	matches := make([]challengeFlagMatch, len(derived))
	for i, flag := range derived {
		matches[i] = challengeFlagMatch{Value: flag.Value, Binding: flag.Binding, Created: true}
	}
	existingUsed := make([]bool, len(existing))

	// 值与绑定都一致：复用现有 flag
	for i, flag := range derived {
		for j, old := range existing {
			if existingUsed[j] {
				continue
			}
			if old.Value == flag.Value && old.Binding == flag.Binding {
				existingUsed[j] = true
				matches[i] = challengeFlagMatch{ExistingID: old.ID, Value: old.Value, Binding: old.Binding}
				break
			}
		}
	}
	// 绑定一致但值变化：原地更新，保留 flag id
	for i, flag := range derived {
		if !matches[i].Created {
			continue
		}
		for j, old := range existing {
			if existingUsed[j] || old.Binding != flag.Binding {
				continue
			}
			existingUsed[j] = true
			matches[i] = challengeFlagMatch{ExistingID: old.ID, Value: flag.Value, Binding: flag.Binding, Changed: true}
			break
		}
	}
	removed := make([]uint, 0)
	for j, old := range existing {
		if !existingUsed[j] {
			removed = append(removed, old.ID)
		}
	}
	return matches, removed
}

// ReconcileChallengeFlags 用 compose 派生出的 flag 集合与题库现有 flag 对齐。
//
// 比赛题目在加入时即持有 flag 副本，因此题库侧的增删改不会再影响任何比赛的解题记录，
// 这里不再需要"删了重建"。
func ReconcileChallengeFlags(tx *gorm.DB, challenge model.Challenge, derived []model.ChallengeFlag) model.RetVal {
	repo := db.InitChallengeFlagRepo(tx)
	matches, removed := matchChallengeFlags(challenge.ChallengeFlags, derived)
	for _, match := range matches {
		if match.Created {
			if _, ret := repo.Create(model.ChallengeFlag{
				ChallengeID: challenge.ID,
				Value:       match.Value,
				Binding:     match.Binding,
			}); !ret.OK {
				return ret
			}
			continue
		}
		if !match.Changed {
			continue
		}
		if ret := repo.Update(match.ExistingID, db.UpdateChallengeFlagOptions{
			Value:   &match.Value,
			Binding: &match.Binding,
		}); !ret.OK {
			return ret
		}
	}
	// 消失的 flag：只软删除题库侧
	if len(removed) > 0 {
		if ret := repo.Delete(removed...); !ret.OK {
			return ret
		}
	}
	return model.SuccessRetVal()
}
