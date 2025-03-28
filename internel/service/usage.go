package service

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"gorm.io/gorm"
)

// CreateUsage 批量将题目添加至比赛, 采用局部回滚, 返回值 bool string 永为 true "Success"
func CreateUsage(tx *gorm.DB, contest model.Contest, form f.CreateUsageForm) ([]model.Usage, []string, bool, string) {
	usageRepo := db.InitUsageRepo(tx)
	challengeRepo := db.InitChallengeRepo(tx)
	flagRepo := db.InitFlagRepo(tx)
	usages := make([]model.Usage, 0)
	failed := make([]string, 0)
	for _, challengeID := range form.ChallengeID {
		// 局部回滚
		tx2 := tx.Begin()
		challenge, ok, _ := challengeRepo.GetByID(challengeID, false, 0)
		if !ok {
			failed = append(failed, challengeID)
			tx2.Rollback()
			continue
		}
		usage, ok, _ := usageRepo.Create(db.CreateUsageOptions{
			ContestID:   contest.ID,
			ChallengeID: challengeID,
			Name:        challenge.Name,
			Desc:        challenge.Desc,
		})
		if !ok {
			failed = append(failed, challengeID)
			tx2.Rollback()
			continue
		}
		options := db.CreateFlagOptions{
			ContestID:    contest.ID,
			UsageID:      usage.ID,
			Score:        1000,
			CurrentScore: 1000,
			Decay:        100,
			MinScore:     100,
			ScoreType:    0,
			Attempt:      0,
		}
		switch challenge.Type {
		case model.StaticChallenge, model.DynamicChallenge:
			for _, flag := range challenge.Flags {
				options.Value = flag
				_, ok, _ = flagRepo.Create(options)
				if !ok {
					failed = append(failed, challengeID)
					tx2.Rollback()
					break
				}
			}
		case model.DockerChallenge:
			for _, flag := range challenge.Docker.Flags {
				options.Value = flag
				_, ok, _ = flagRepo.Create(options)
				if !ok {
					failed = append(failed, challengeID)
					tx2.Rollback()
					break
				}
			}
		case model.DockersChallenge:
			for _, docker := range challenge.Dockers {
				for _, flag := range docker.Flags {
					options.Value = flag
					_, ok, _ = flagRepo.Create(options)
					if !ok {
						failed = append(failed, challengeID)
						tx2.Rollback()
						break
					}
				}
			}
		default:
			failed = append(failed, challengeID)
			tx2.Rollback()
		}
		tx2.Commit()
		usages = append(usages, usage)
	}
	return usages, failed, true, "Success"
}

func UpdateUsage(tx *gorm.DB, usage model.Usage, form f.UpdateUsageForm) (bool, string) {
	repo := db.InitUsageRepo(tx)
	return repo.Update(usage.ID, db.UpdateUsageOptions{
		Name:  form.Name,
		Desc:  form.Desc,
		Hints: form.Hints,
		Tags:  form.Tags,
	})
}

func DeleteUsage(tx *gorm.DB, usage model.Usage) (bool, string) {
	repo := db.InitUsageRepo(tx)
	return repo.Delete(usage.ID)
}
