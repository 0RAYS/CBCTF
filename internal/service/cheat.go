package service

import (
	"fmt"
	"maps"
	"net/netip"
	"slices"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/model"
	"CBCTF/internal/prometheus"
	"CBCTF/internal/traffic"
)

// cheatTimelineLimit 限制每条线索保存的观测事件数量。
// 事件按时间先后保留最早的若干条，既能定位共用行为的起点，也不会让线索记录无限膨胀。
const cheatTimelineLimit = 50

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
	maps.Copy(countOptions.Conditions, options.Conditions)
	countOptions.Conditions["checked"] = true
	checked, ret := db.InitCheatRepo(tx).Count(countOptions)
	if !ret.OK {
		return nil, 0, 0, ret
	}
	return cheats, count, checked, model.SuccessRetVal()
}

func UpdateCheat(tx *gorm.DB, cheat model.Cheat, form dto.UpdateCheatForm, reviewer uint) model.RetVal {
	options := db.UpdateCheatRepo{
		Reason:  form.Reason,
		Type:    form.Type,
		Checked: form.Checked,
		Comment: form.Comment,
	}
	// 结论必须能追溯到具体的管理员，否则复核结果无法追责。
	if form.Checked != nil {
		now := time.Now()
		options.ReviewedBy = &reviewer
		options.ReviewedAt = &now
	}
	return db.InitCheatRepo(tx).Update(cheat.ID, options)
}

func DeleteContestCheats(tx *gorm.DB, contest model.Contest) model.RetVal {
	return db.InitCheatRepo(tx).DeleteByContestID(contest.ID)
}

func DeleteCheat(tx *gorm.DB, cheat model.Cheat) model.RetVal {
	return db.InitCheatRepo(tx).Delete(cheat.ID)
}

func checkWhitelistIP(ip string) bool {
	ip = traffic.NormalizeTrafficIP(ip)
	addr, err := netip.ParseAddr(ip)
	return err != nil || slices.ContainsFunc(config.Env.Cheat.IP.Whitelist, func(cidr string) bool {
		if strings.Contains(cidr, "/") {
			prefix, err := netip.ParsePrefix(cidr)
			if err != nil {
				return false
			}
			return prefix.Contains(addr)
		}
		return traffic.NormalizeTrafficIP(cidr) == ip
	})
}

