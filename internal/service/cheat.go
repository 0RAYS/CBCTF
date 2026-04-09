package service

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/model"
	"CBCTF/internal/prometheus"
	"fmt"
	"net/netip"
	"slices"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

func ListCheats(tx *gorm.DB, contest model.Contest, form dto.GetCheatsForm) ([]model.Cheat, int64, int64, model.RetVal) {
	options := db.GetOptions{
		Conditions: map[string]any{"contest_id": contest.ID},
		Sort:       []string{"id DESC"},
	}
	if form.Type != "" {
		options.Conditions["type"] = form.Type
	}
	if form.ReasonType != "" {
		options.Conditions["reason_type"] = form.ReasonType
	}
	cheats, count, ret := db.InitCheatRepo(tx).List(form.Limit, form.Offset, options)
	if !ret.OK {
		return nil, 0, 0, ret
	}
	countOptions := db.CountOptions{Conditions: map[string]any{}}
	for key, value := range options.Conditions {
		countOptions.Conditions[key] = value
	}
	countOptions.Conditions["checked"] = true
	checked, ret := db.InitCheatRepo(tx).Count(countOptions)
	if !ret.OK {
		return nil, 0, 0, ret
	}
	return cheats, count, checked, model.SuccessRetVal()
}

func UpdateCheat(tx *gorm.DB, cheat model.Cheat, form dto.UpdateCheatForm) model.RetVal {
	return db.InitCheatRepo(tx).Update(cheat.ID, db.UpdateCheatRepo{
		Reason:  form.Reason,
		Type:    form.Type,
		Checked: form.Checked,
		Comment: form.Comment,
	})
}

func DeleteContestCheats(tx *gorm.DB, contest model.Contest) model.RetVal {
	return db.InitCheatRepo(tx).DeleteByContestID(contest.ID)
}

func DeleteCheat(tx *gorm.DB, cheat model.Cheat) model.RetVal {
	return db.InitCheatRepo(tx).Delete(cheat.ID)
}

type deviceInfo struct {
	UserID    uint
	FirstTime time.Time
}

type ipUserInfo struct {
	UserID uint
	Time   time.Time
}

func shouldKeepUserGroup(contestID uint, userIDs []uint, teamRepo *db.TeamRepo) bool {
	userTeamMap, ret := teamRepo.GetUserTeamMap(contestID, userIDs...)
	if !ret.OK {
		return false
	}

	teamSet := make(map[uint]struct{})
	missingTeam := false
	for _, userID := range userIDs {
		teamID, ok := userTeamMap[userID]
		if !ok || teamID == 0 {
			missingTeam = true
			continue
		}
		teamSet[teamID] = struct{}{}
	}
	if len(teamSet) == 0 {
		return missingTeam && len(userIDs) > 1
	}
	return len(teamSet) > 1 || missingTeam
}

func CheckSameDevice(tx *gorm.DB, contest model.Contest) {
	rows, ret := db.InitDeviceRepo(tx).ListSharedContestDevices(contest.ID, contest.Start, contest.Start.Add(contest.Duration))
	if !ret.OK {
		return
	}

	deviceUserMap := make(map[string][]deviceInfo)
	for _, row := range rows {
		deviceUserMap[row.Magic] = append(deviceUserMap[row.Magic], deviceInfo{
			UserID:    row.UserID,
			FirstTime: row.FirstTime,
		})
	}

	teamRepo := db.InitTeamRepo(tx)
	repo := db.InitCheatRepo(tx)
	for magic, infos := range deviceUserMap {
		if len(infos) <= 1 {
			continue
		}

		userIDs := make([]uint, 0, len(infos))
		for _, info := range infos {
			userIDs = append(userIDs, info.UserID)
		}
		if !shouldKeepUserGroup(contest.ID, userIDs, teamRepo) {
			continue
		}

		var (
			users    []string
			earliest time.Time
		)
		for i, info := range infos {
			users = append(users, strconv.Itoa(int(info.UserID)))
			if i == 0 || info.FirstTime.Before(earliest) {
				earliest = info.FirstTime
			}
		}

		repo.Create(db.CreateCheatOptions{
			ContestID:  contest.ID,
			Model:      model.CheatRefModel{model.ModelName(model.User{}): userIDs},
			Magic:      magic,
			Reason:     fmt.Sprintf(string(model.SameDeviceMagicTmpl), fmt.Sprintf("User %s", strings.Join(users, ","))),
			ReasonType: model.ReasonTypeSameDeviceType,
			Type:       model.SuspiciousType,
			Time:       earliest,
		})
		prometheus.RecordCheatDetection(string(model.ReasonTypeSameDeviceType))
	}
}

