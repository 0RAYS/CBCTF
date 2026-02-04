package cheat

import (
	"CBCTF/internal/db"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"database/sql"
	"fmt"
	"slices"
	"strings"
)

// CheckWrongFlag 检查是否提交别队 flag
func CheckWrongFlag(contest model.Contest) {
	questions, _, ret := db.InitContestChallengeRepo(db.DB).List(-1, -1, db.GetOptions{
		Selects:    []string{"id", "type"},
		Conditions: map[string]any{"contest_id": contest.ID, "type": model.QuestionChallengeType},
	})
	if !ret.OK {
		log.Logger.Warning("Failed to get questions challenge, CheckWrongFlag maybe wrong")
	}
	teams, _, ret := db.InitTeamRepo(db.DB).List(-1, -1, db.GetOptions{
		Selects:    []string{"id"},
		Conditions: map[string]any{"contest_id": contest.ID},
	})
	if !ret.OK {
		return
	}
	teamIDs := make([]uint, len(teams))
	for i, team := range teams {
		teamIDs[i] = team.ID
	}
	teamFlags, _, ret := db.InitTeamFlagRepo(db.DB).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"team_id": teamIDs},
		Selects:    []string{"id", "team_id", "value"},
	})
	if !ret.OK {
		return
	}
	submissions, _, ret := db.InitSubmissionRepo(db.DB).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"team_id": teamIDs},
		Selects:    []string{"id", "team_id", "solved", "ip", "value", "contest_challenge_id", "created_at"},
	})
	if !ret.OK {
		return
	}
	flagTeamIDMap := make(map[string][]uint)
	for _, teamFlag := range teamFlags {
		if _, ok := flagTeamIDMap[teamFlag.Value]; ok {
			if !slices.Contains(flagTeamIDMap[teamFlag.Value], teamFlag.TeamID) {
				flagTeamIDMap[teamFlag.Value] = append(flagTeamIDMap[teamFlag.Value], teamFlag.TeamID)
			}
		} else {
			flagTeamIDMap[teamFlag.Value] = []uint{teamFlag.TeamID}
		}
	}
	cheatRepo := db.InitCheatRepo(db.DB)
	for _, submission := range submissions {
		if submission.Solved || slices.ContainsFunc(questions, func(q model.ContestChallenge) bool {
			return q.ID == submission.ContestChallengeID
		}) {
			continue
		}

		// 检查是否提交了其他队伍的 flag
		if teamIDL, ok := flagTeamIDMap[submission.Value]; ok && !slices.Contains(teamIDL, submission.TeamID) {
			var tmp strings.Builder
			for _, teamID := range teamIDL {
				tmp.WriteString(fmt.Sprintf("Team-%d, ", teamID))
			}
			cheatRepo.Create(db.CreateCheatOptions{
				TeamID:    sql.Null[uint]{V: submission.TeamID, Valid: true},
				ContestID: sql.Null[uint]{V: contest.ID, Valid: true},
				IP:        submission.IP,
				Comment:   submission.Value,
				Reason:    fmt.Sprintf(model.SubmitOtherTeamFlag, submission.TeamID, strings.Trim(tmp.String(), ", "), contest.ID),
				Type:      model.Cheater,
				Checked:   false,
				Time:      submission.CreatedAt,
			})
		}
	}
}
