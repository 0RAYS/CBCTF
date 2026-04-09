package service

import (
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/task"
	"CBCTF/internal/utils"
	"CBCTF/internal/view"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

func CreateTeamFlags(tx *gorm.DB, team model.Team, contest model.Contest) model.RetVal {
	contestChallenges, _, ret := db.InitContestChallengeRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_id": contest.ID},
		Preloads:   map[string]db.GetOptions{"ContestFlags": {}, "Challenge": {}},
	})
	if !ret.OK {
		return ret
	}
	for _, contestChallenge := range contestChallenges {
		ret = db.WithTransactionDB(tx, func(tx2 *gorm.DB) model.RetVal {
			teamFlags, ret := CreateTeamFlag(tx2, team, contest, contestChallenge)
			if !ret.OK {
				return ret
			}
			if contestChallenge.Type == model.DynamicChallengeType {
				generator, ret := GetGenerator(tx2, contest.ID, contestChallenge.Challenge)
				if !ret.OK {
					return model.RetVal{
						Msg: i18n.Model.CreateError,
						Attr: map[string]any{
							"Model": model.ModelName(model.TeamFlag{}),
							"Error": fmt.Sprintf("generate attachment failed: %s", ret.Msg),
						},
					}
				}
				if _, err := task.EnqueueGenAttachmentTask(team.CaptainID, generator, contestChallenge.Challenge, team, teamFlags); err != nil {
					log.Logger.Warningf("Failed to enqueue gen attachment task: %s", err)
					return model.RetVal{
						Msg: i18n.Model.CreateError,
						Attr: map[string]any{
							"Model": model.ModelName(model.TeamFlag{}),
							"Error": err.Error(),
						},
					}
				}
			}
			return model.SuccessRetVal()
		})
		if !ret.OK {
			if ret.Attr == nil {
				ret.Attr = map[string]any{}
			}
			ret.Attr["Model"] = model.ModelName(model.TeamFlag{})
			if _, ok := ret.Attr["Error"]; !ok {
				ret.Attr["Error"] = ret.Msg
			}
			ret.Msg = i18n.Model.CreateError
			return ret
		}
	}
	return model.SuccessRetVal()
}

// CreateTeamFlag model.ContestChallenge Preload model.ContestFlag
func CreateTeamFlag(tx *gorm.DB, team model.Team, contest model.Contest, contestChallenge model.ContestChallenge) ([]model.TeamFlag, model.RetVal) {
	teamFlagRepo := db.InitTeamFlagRepo(tx)
	teamFlagL := make([]model.TeamFlag, 0)
	for _, contestFlag := range contestChallenge.ContestFlags {
		teamFlag, ret := teamFlagRepo.Get(db.GetOptions{
			Conditions: map[string]any{"team_id": team.ID, "contest_flag_id": contestFlag.ID},
		})
		if ret.OK {
			teamFlagL = append(teamFlagL, teamFlag)
			continue
		}
		options := db.CreateTeamFlagOptions{
			TeamID:          team.ID,
			ContestFlagID:   contestFlag.ID,
			ChallengeFlagID: contestFlag.ChallengeFlagID,
			Solved:          false,
		}
		if result := model.StaticFlagTmpl.FindAllStringSubmatch(contestFlag.Value, 1); len(result) > 0 {
			options.Value = result[0][1]
		} else if result = model.DynamicFlagTmpl.FindAllStringSubmatch(contestFlag.Value, 1); len(result) > 0 {
			options.Value = utils.RandFlag(result[0][1])
		} else if result = model.UUIDFlagTmpl.FindAllStringSubmatch(contestFlag.Value, 1); len(result) > 0 {
			options.Value = utils.UUID()
		} else {
			options.Value = contestFlag.Value
		}
		if prefix := contest.Prefix; prefix != "" && contestChallenge.Type != model.QuestionChallengeType {
			options.Value = fmt.Sprintf("%s{%s}", contest.Prefix, options.Value)
		}
		teamFlag, ret = teamFlagRepo.Create(options)
		if !ret.OK {
			errMsg, ok := ret.Attr["Error"].(string)
			if !ok || !strings.Contains(strings.ToLower(errMsg), "duplicate key") {
				return nil, ret
			}
			teamFlag, ret = teamFlagRepo.Get(db.GetOptions{
				Conditions: map[string]any{"team_id": team.ID, "contest_flag_id": contestFlag.ID},
			})
			if !ret.OK {
				return nil, ret
			}
		}
		teamFlagL = append(teamFlagL, teamFlag)
	}
	return teamFlagL, model.SuccessRetVal()
}

// UpdateTeamFlag model.ContestChallenge Preload model.ContestFlag
func UpdateTeamFlag(tx *gorm.DB, team model.Team, contest model.Contest, contestChallenge model.ContestChallenge) ([]model.TeamFlag, model.RetVal) {
	submissionRepo := db.InitSubmissionRepo(tx)
	submissions, _, ret := submissionRepo.List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"team_id": team.ID, "contest_challenge_id": contestChallenge.ID},
	})
	if !ret.OK {
		return nil, ret
	}
	submissionIDL := make([]uint, 0)
	for _, submission := range submissions {
		submissionIDL = append(submissionIDL, submission.ID)
	}
	teamFlagIDL := make([]uint, 0)
	teamFlagRepo := db.InitTeamFlagRepo(tx)
	for _, contestFlag := range contestChallenge.ContestFlags {
		teamFlag, ret := teamFlagRepo.Get(db.GetOptions{
			Conditions: map[string]any{"team_id": team.ID, "contest_flag_id": contestFlag.ID},
		})
		if !ret.OK {
			return nil, ret
		}
		teamFlagIDL = append(teamFlagIDL, teamFlag.ID)
	}
	if ret = submissionRepo.Delete(submissionIDL...); !ret.OK {
		return nil, ret
	}
	if ret = teamFlagRepo.Delete(teamFlagIDL...); !ret.OK {
		return nil, ret
	}
	return CreateTeamFlag(tx, team, contest, contestChallenge)
}

