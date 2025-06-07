package service

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"CBCTF/internel/utils"
	"gorm.io/gorm"
)

// CreateTeamFlag 需要预加载 ContestFlags
func CreateTeamFlag(tx *gorm.DB, team model.Team, contestChallenge model.ContestChallenge) ([]model.TeamFlag, bool, string) {
	teamFlagRepo := db.InitTeamFlagRepo(tx)
	teamFlagL := make([]model.TeamFlag, 0)
	for _, contestFlag := range contestChallenge.ContestFlags {
		teamFlag, ok, msg := teamFlagRepo.GetWithConditions(db.GetOptions{
			{Key: "team_id", Value: team.ID, Op: "and"},
			{Key: "contest_flag_id", Value: contestFlag.ID, Op: "and"},
		})
		if ok {
			teamFlagL = append(teamFlagL, teamFlag)
			continue
		}
		options := db.CreateTeamFlagOptions{
			TeamID:        team.ID,
			ContestFlagID: contestFlag.ID,
			Solved:        false,
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
	submissions, _, ok, msg := submissionRepo.ListWithConditions(-1, -1, db.GetOptions{
		{Key: "team_id", Value: team.ID, Op: "and"},
		{Key: "contest_challenge_id", Value: contestChallenge.ID, Op: "and"},
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
		teamFlag, ok, msg := teamFlagRepo.GetWithConditions(db.GetOptions{
			{Key: "team_id", Value: team.ID, Op: "and"},
			{Key: "contest_flag_id", Value: contestFlag.ID, Op: "and"},
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
		if _, ok, _ := teamFlagRepo.GetWithConditions(db.GetOptions{
			{Key: "team_id", Value: team.ID, Op: "and"},
			{Key: "contest_flag_id", Value: contestFlag.ID, Op: "and"},
		}); !ok {
			return false
		}
	}
	return true
}