// shouldKeepUserGroup 判断共用同一 IP 的用户是否分属不同队伍。
// 同队成员共用网络是常态，跨队共用才是需要复核的线索。
func shouldKeepUserGroup(userTeamMap map[uint]uint, userIDs []uint) bool {
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

func CheckWebReqIP(tx *gorm.DB, contest model.Contest) model.RetVal {
	start, end := contest.Start, contest.Start.Add(contest.Duration)
	rows, ret := db.InitRequestRepo(tx).ListSharedContestUserIPs(contest.ID, start, end)
	if !ret.OK {
		return ret
	}

	byIP := make(map[string][]db.ContestSharedIPUser)
	for _, row := range rows {
		if checkWhitelistIP(row.IP) {
			continue
		}
		byIP[row.IP] = append(byIP[row.IP], row)
	}
	if len(byIP) == 0 {
		return model.SuccessRetVal()
	}

	ips := slices.Sorted(maps.Keys(byIP))
	requests, ret := db.InitRequestRepo(tx).ListContestIPRequests(contest.ID, ips, start, end, cheatTimelineLimit)
	if !ret.OK {
		return ret
	}
	submissions, ret := db.InitSubmissionRepo(tx).ListContestIPSubmissions(contest.ID, ips, start, end, cheatTimelineLimit)
	if !ret.OK {
		return ret
	}

	teamRepo := db.InitTeamRepo(tx)
	cheatRepo := db.InitCheatRepo(tx)
	batch := model.NewBatch(len(byIP))
	for _, ip := range ips {
		if tx.Statement.Context.Err() != nil {
			return batch.Result(tx.Statement.Context)
		}
		users := byIP[ip]
		if len(users) <= 1 {
			batch.Skip(ip, "single_user")
			continue
		}

		userIDs := make([]uint, 0, len(users))
		for _, user := range users {
			userIDs = append(userIDs, user.UserID)
		}
		userTeams, ret := teamRepo.GetUserTeamMap(contest.ID, userIDs...)
		if !ret.OK {
			batch.Fail(ip, "read_team_membership", ret)
			return batch.Result(tx.Statement.Context)
		}
		if !shouldKeepUserGroup(userTeams, userIDs) {
			batch.Skip(ip, "same_team")
			continue
		}

		evidence := buildSameWebIPEvidence(ip, users, userTeams, requests, submissions)
		labels := make([]string, 0, len(userIDs))
		for _, userID := range userIDs {
			labels = append(labels, strconv.FormatUint(uint64(userID), 10))
		}
		subject := fmt.Sprintf("User %s (%d requests, %d submissions, min gap %s)",
			strings.Join(labels, ","),
			users[0].TotalRequests,
			len(evidence.Web.Submissions),
			formatGap(evidence.Metrics.MinGapSeconds),
		)
		cheat := model.Cheat{
			ContestID:  contest.ID,
			Model:      map[string][]uint{model.Name(model.User{}): userIDs},
			IP:         ip,
			Comment:    ip,
			Reason:     fmt.Sprintf(string(model.ReqWebSameIPTmpl), subject),
			ReasonType: model.ReasonTypeSameWebIPType,
			Type:       model.SuspiciousType,
			Checked:    false,
			Time:       evidence.Metrics.FirstSeen,
			Evidence:   evidence,
			DedupKey:   ip,
		}
		cheat.LegacyHash = model.LegacyCheatHash(contest.ID, cheat.Model, cheat.ReasonType, cheat.IP, cheat.Comment)
		if !recordCheat(batch, ip, cheatRepo, cheat, contest.ID, model.ReasonTypeSameWebIPType) {
			return batch.Result(tx.Statement.Context)
		}
	}
	return batch.Result(tx.Statement.Context)
}

func CheckVictimReqIP(tx *gorm.DB, contest model.Contest) model.RetVal {
	start, end := contest.Start, contest.Start.Add(contest.Duration)
	rows, ret := db.InitTrafficRepo(tx).ListSharedContestVictimIPs(contest.ID, start, end)
	if !ret.OK {
		return ret
	}

	byIP := make(map[string][]db.TeamVictimIP)
	for _, row := range rows {
		if !traffic.IsPublicTrafficIP(row.SrcIP) || checkWhitelistIP(row.SrcIP) {
			continue
		}
		byIP[row.SrcIP] = append(byIP[row.SrcIP], row)
	}
	if len(byIP) == 0 {
		return model.SuccessRetVal()
	}

	ips := slices.Sorted(maps.Keys(byIP))
	accesses, ret := db.InitTrafficRepo(tx).ListContestVictimAccesses(contest.ID, ips, start, end, cheatTimelineLimit)
	if !ret.OK {
		return ret
	}

	cheatRepo := db.InitCheatRepo(tx)
	batch := model.NewBatch(len(byIP))
	for _, ip := range ips {
		if tx.Statement.Context.Err() != nil {
			return batch.Result(tx.Statement.Context)
		}
		teams := byIP[ip]
		if len(teams) <= 1 {
			batch.Skip(ip, "single_team")
			continue
		}

		evidence := buildSameVictimIPEvidence(ip, teams, accesses)
		labels := make([]string, 0, len(teams))
		teamIDs := make([]uint, 0, len(teams))
		for _, team := range teams {
			labels = append(labels, strconv.FormatUint(uint64(team.TeamID), 10))
			teamIDs = append(teamIDs, team.TeamID)
		}
		subject := fmt.Sprintf("Team %s (%d accesses, %d victims, min gap %s)",
			strings.Join(labels, ","),
			evidence.Metrics.Events,
			len(evidence.Victim.Victims),
			formatGap(evidence.Metrics.MinGapSeconds),
		)
		cheat := model.Cheat{
			ContestID:  contest.ID,
			Model:      map[string][]uint{model.Name(model.Team{}): teamIDs},
			IP:         ip,
			Comment:    ip,
			Reason:     fmt.Sprintf(string(model.ReqVictimSameIPTmpl), subject),
			ReasonType: model.ReasonTypeSameVictimIPType,
			Type:       model.SuspiciousType,
			Checked:    false,
			Time:       evidence.Metrics.FirstSeen,
			Evidence:   evidence,
			DedupKey:   ip,
		}
		cheat.LegacyHash = model.LegacyCheatHash(contest.ID, cheat.Model, cheat.ReasonType, cheat.IP, cheat.Comment)
		if !recordCheat(batch, ip, cheatRepo, cheat, contest.ID, model.ReasonTypeSameVictimIPType) {
			return batch.Result(tx.Statement.Context)
		}
	}
	return batch.Result(tx.Statement.Context)
}

// wrongFlagGroup 是同一支队伍针对同一题目重复提交他人 flag 的聚合。
// 一次提交可能同时命中多支队伍的 flag，因此提交记录按 ID 去重后保存。
type wrongFlagGroup struct {
	ContestChallengeID uint
	ChallengeName      string
	SubmittingTeamID   uint
	Value              string
	MatchedTeamIDs     []uint
	Submissions        []db.ContestWrongFlagSubmission
	seen               map[uint]bool
}

func CheckWrongFlag(tx *gorm.DB, contest model.Contest) model.RetVal {
	rows, ret := db.InitTeamFlagRepo(tx).ListContestWrongFlagSubmissions(contest.ID, contest.Start, contest.Start.Add(contest.Duration))
	if !ret.OK {
		return ret
	}

	groups := make(map[string]*wrongFlagGroup)
	for _, row := range rows {
		key := fmt.Sprintf("%d|%d|%s", row.ContestChallengeID, row.TeamID, row.Value)
		group, ok := groups[key]
		if !ok {
			group = &wrongFlagGroup{
				ContestChallengeID: row.ContestChallengeID,
				ChallengeName:      row.ChallengeName,
				SubmittingTeamID:   row.TeamID,
				Value:              row.Value,
				MatchedTeamIDs:     []uint{},
				seen:               make(map[uint]bool),
			}
			groups[key] = group
		}
		if !slices.Contains(group.MatchedTeamIDs, row.MatchedTeamID) {
			group.MatchedTeamIDs = append(group.MatchedTeamIDs, row.MatchedTeamID)
		}
		if group.seen[row.SubmissionID] {
			continue
		}
		group.seen[row.SubmissionID] = true
		group.Submissions = append(group.Submissions, row)
	}

	cheatRepo := db.InitCheatRepo(tx)
	batch := model.NewBatch(len(groups))
	for _, group := range groups {
		if tx.Statement.Context.Err() != nil {
			return batch.Result(tx.Statement.Context)
		}
		key := fmt.Sprintf("%d|%d|%s", group.ContestChallengeID, group.SubmittingTeamID, group.Value)
		if len(group.Submissions) == 0 {
			batch.Skip(key, "no_submission")
			continue
		}

		evidence := buildWrongFlagEvidence(group)
		teamIDs := append(slices.Clone(group.MatchedTeamIDs), group.SubmittingTeamID)
		var others strings.Builder
		for _, teamID := range group.MatchedTeamIDs {
			others.WriteString("Team-" + strconv.FormatUint(uint64(teamID), 10) + ", ")
		}
		challenge := group.ChallengeName
		if challenge == "" {
			challenge = "#" + strconv.FormatUint(uint64(group.ContestChallengeID), 10)
		}
		cheat := model.Cheat{
			ContestID: contest.ID,
			Model:     map[string][]uint{model.Name(model.Team{}): teamIDs},
			IP:        evidence.Flag.Submissions[0].IP,
			Comment:   group.Value,
			Reason: fmt.Sprintf(string(model.SubmitOtherTeamFlagTmpl), group.SubmittingTeamID,
				strings.Trim(others.String(), ", "), challenge, contest.ID),
			ReasonType: model.ReasonTypeWrongFlagType,
			Type:       model.CheaterType,
			Checked:    false,
			Time:       evidence.Flag.Submissions[0].At,
			Evidence:   evidence,
			DedupKey:   key,
		}
		cheat.LegacyHash = model.LegacyCheatHash(contest.ID, cheat.Model, cheat.ReasonType, cheat.IP, cheat.Comment)
		if !recordCheat(batch, key, cheatRepo, cheat, contest.ID, model.ReasonTypeWrongFlagType) {
			return batch.Result(tx.Statement.Context)
		}
	}
	return batch.Result(tx.Statement.Context)
}

// recordCheat 写入线索。已存在的线索就地刷新证据链，
// 使重复扫描只更新同一条记录，而不是不断堆积重复线索。
func recordCheat(batch *model.BatchResult, id string, repo *db.CheatRepo, cheat model.Cheat, contestID uint, reason model.CheatReasonType) bool {
	existing, created, ret := repo.Create(cheat)
	if !ret.OK {
		return recordCheatOutcome(batch, id, ret, false, contestID, reason)
	}
	if created {
		return recordCheatOutcome(batch, id, ret, true, contestID, reason)
	}
	if ret = repo.Update(existing.ID, db.RefreshCheatRepo{
		Hash:     cheat.Hash,
		Evidence: cheat.Evidence,
		Time:     cheat.Time,
	}); !ret.OK {
		return recordCheatOutcome(batch, id, ret, false, contestID, reason)
	}
	return recordCheatOutcome(batch, id, ret, false, contestID, reason)
}

// recordCheatOutcome 记录一条线索的处理结果，只有新落库的线索才计入检测计数。
func recordCheatOutcome(batch *model.BatchResult, id string, ret model.RetVal, created bool, contestID uint, reason model.CheatReasonType) bool {
	if !ret.OK {
		batch.Fail(id, "write_evidence", ret)
		return !model.BatchDependencyFailed(ret)
	}
	if !created {
		batch.Skip(id, "existing_evidence")
		return true
	}
	batch.Success(id, "created")
	prometheus.RecordCheatDetection(contestID, string(reason))
	return true
}

func RunCheatChecks(tx *gorm.DB, contest model.Contest) model.RetVal {
	checks := []struct {
		name string
		run  func(*gorm.DB, model.Contest) model.RetVal
	}{{"web_ip", CheckWebReqIP}, {"victim_ip", CheckVictimReqIP}, {"wrong_flag", CheckWrongFlag}}
	batch := model.NewBatch(len(checks))
	for _, check := range checks {
		if tx.Statement.Context.Err() != nil {
			return batch.Result(tx.Statement.Context)
		}
		ret := check.run(tx, contest)
		batch.Record(check.name, "scan", ret)
		if !ret.OK && model.BatchDependencyFailed(ret) {
			return batch.Result(tx.Statement.Context)
		}
	}
	return batch.Result(tx.Statement.Context)
}