func CheckIfGenerated(tx *gorm.DB, team model.Team, contestFlags []model.ContestFlag) bool {
	if len(contestFlags) == 0 {
		return true
	}
	contestFlagIDL := make([]uint, 0, len(contestFlags))
	for _, contestFlag := range contestFlags {
		contestFlagIDL = append(contestFlagIDL, contestFlag.ID)
	}
	count, ret := db.InitTeamFlagRepo(tx).CountGenerated(team.ID, contestFlagIDL...)
	return ret.OK && count == int64(len(contestFlags))
}

func ListTeamFlagViews(tx *gorm.DB, team model.Team) ([]view.TeamFlagChallengeView, model.RetVal) {
	teamFlags, _, ret := db.InitTeamFlagRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"team_id": team.ID},
		Preloads: map[string]db.GetOptions{"ContestFlag": {
			Preloads: map[string]db.GetOptions{"ContestChallenge": {}},
		}},
	})
	if !ret.OK {
		return nil, ret
	}

	type groupData struct {
		index int
		view  view.TeamFlagChallengeView
	}
	groupMap := make(map[uint]*groupData)
	result := make([]view.TeamFlagChallengeView, 0)
	for _, flag := range teamFlags {
		id := flag.ContestFlag.ContestChallengeID
		group, ok := groupMap[id]
		if !ok {
			result = append(result, view.TeamFlagChallengeView{
				Name:     flag.ContestFlag.ContestChallenge.Name,
				Type:     flag.ContestFlag.ContestChallenge.Type,
				Category: flag.ContestFlag.ContestChallenge.Category,
				Hidden:   flag.ContestFlag.ContestChallenge.Hidden,
				Flags:    make([]view.TeamFlagInfoView, 0),
			})
			group = &groupData{index: len(result) - 1, view: result[len(result)-1]}
			groupMap[id] = group
		}
		result[group.index].Flags = append(result[group.index].Flags, view.TeamFlagInfoView{
			Value:        flag.Value,
			Solved:       flag.Solved,
			Template:     flag.ContestFlag.Value,
			InitScore:    flag.ContestFlag.Score,
			CurrentScore: flag.ContestFlag.CurrentScore,
			Decay:        flag.ContestFlag.Decay,
			MinScore:     flag.ContestFlag.MinScore,
			Solvers:      flag.ContestFlag.Solvers,
		})
	}
	return result, model.SuccessRetVal()
}

func InitTeamChallenge(tx *gorm.DB, user model.User, team model.Team, contest model.Contest, challenge model.Challenge, contestChallenge model.ContestChallenge) model.RetVal {
	contestFlags, _, ret := db.InitContestFlagRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_challenge_id": contestChallenge.ID},
	})
	if !ret.OK {
		return ret
	}
	contestChallenge.ContestFlags = contestFlags
	return db.WithTransactionDB(tx, func(tx2 *gorm.DB) model.RetVal {
		teamFlags, createRet := CreateTeamFlag(tx2, team, contest, contestChallenge)
		if !createRet.OK {
			return createRet
		}
		if challenge.Type != model.DynamicChallengeType {
			return model.SuccessRetVal()
		}

		generator, generatorRet := GetGenerator(tx2, contest.ID, challenge)
		if !generatorRet.OK {
			return generatorRet
		}
		if _, err := task.EnqueueGenAttachmentTask(user.ID, generator, challenge, team, teamFlags); err != nil {
			log.Logger.Warningf("Failed to enqueue gen attachment task: %s", err)
			return model.RetVal{Msg: i18n.Task.EnqueueError, Attr: map[string]any{"Error": err.Error()}}
		}
		return model.SuccessRetVal()
	})
}

func ResetTeamChallenge(tx *gorm.DB, user model.User, team model.Team, contest model.Contest, challenge model.Challenge, contestChallenge model.ContestChallenge) model.RetVal {
	contestFlags, _, ret := db.InitContestFlagRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_challenge_id": contestChallenge.ID},
	})
	if !ret.OK {
		return ret
	}
	contestChallenge.ContestFlags = contestFlags
	ret = db.WithTransactionDB(tx, func(tx2 *gorm.DB) model.RetVal {
		teamFlags, updateRet := UpdateTeamFlag(tx2, team, contest, contestChallenge)
		if !updateRet.OK {
			return updateRet
		}
		if challenge.Type != model.DynamicChallengeType {
			return model.SuccessRetVal()
		}

		generator, generatorRet := GetGenerator(tx2, contest.ID, challenge)
		if !generatorRet.OK {
			return generatorRet
		}
		if _, err := task.EnqueueGenAttachmentTask(user.ID, generator, challenge, team, teamFlags); err != nil {
			log.Logger.Warningf("Failed to enqueue gen attachment task: %s", err)
			return model.RetVal{Msg: i18n.Task.EnqueueError, Attr: map[string]any{"Error": err.Error()}}
		}
		return model.SuccessRetVal()
	})
	if !ret.OK {
		return ret
	}

	if challenge.Type == model.PodsChallengeType {
		go func() {
			victim, victimRet := db.InitVictimRepo(tx).HasAliveVictim(team.ID, challenge.ID)
			if !victimRet.OK {
				return
			}
			_ = ForceStopVictim(tx, victim)
		}()
	}
	return model.SuccessRetVal()
}
