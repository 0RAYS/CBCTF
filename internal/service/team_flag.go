package service

import (
	"CBCTF/internal/db"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/task"
	"CBCTF/internal/utils"
	"errors"
	"fmt"

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
		_ = tx.Transaction(func(tx2 *gorm.DB) error {
			teamFlags, ret := CreateTeamFlag(tx2, team, contest, contestChallenge)
			if err, ok := ret.Attr["Error"]; ok && !ret.OK {
				return errors.New(err.(string))
			}
			if contestChallenge.Type == model.DynamicChallengeType {
				if _, err := task.EnqueueGenAttachmentTask(team.CaptainID, contestChallenge.Challenge, team, teamFlags); err != nil {
					log.Logger.Warningf("Failed to enqueue gen attachment task: %s", err)
					return err
				}
			}
			return nil
		})
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
			return nil, ret
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
	teamFlagRepo := db.InitTeamFlagRepo(tx)
	for _, contestFlag := range contestFlags {
		if _, ret := teamFlagRepo.Get(db.GetOptions{
			Conditions: map[string]any{"team_id": team.ID, "contest_flag_id": contestFlag.ID},
		}); !ret.OK {
			return false
		}
	}
	return true
}
