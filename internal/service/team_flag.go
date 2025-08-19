package service

import (
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/task"
	"CBCTF/internal/utils"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func CreateTeamFlags(tx *gorm.DB, team model.Team, contest model.Contest) (bool, string) {
	contestChallenges, _, ok, msg := db.InitContestChallengeRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_id": contest.ID},
		Preloads:   map[string]db.GetOptions{"ContestFlags": {}, "Challenge": {}},
	})
	if !ok {
		return false, msg
	}
	for _, contestChallenge := range contestChallenges {
		_ = tx.Transaction(func(tx2 *gorm.DB) error {
			teamFlags, ok, msg := CreateTeamFlag(tx2, team, contestChallenge)
			if !ok {
				return errors.New(msg)
			}
			if contestChallenge.Type == model.DynamicChallengeType {
				if _, err := task.EnqueueGenAttachmentTask(team.CaptainID, contestChallenge, team, teamFlags); err != nil {
					log.Logger.Warningf("Failed to enqueue gen attachment task: %s", err)
					return errors.New(i18n.EnqueueTaskError)
				}
			}
			return nil
		})
	}
	return true, i18n.Success
}

// CreateTeamFlag 需要预加载 ContestFlags
func CreateTeamFlag(tx *gorm.DB, team model.Team, contestChallenge model.ContestChallenge) ([]model.TeamFlag, bool, string) {
	teamFlagRepo := db.InitTeamFlagRepo(tx)
	teamFlagL := make([]model.TeamFlag, 0)
	for _, contestFlag := range contestChallenge.ContestFlags {
		teamFlag, ok, msg := teamFlagRepo.Get(db.GetOptions{
			Conditions: map[string]any{"team_id": team.ID, "contest_flag_id": contestFlag.ID},
		})
		if ok {
			teamFlagL = append(teamFlagL, teamFlag)
			continue
		}
		options := db.CreateTeamFlagOptions{
			TeamID:          team.ID,
			ContestFlagID:   contestFlag.ID,
			ChallengeFlagID: contestFlag.ChallengeFlagID,
			Solved:          false,
		}
		if result := model.StaticFlag.FindAllStringSubmatch(contestFlag.Value, 1); len(result) > 0 {
			options.Value = result[0][1]
		} else if result = model.DynamicFlag.FindAllStringSubmatch(contestFlag.Value, 1); len(result) > 0 {
			options.Value = utils.RandFlag(result[0][1])
		} else if result = model.UUIDFlag.FindAllStringSubmatch(contestFlag.Value, 1); len(result) > 0 {
			options.Value = utils.UUID()
		} else {
			options.Value = contestFlag.Value
		}
		if prefix := contestChallenge.Contest.Prefix; prefix != "" && contestChallenge.Type != model.QuestionChallengeType {
			options.Value = fmt.Sprintf("%s{%s}", contestChallenge.Contest.Prefix, options.Value)
		}
		teamFlag, ok, msg = teamFlagRepo.Create(options)
		if !ok {
			return teamFlagL, false, msg
		}
		teamFlagL = append(teamFlagL, teamFlag)
	}
	return teamFlagL, true, i18n.Success
}

// UpdateTeamFlag 需要预加载 ContestFlags
func UpdateTeamFlag(tx *gorm.DB, team model.Team, contestChallenge model.ContestChallenge) ([]model.TeamFlag, bool, string) {
	submissionRepo := db.InitSubmissionRepo(tx)
	submissions, _, ok, msg := submissionRepo.List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"team_id": team.ID, "contest_challenge_id": contestChallenge.ID},
	})
	if !ok {
		return make([]model.TeamFlag, 0), false, msg
	}
	submissionIDL := make([]uint, 0)
	for _, submission := range submissions {
		submissionIDL = append(submissionIDL, submission.ID)
	}
	teamFlagIDL := make([]uint, 0)
	teamFlagRepo := db.InitTeamFlagRepo(tx)
	for _, contestFlag := range contestChallenge.ContestFlags {
		teamFlag, ok, msg := teamFlagRepo.Get(db.GetOptions{
			Conditions: map[string]any{"team_id": team.ID, "contest_flag_id": contestFlag.ID},
		})
		if !ok {
			return make([]model.TeamFlag, 0), false, msg
		}
		teamFlagIDL = append(teamFlagIDL, teamFlag.ID)
	}
	if ok, msg = submissionRepo.Delete(submissionIDL...); !ok {
		return make([]model.TeamFlag, 0), false, msg
	}
	if ok, msg = teamFlagRepo.Delete(teamFlagIDL...); !ok {
		return make([]model.TeamFlag, 0), false, msg
	}
	return CreateTeamFlag(tx, team, contestChallenge)
}

// CheckIfGenerated contestChallenge 需要预加载 ContestFlags
func CheckIfGenerated(tx *gorm.DB, team model.Team, contestChallenge model.ContestChallenge) bool {
	teamFlagRepo := db.InitTeamFlagRepo(tx)
	for _, contestFlag := range contestChallenge.ContestFlags {
		if _, ok, _ := teamFlagRepo.Get(db.GetOptions{
			Conditions: map[string]any{"team_id": team.ID, "contest_flag_id": contestFlag.ID},
		}); !ok {
			return false
		}
	}
	return true
}
