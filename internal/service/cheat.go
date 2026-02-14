package service

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"fmt"
	"net/netip"
	"slices"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

func CheckSameDevice(tx *gorm.DB, contest model.Contest) {
	userIDL, ret := db.GetUserIDByContestID(tx, contest.ID)
	if !ret.OK {
		return
	}
	deviceUserMap := make(map[string][]uint)
	devices, _, ret := db.InitDeviceRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"user_id": userIDL},
	})
	if !ret.OK {
		return
	}
	for _, device := range devices {
		if !slices.Contains(deviceUserMap[device.Magic], device.UserID) {
			deviceUserMap[device.Magic] = append(deviceUserMap[device.Magic], device.UserID)
		}
	}
	repo := db.InitCheatRepo(tx)
	for magic, userIDs := range deviceUserMap {
		if len(userIDs) > 1 {
			var str []string
			for _, uid := range userIDs {
				str = append(str, strconv.Itoa(int(uid)))
			}
			repo.Create(db.CreateCheatOptions{
				ContestID:  contest.ID,
				Model:      model.CheatRefModel{model.User{}.ModelName(): userIDs},
				Magic:      magic,
				Reason:     fmt.Sprintf(model.SameDeviceMagic, fmt.Sprintf("User %s", strings.Join(str, ","))),
				ReasonType: model.ReasonTypeSameDevice,
				Type:       model.Suspicious,
				Time:       devices[0].CreatedAt,
			})
		}
	}
}

// CheckWrongFlag 检查是否提交别队 flag
func CheckWrongFlag(tx *gorm.DB, contest model.Contest) {
	questions, _, ret := db.InitContestChallengeRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_id": contest.ID, "type": model.QuestionChallengeType},
	})
	if !ret.OK {
		log.Logger.Warning("Failed to get questions challenge, CheckWrongFlag maybe wrong")
	}
	teams, _, ret := db.InitTeamRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_id": contest.ID},
	})
	if !ret.OK {
		return
	}
	teamIDs := make([]uint, len(teams))
	for i, team := range teams {
		teamIDs[i] = team.ID
	}
	teamFlags, _, ret := db.InitTeamFlagRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"team_id": teamIDs},
	})
	if !ret.OK {
		return
	}
	submissions, _, ret := db.InitSubmissionRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"team_id": teamIDs},
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
	cheatRepo := db.InitCheatRepo(tx)
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
				ContestID:  contest.ID,
				Model:      model.CheatRefModel{model.Team{}.ModelName(): append(teamIDL, submission.TeamID)},
				IP:         submission.IP,
				Comment:    submission.Value,
				Reason:     fmt.Sprintf(model.SubmitOtherTeamFlag, submission.TeamID, strings.Trim(tmp.String(), ", "), contest.ID),
				ReasonType: model.ReasonTypeWrongFlag,
				Type:       model.Cheater,
				Checked:    false,
				Time:       submission.CreatedAt,
			})
		}
	}
}

func checkWhitelistIP(ip string) bool {
	addr, err := netip.ParseAddr(ip)
	return err != nil || slices.ContainsFunc(config.Env.Cheat.IP.Whitelist, func(cidr string) bool {
		if strings.Contains(cidr, "/") {
			prefix, err := netip.ParsePrefix(cidr)
			if err != nil {
				return false
			}
			return prefix.Contains(addr)
		}
		return cidr == ip
	})
}

// CheckWebReqIP 检查用户访问 Web 的 IP
func CheckWebReqIP(tx *gorm.DB, contest model.Contest) {
	userIDL, ret := db.GetUserIDByContestID(tx, contest.ID)
	if !ret.OK {
		return
	}
	repo := db.InitRequestRepo(tx)
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

	cheatRepo := db.InitCheatRepo(tx)
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
				ContestID:  contest.ID,
				Model:      model.CheatRefModel{model.User{}.ModelName(): userIDs},
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
func CheckVictimReqIP(tx *gorm.DB, contest model.Contest) {
	teams, _, ret := db.InitTeamRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_id": contest.ID, "banned": false},
	})
	if !ret.OK {
		return
	}
	teamIDs := make([]uint, len(teams))
	for i, team := range teams {
		teamIDs[i] = team.ID
	}

	trafficResults, ret := db.InitTrafficRepo(tx).GetTeamVictimIP(teamIDs...)
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

	cheatRepo := db.InitCheatRepo(tx)
	for ip, v := range ipTeamMap {
		if len(v) > 1 {
			var str []string
			teamIDs = make([]uint, 0, len(v))
			var earliest time.Time
			for i, team := range v {
				str = append(str, strconv.Itoa(int(team.ID)))
				teamIDs = append(teamIDs, team.ID)
				if i == 0 || team.Time.Before(earliest) {
					earliest = team.Time
				}
			}
			cheatRepo.Create(db.CreateCheatOptions{
				ContestID:  contest.ID,
				Model:      model.CheatRefModel{model.Team{}.ModelName(): teamIDs},
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
