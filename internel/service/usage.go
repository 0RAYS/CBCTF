package service

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"errors"
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
		_ = tx.Transaction(func(tx *gorm.DB) error {
			tx2 := tx.Begin()
			challenge, ok, msg := challengeRepo.GetByID(challengeID, false, 0)
			if !ok {
				failed = append(failed, challengeID)
				return errors.New(msg)
			}
			usage, ok, msg := usageRepo.Create(db.CreateUsageOptions{
				ContestID:   contest.ID,
				ChallengeID: challengeID,
				Name:        challenge.Name,
				Desc:        challenge.Desc,
				Docker:      challenge.Docker,
				Dockers:     challenge.Dockers,
				Attempt:     0,
			})
			if !ok {
				failed = append(failed, challengeID)
				return errors.New(msg)
			}
			options := db.CreateFlagOptions{
				ContestID:    contest.ID,
				UsageID:      usage.ID,
				Score:        1000,
				CurrentScore: 1000,
				Decay:        100,
				MinScore:     100,
				ScoreType:    0,
				Blood:        make(model.Uints, 3),
			}
			switch challenge.Type {
			case model.StaticChallenge, model.DynamicChallenge:
				for _, s := range challenge.Flags {
					options.Value = s
					_, ok, msg = flagRepo.Create(options)
					if !ok {
						failed = append(failed, challengeID)
						return errors.New(msg)
					}
				}
			case model.DockerChallenge:
				for _, s := range challenge.Docker.Flags {
					options.Value = s
					flag, ok, msg := flagRepo.Create(options)
					if !ok {
						failed = append(failed, challengeID)
						return errors.New(msg)
					}
					usage.Docker.FlagsID = append(usage.Docker.FlagsID, flag.ID)
				}
				if ok, msg := usageRepo.Update(usage.ID, db.UpdateUsageOptions{
					Docker: &challenge.Docker,
				}); !ok {
					failed = append(failed, challengeID)
					return errors.New(msg)
				}
			case model.DockersChallenge:
				for i, docker := range challenge.Dockers {
					for _, s := range docker.Flags {
						options.Value = s
						flag, ok, _ := flagRepo.Create(options)
						if !ok {
							failed = append(failed, challengeID)
							tx2.Rollback()
							break
						}
						usage.Dockers[i].FlagsID = append(usage.Dockers[i].FlagsID, flag.ID)
					}
				}
				if ok, msg := usageRepo.Update(usage.ID, db.UpdateUsageOptions{
					Dockers: &challenge.Dockers,
				}); !ok {
					failed = append(failed, challengeID)
					return errors.New(msg)
				}
			default:
				failed = append(failed, challengeID)
				return errors.New("InvalidChallengeType")
			}
			usages = append(usages, usage)
			return nil
		})
	}
	return usages, failed, true, "Success"
}

func UpdateUsage(tx *gorm.DB, usage model.Usage, form f.UpdateUsageForm) (bool, string) {
	repo := db.InitUsageRepo(tx)
	return repo.Update(usage.ID, db.UpdateUsageOptions{
		Name:    form.Name,
		Desc:    form.Desc,
		Attempt: form.Attempt,
		Hidden:  form.Hidden,
		Hints:   form.Hints,
		Tags:    form.Tags,
		Docker:  form.Docker,
		Dockers: form.Dockers,
	})
}

func DeleteUsage(tx *gorm.DB, usage model.Usage) (bool, string) {
	repo := db.InitUsageRepo(tx)
	return repo.Delete(usage.ID)
}
