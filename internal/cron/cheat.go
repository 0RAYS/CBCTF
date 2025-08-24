package cron

import (
	"CBCTF/internal/db"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"database/sql"
	"fmt"
	"net"
	"slices"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

func checkCheat(c *cron.Cron) {
	function := exec("CheckCheat", func() {
		contests, _, ok, _ := db.InitContestRepo(db.DB).List(-1, -1, db.GetOptions{
			Selects: []string{"id", "start", "duration"},
		})
		if !ok {
			return
		}
		for _, contest := range contests {
			if time.Now().Sub(contest.Start.Add(contest.Duration)) > 15*time.Minute {
				continue
			}
			checkRemoteIP(contest)
			checkWrongFlag(contest)
		}
	})
	function()
	c.Schedule(cron.Every(10*time.Minute), cron.FuncJob(function))
}

// checkRemoteIP 检查用户涉及的 IP, 包含访问站点和靶机访问
func checkRemoteIP(contest model.Contest) {
	teams, _, ok, _ := db.InitTeamRepo(db.DB).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_id": contest.ID},
		Selects:    []string{"id"},
		Preloads: map[string]db.GetOptions{
			"Users": {
				Selects:  []string{"id"},
				Preloads: map[string]db.GetOptions{"Devices": {Selects: []string{"id", "user_id", "magic"}}},
			},
		},
	})
	if !ok {
		return
	}
	type tmp struct {
		Time time.Time
		ID   uint
	}
	ipTeamMap := make(map[string][]tmp)
	requestRepo := db.InitRequestRepo(db.DB)
	victimRepo := db.InitVictimRepo(db.DB)
	for _, team := range teams {
		for _, user := range team.Users {
			for _, device := range user.Devices {
				requests, _, ok, _ := requestRepo.GetByMagic(device.Magic)
				if !ok {
					continue
				}
				for _, request := range requests {
					netIP := net.ParseIP(request.IP)
					if netIP == nil {
						continue
					}
					if netIP.IsLoopback() {
						continue
					}
					if !slices.ContainsFunc(ipTeamMap[request.IP], func(s tmp) bool {
						return s.ID == team.ID
					}) {
						ipTeamMap[request.IP] = append(ipTeamMap[request.IP], tmp{Time: request.Time, ID: team.ID})
					}
				}
			}
		}
		victims, _, ok, _ := victimRepo.List(-1, -1, db.GetOptions{
			Selects:    []string{"id", "team_id"},
			Conditions: map[string]any{"team_id": team.ID},
			Deleted:    true,
			Preloads:   map[string]db.GetOptions{"Traffics": {Selects: []string{"id", "victim_id", "src_ip", "created_at"}}},
		})
		if !ok {
			continue
		}
		for _, victim := range victims {
			for _, traffics := range victim.Traffics {
				netIP := net.ParseIP(traffics.SrcIP)
				if netIP == nil {
					continue
				}
				if netIP.IsLoopback() {
					continue
				}
				if !slices.ContainsFunc(ipTeamMap[traffics.SrcIP], func(s tmp) bool {
					return s.ID == victim.TeamID.V
				}) {
					// 靶机流量的时间此处实际上为靶机关闭的时间, 但影响不大
					ipTeamMap[traffics.SrcIP] = append(ipTeamMap[traffics.SrcIP], tmp{Time: traffics.CreatedAt, ID: victim.TeamID.V})
				}
			}
		}
	}
	cheatRepo := db.InitCheatRepo(db.DB)
	for ip, v := range ipTeamMap {
		if len(v) > 1 {
			var str strings.Builder
			for _, team := range v {
				str.WriteString(fmt.Sprintf("Team-%d, ", team.ID))
			}
			for _, team := range v {
				cheatRepo.Create(db.CreateCheatOptions{
					TeamID:    sql.Null[uint]{V: team.ID, Valid: true},
					ContestID: sql.Null[uint]{V: contest.ID, Valid: true},
					IP:        ip,
					Comment:   ip,
					Reason:    fmt.Sprintf(model.SameIP, strings.Trim(str.String(), ", ")),
					Type:      model.Suspicious,
					Checked:   false,
					Time:      team.Time,
				})
			}
		}
	}
}

// checkWrongFlag 检查是否提交别队 flag
func checkWrongFlag(contest model.Contest) {
	questions, _, ok, _ := db.InitContestChallengeRepo(db.DB).List(-1, -1, db.GetOptions{
		Selects:    []string{"id", "type"},
		Conditions: map[string]any{"contest_id": contest.ID, "type": model.QuestionChallengeType},
	})
	if !ok {
		log.Logger.Warning("Failed to get questions challenge, checkWrongFlag maybe wrong")
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
