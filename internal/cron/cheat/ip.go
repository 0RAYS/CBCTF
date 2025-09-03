package cheat

import (
	"CBCTF/internal/db"
	"CBCTF/internal/model"
	"database/sql"
	"fmt"
	"net"
	"slices"
	"strings"
	"time"
)

// CheckVictimReqIP 检查用户访问靶机的 IP
func CheckVictimReqIP(contest model.Contest) {
	teams, _, ok, _ := db.InitTeamRepo(db.DB).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_id": contest.ID},
		Selects:    []string{"id"},
	})
	if !ok {
		return
	}
	type tmp struct {
		Time time.Time
		ID   uint
	}
	ipTeamMap := make(map[string][]tmp)
	victimRepo := db.InitVictimRepo(db.DB)
	for _, team := range teams {
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
					Reason:    fmt.Sprintf(model.ReqVictimSameIP, strings.Trim(str.String(), ", ")),
					Type:      model.Suspicious,
					Checked:   false,
					Time:      team.Time,
				})
			}
		}
	}
}
