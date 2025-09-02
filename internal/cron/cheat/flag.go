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
	questions, _, ok, _ := db.InitContestChallengeRepo(db.DB).List(-1, -1, db.GetOptions{
		Selects:    []string{"id", "type"},
		Conditions: map[string]any{"contest_id": contest.ID, "type": model.QuestionChallengeType},
	})
	if !ok {
		log.Logger.Warning("Failed to get questions challenge, CheckWrongFlag maybe wrong")
	}
	teams, _, ok, _ := db.InitTeamRepo(db.DB).List(-1, -1, db.GetOptions{
		Selects:    []string{"id"},
		Conditions: map[string]any{"contest_id": contest.ID},
		Preloads: map[string]db.GetOptions{
			"TeamFlags":   {Selects: []string{"id", "team_id", "value"}},
			"Submissions": {Selects: []string{"id", "team_id", "solved", "ip", "value", "contest_challenge_id", "created_at"}},
		},
	})
	if !ok {
		return
	}
	flagTeamIDMap := make(map[string][]uint)
	for _, team := range teams {
		for _, teamFlag := range team.TeamFlags {
			if _, ok = flagTeamIDMap[teamFlag.Value]; ok {
				if !slices.Contains(flagTeamIDMap[teamFlag.Value], teamFlag.TeamID) {
					flagTeamIDMap[teamFlag.Value] = append(flagTeamIDMap[teamFlag.Value], team.ID)
				}
			} else {
				flagTeamIDMap[teamFlag.Value] = []uint{team.ID}
			}
		}
	}
	cheatRepo := db.InitCheatRepo(db.DB)
	for _, team := range teams {
		for _, submission := range team.Submissions {
			if submission.Solved || slices.ContainsFunc(questions, func(q model.ContestChallenge) bool {
				return q.ID == submission.ContestChallengeID
			}) {
				continue
			}
			var tmp strings.Builder
			if teamIDL, ok := flagTeamIDMap[submission.Value]; ok {
				if !slices.Contains(flagTeamIDMap[submission.Value], team.ID) {
					for _, teamID := range teamIDL {
						tmp.WriteString(fmt.Sprintf("Team-%d, ", teamID))
					}
				}
			}
			cheatRepo.Create(db.CreateCheatOptions{
				TeamID:    sql.Null[uint]{V: team.ID, Valid: true},
				ContestID: sql.Null[uint]{V: contest.ID, Valid: true},
				IP:        submission.IP,
				Comment:   submission.Value,
				Reason:    fmt.Sprintf(model.SubmitOtherTeamFlag, team.ID, strings.Trim(tmp.String(), ", "), contest.ID),
				Type:      model.Cheater,
				Checked:   false,
				Time:      submission.CreatedAt,
			})
		}
	}
}
