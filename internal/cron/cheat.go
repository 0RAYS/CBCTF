package cron

import (
	"CBCTF/internal/db"
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
			if !contest.IsRunning() {
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
	ipTeamIDMap := make(map[string][]uint)
	requestRepo := db.InitRequestRepo(db.DB)
	victimRepo := db.InitVictimRepo(db.DB)
	for _, team := range teams {
		for _, user := range team.Users {
			for _, device := range user.Devices {
				ipL, ok, _ := requestRepo.GetIPByMagic(device.Magic)
				if !ok {
					continue
				}
				for _, ip := range ipL {
					netIP := net.ParseIP(ip)
					if netIP == nil {
						continue
					}
					if netIP.IsLoopback() {
						continue
					}
					if !slices.Contains(ipTeamIDMap[ip], team.ID) {
						ipTeamIDMap[ip] = append(ipTeamIDMap[ip], team.ID)
					}
				}
			}
		}
		victims, _, ok, _ := victimRepo.List(-1, -1, db.GetOptions{
			Selects:    []string{"id", "team_id"},
			Conditions: map[string]any{"team_id": team.ID},
			Deleted:    true,
			Preloads:   map[string]db.GetOptions{"Traffics": {Selects: []string{"id", "victim_id", "src_ip"}}},
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
				if !slices.Contains(ipTeamIDMap[traffics.SrcIP], victim.TeamID) {
					ipTeamIDMap[traffics.SrcIP] = append(ipTeamIDMap[traffics.SrcIP], victim.TeamID)
				}
			}
		}
	}
	cheatRepo := db.InitCheatRepo(db.DB)
	for ip, teamIDL := range ipTeamIDMap {
		if len(teamIDL) > 1 {
			var tmp strings.Builder
			for _, teamID := range teamIDL {
				tmp.WriteString(fmt.Sprintf("Team-%d, ", teamID))
			}
			for _, teamID := range teamIDL {
				cheatRepo.Create(db.CreateCheatOptions{
					TeamID:    sql.Null[uint]{V: teamID, Valid: true},
					ContestID: sql.Null[uint]{V: contest.ID, Valid: true},
					IP:        ip,
					Reason:    fmt.Sprintf(model.SameIP, strings.Trim(tmp.String(), ", "), ip),
					Type:      model.Suspicious,
					Checked:   false,
				})
			}
		}
	}
}

// checkWrongFlag 检查是否提交别队 flag
func checkWrongFlag(contest model.Contest) {
	teams, _, ok, _ := db.InitTeamRepo(db.DB).List(-1, -1, db.GetOptions{
		Selects:    []string{"id"},
		Conditions: map[string]any{"contest_id": contest.ID},
		Preloads: map[string]db.GetOptions{
			"TeamFlags":   {Selects: []string{"id", "team_id", "value"}},
			"Submissions": {Selects: []string{"id", "team_id", "solved", "ip", "value"}},
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
	for _, team := range teams {
		for _, submission := range team.Submissions {
			if submission.Solved {
				continue
			}
			if teamIDL, ok := flagTeamIDMap[submission.Value]; ok {
				if !slices.Contains(flagTeamIDMap[submission.Value], team.ID) {
					var tmp strings.Builder
					for _, teamID := range teamIDL {
						tmp.WriteString(fmt.Sprintf("Team-%d, ", teamID))
					}
					cheatRepo := db.InitCheatRepo(db.DB)
					cheatRepo.Create(db.CreateCheatOptions{
						TeamID:    sql.Null[uint]{V: team.ID, Valid: true},
						ContestID: sql.Null[uint]{V: contest.ID, Valid: true},
						IP:        submission.IP,
						Comment:   fmt.Sprintf("submission-%d", submission.ID),
						Reason:    fmt.Sprintf(model.SubmitOtherTeamFlag, team.ID, strings.Trim(tmp.String(), ", "), contest.ID),
						Type:      model.Cheater,
						Checked:   false,
					})
				}
			}
		}
	}
}
