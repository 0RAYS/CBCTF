package cheat

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/model"
	"database/sql"
	"fmt"
	"net"
	"slices"
	"strconv"
	"strings"
	"time"
)

// CheckWebReqIP 检查用户访问 Web 的 IP
func CheckWebReqIP(contest model.Contest) {
	userIDL, ok, _ := db.GetUserIDByContestID(db.DB, contest.ID)
	if !ok {
		return
	}
	repo := db.InitRequestRepo(db.DB)
	ipUserMap := make(map[string][]uint)
	for _, userID := range userIDL {
		ipL, ok, _ := repo.GetUserIP(userID)
		if !ok {
			continue
		}
		for _, ip := range ipL {
			if addr := net.ParseIP(ip); addr == nil || slices.ContainsFunc(config.Env.Cheat.IP.Whitelist, func(cidr string) bool {
				if strings.Contains(cidr, "/") {
					_, network, err := net.ParseCIDR(cidr)
					if err != nil {
						return false
					}
					return network.Contains(addr)
				} else {
					return cidr == ip
				}
			}) {
				continue
			}
			if !slices.Contains(ipUserMap[ip], userID) {
				ipUserMap[ip] = append(ipUserMap[ip], userID)
			}
		}
	}
	cheatRepo := db.InitCheatRepo(db.DB)
	for ip, users := range ipUserMap {
		if len(users) > 1 {
			var str []string
			for _, user := range users {
				str = append(str, strconv.Itoa(int(user)))
			}
			for _, user := range users {
				first, ok, _ := repo.Get(db.GetOptions{
					Conditions: map[string]any{"id": user, "ip": ip},
					Selects:    []string{"id", "time"},
				})
				if !ok {
					continue
				}
				cheatRepo.Create(db.CreateCheatOptions{
					UserID:    sql.Null[uint]{V: user, Valid: true},
					ContestID: sql.Null[uint]{V: contest.ID, Valid: true},
					IP:        ip,
					Comment:   ip,
					Reason:    fmt.Sprintf(model.ReqWebSameIP, fmt.Sprintf("User %s", strings.Join(str, ","))),
					Type:      model.Suspicious,
					Checked:   false,
					Time:      first.Time,
				})
			}
		}
	}
}

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
	trafficRepo := db.InitTrafficRepo(db.DB)
	for _, team := range teams {
		victims, _, ok, _ := victimRepo.List(-1, -1, db.GetOptions{
			Conditions: map[string]any{"team_id": team.ID},
			Selects:    []string{"id", "team_id", "deleted_at"},
		})
		if !ok {
			continue
		}
		for _, victim := range victims {
			ipL, ok, _ := trafficRepo.GetVictimReqIP(victim.ID)
			if !ok {
				continue
			}
			for _, ip := range ipL {
				if addr := net.ParseIP(ip); addr == nil || slices.ContainsFunc(config.Env.Cheat.IP.Whitelist, func(cidr string) bool {
					if strings.Contains(cidr, "/") {
						_, network, err := net.ParseCIDR(cidr)
						if err != nil {
							return false
						}
						return network.Contains(addr)
					} else {
						return cidr == ip
					}
				}) {
					continue
				}
				if !slices.ContainsFunc(ipTeamMap[ip], func(s tmp) bool {
					return s.ID == team.ID
				}) {
					// 靶机流量的时间此处实际上为靶机关闭的时间, 但影响不大
					ipTeamMap[ip] = append(ipTeamMap[ip], tmp{Time: victim.DeletedAt.Time, ID: team.ID})
				}
			}
		}
	}
	cheatRepo := db.InitCheatRepo(db.DB)
	for ip, v := range ipTeamMap {
		if len(v) > 1 {
			var str []string
			for _, team := range teams {
				str = append(str, strconv.Itoa(int(team.ID)))
			}
			for _, team := range v {
				cheatRepo.Create(db.CreateCheatOptions{
					TeamID:    sql.Null[uint]{V: team.ID, Valid: true},
					ContestID: sql.Null[uint]{V: contest.ID, Valid: true},
					IP:        ip,
					Comment:   ip,
					Reason:    fmt.Sprintf(model.ReqVictimSameIP, fmt.Sprintf("Team %s", strings.Join(str, ","))),
					Type:      model.Suspicious,
					Checked:   false,
					Time:      team.Time,
				})
			}
		}
	}
}
