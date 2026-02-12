package cheat

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/model"
	"fmt"
	"net"
	"slices"
	"strconv"
	"strings"
	"time"
)

func checkWhitelistIP(ip string) bool {
	addr := net.ParseIP(ip)
	return addr == nil || slices.ContainsFunc(config.Env.Cheat.IP.Whitelist, func(cidr string) bool {
		if strings.Contains(cidr, "/") {
			_, network, err := net.ParseCIDR(cidr)
			if err != nil {
				return false
			}
			return network.Contains(addr)
		}
		return cidr == ip
	})
}

// CheckWebReqIP 检查用户访问 Web 的 IP
func CheckWebReqIP(contest model.Contest) {
	userIDL, ret := db.GetUserIDByContestID(db.DB, contest.ID)
	if !ret.OK {
		return
	}
	repo := db.InitRequestRepo(db.DB)
	type ipUserInfo struct {
		UserID uint
		Time   time.Time
	}
	ipUserMap := make(map[string][]ipUserInfo)
	userIPL, ret := repo.GetUserIP(userIDL...)
	if !ret.OK {
		return
	}

	for _, result := range userIPL {
		if checkWhitelistIP(result.IP) {
			continue
		}
		if !slices.ContainsFunc(ipUserMap[result.IP], func(info ipUserInfo) bool {
			return info.UserID == result.UserID
		}) {
			ipUserMap[result.IP] = append(ipUserMap[result.IP], ipUserInfo{UserID: result.UserID, Time: result.FirstTime})
		}
	}

	cheatRepo := db.InitCheatRepo(db.DB)
	for ip, users := range ipUserMap {
		if len(users) > 1 {
			var str []string
			userIDs := make([]uint, 0, len(users))
			var earliest time.Time
			for i, user := range users {
				str = append(str, strconv.Itoa(int(user.UserID)))
				userIDs = append(userIDs, user.UserID)
				if i == 0 || user.Time.Before(earliest) {
					earliest = user.Time
				}
			}
			cheatRepo.Create(db.CreateCheatOptions{
				Model: model.CheatRefModel{
					model.User{}.ModelName(): userIDs,
					contest.ModelName():      {contest.ID},
				},
				IP:         ip,
				Comment:    ip,
				Reason:     fmt.Sprintf(model.ReqWebSameIP, fmt.Sprintf("User %s", strings.Join(str, ","))),
				ReasonType: model.ReasonTypeSameWebIP,
				Type:       model.Suspicious,
				Checked:    false,
				Time:       earliest,
			})
		}
	}
}

// CheckVictimReqIP 检查用户访问靶机的 IP
func CheckVictimReqIP(contest model.Contest) {
	teams, _, ret := db.InitTeamRepo(db.DB).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_id": contest.ID, "banned": false},
		Selects:    []string{"id"},
	})
	if !ret.OK {
		return
	}
	teamIDs := make([]uint, len(teams))
	for i, team := range teams {
		teamIDs[i] = team.ID
	}

	trafficResults, ret := db.InitTrafficRepo(db.DB).GetTeamVictimIP(teamIDs...)
	if !ret.OK {
		return
	}

	type tmp struct {
		Time time.Time
		ID   uint
	}
	ipTeamMap := make(map[string][]tmp)

	// 构建 IP -> teams 映射，同时过滤白名单
	for _, result := range trafficResults {
		if checkWhitelistIP(result.SrcIP) {
			continue
		}
		if !slices.ContainsFunc(ipTeamMap[result.SrcIP], func(s tmp) bool {
			return s.ID == result.TeamID
		}) {
			// 靶机流量的时间此处实际上为靶机关闭的时间, 但影响不大
			ipTeamMap[result.SrcIP] = append(ipTeamMap[result.SrcIP], tmp{Time: result.StopTime.Time, ID: result.TeamID})
		}
	}

	cheatRepo := db.InitCheatRepo(db.DB)
	for ip, v := range ipTeamMap {
		if len(v) > 1 {
			var str []string
			teamIDs := make([]uint, 0, len(v))
			var earliest time.Time
			for i, team := range v {
				str = append(str, strconv.Itoa(int(team.ID)))
				teamIDs = append(teamIDs, team.ID)
				if i == 0 || team.Time.Before(earliest) {
					earliest = team.Time
				}
			}
			cheatRepo.Create(db.CreateCheatOptions{
				Model: model.CheatRefModel{
					model.Team{}.ModelName(): teamIDs,
					contest.ModelName():      {contest.ID},
				},
				IP:         ip,
				Comment:    ip,
				Reason:     fmt.Sprintf(model.ReqVictimSameIP, fmt.Sprintf("Team %s", strings.Join(str, ","))),
				ReasonType: model.ReasonTypeSameVictimIP,
				Type:       model.Suspicious,
				Checked:    false,
				Time:       earliest,
			})
		}
	}
}