func CheckWrongFlag(tx *gorm.DB, contest model.Contest) {
	rows, ret := db.InitTeamFlagRepo(tx).ListContestWrongFlagSubmissions(contest.ID, contest.Start, contest.Start.Add(contest.Duration))
	if !ret.OK {
		return
	}

	type submissionDetail struct {
		TeamID    uint
		IP        string
		Value     string
		CreatedAt time.Time
		Others    []uint
	}

	submissionMap := make(map[uint]*submissionDetail)
	for _, row := range rows {
		detail, ok := submissionMap[row.SubmissionID]
		if !ok {
			detail = &submissionDetail{
				TeamID:    row.TeamID,
				IP:        row.IP,
				Value:     row.Value,
				CreatedAt: row.CreatedAt,
				Others:    make([]uint, 0),
			}
			submissionMap[row.SubmissionID] = detail
		}
		if !slices.Contains(detail.Others, row.MatchedTeamID) {
			detail.Others = append(detail.Others, row.MatchedTeamID)
		}
	}

	cheatRepo := db.InitCheatRepo(tx)
	for _, submission := range submissionMap {
		if len(submission.Others) == 0 {
			continue
		}

		var tmp strings.Builder
		teamIDs := make([]uint, 0, len(submission.Others)+1)
		for _, teamID := range submission.Others {
			tmp.WriteString(fmt.Sprintf("Team-%d, ", teamID))
			teamIDs = append(teamIDs, teamID)
		}
		teamIDs = append(teamIDs, submission.TeamID)

		cheatRepo.Create(db.CreateCheatOptions{
			ContestID:  contest.ID,
			Model:      model.CheatRefModel{model.ModelName(model.Team{}): teamIDs},
			IP:         submission.IP,
			Comment:    submission.Value,
			Reason:     fmt.Sprintf(string(model.SubmitOtherTeamFlagTmpl), submission.TeamID, strings.Trim(tmp.String(), ", "), contest.ID),
			ReasonType: model.ReasonTypeWrongFlagType,
			Type:       model.CheaterType,
			Checked:    false,
			Time:       submission.CreatedAt,
		})
		prometheus.RecordCheatDetection(string(model.ReasonTypeWrongFlagType))
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

func CheckWebReqIP(tx *gorm.DB, contest model.Contest) {
	rows, ret := db.InitRequestRepo(tx).ListSharedContestUserIPs(contest.ID, contest.Start, contest.Start.Add(contest.Duration))
	if !ret.OK {
		return
	}

	ipUserMap := make(map[string][]ipUserInfo)
	for _, row := range rows {
		if checkWhitelistIP(row.IP) {
			continue
		}
		ipUserMap[row.IP] = append(ipUserMap[row.IP], ipUserInfo{
			UserID: row.UserID,
			Time:   row.FirstTime,
		})
	}

	teamRepo := db.InitTeamRepo(tx)
	cheatRepo := db.InitCheatRepo(tx)
	for ip, users := range ipUserMap {
		if len(users) <= 1 {
			continue
		}

		userIDs := make([]uint, 0, len(users))
		for _, user := range users {
			userIDs = append(userIDs, user.UserID)
		}
		if !shouldKeepUserGroup(contest.ID, userIDs, teamRepo) {
			continue
		}

		var (
			str      []string
			earliest time.Time
		)
		for i, user := range users {
			str = append(str, strconv.Itoa(int(user.UserID)))
			if i == 0 || user.Time.Before(earliest) {
				earliest = user.Time
			}
		}

		cheatRepo.Create(db.CreateCheatOptions{
			ContestID:  contest.ID,
			Model:      model.CheatRefModel{model.ModelName(model.User{}): userIDs},
			IP:         ip,
			Comment:    ip,
			Reason:     fmt.Sprintf(string(model.ReqWebSameIPTmpl), fmt.Sprintf("User %s", strings.Join(str, ","))),
			ReasonType: model.ReasonTypeSameWebIPType,
			Type:       model.SuspiciousType,
			Checked:    false,
			Time:       earliest,
		})
		prometheus.RecordCheatDetection(string(model.ReasonTypeSameWebIPType))
	}
}

func CheckVictimReqIP(tx *gorm.DB, contest model.Contest) {
	type teamInfo struct {
		Time time.Time
		ID   uint
	}

	rows, ret := db.InitTrafficRepo(tx).ListSharedContestVictimIPs(contest.ID, contest.Start, contest.Start.Add(contest.Duration))
	if !ret.OK {
		return
	}

	ipTeamMap := make(map[string][]teamInfo)
	for _, row := range rows {
		if checkWhitelistIP(row.SrcIP) {
			continue
		}
		ipTeamMap[row.SrcIP] = append(ipTeamMap[row.SrcIP], teamInfo{
			Time: row.FirstTime,
			ID:   row.TeamID,
		})
	}

	cheatRepo := db.InitCheatRepo(tx)
	for ip, teams := range ipTeamMap {
		if len(teams) <= 1 {
			continue
		}

		var (
			str      []string
			teamIDs  []uint
			earliest time.Time
		)
		for i, team := range teams {
			str = append(str, strconv.Itoa(int(team.ID)))
			teamIDs = append(teamIDs, team.ID)
			if i == 0 || team.Time.Before(earliest) {
				earliest = team.Time
			}
		}

		cheatRepo.Create(db.CreateCheatOptions{
			ContestID:  contest.ID,
			Model:      model.CheatRefModel{model.ModelName(model.Team{}): teamIDs},
			IP:         ip,
			Comment:    ip,
			Reason:     fmt.Sprintf(string(model.ReqVictimSameIPTmpl), fmt.Sprintf("Team %s", strings.Join(str, ","))),
			ReasonType: model.ReasonTypeSameVictimIPType,
			Type:       model.SuspiciousType,
			Checked:    false,
			Time:       earliest,
		})
		prometheus.RecordCheatDetection(string(model.ReasonTypeSameVictimIPType))
	}
}
