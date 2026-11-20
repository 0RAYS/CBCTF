package service

import (
	"cmp"
	"fmt"
	"slices"
	"time"

	"CBCTF/internal/db"
	"CBCTF/internal/model"
)

// buildSameWebIPEvidence 组装共享网页 IP 的证据链：
// 谁在共用、在哪些题目上留下动作、动作发生在什么时候、彼此间隔多久。
func buildSameWebIPEvidence(ip string, users []db.ContestSharedIPUser, userTeams map[uint]uint, requests []db.ContestIPRequest, submissions []db.ContestIPSubmission) model.CheatEvidence {
	evidence := model.CheatEvidence{
		ReasonType: model.ReasonTypeSameWebIPType,
		Subject:    ip,
		Web:        &model.CheatWebEvidence{},
	}
	for _, user := range users {
		actor := model.CheatActor{
			Kind:      model.CheatActorKindUser,
			ID:        user.UserID,
			Label:     fmt.Sprintf("User-%d", user.UserID),
			FirstSeen: user.FirstTime,
			LastSeen:  user.LastTime,
			Count:     user.RequestCount,
		}
		if teamID, ok := userTeams[user.UserID]; ok && teamID != 0 {
			actor.TeamIDs = []uint{teamID}
		}
		evidence.Actors = append(evidence.Actors, actor)
	}

	for _, request := range requests {
		if request.IP != ip {
			continue
		}
		evidence.Web.Requests = append(evidence.Web.Requests, model.CheatRequest{
			At:        request.Time,
			UserID:    request.UserID,
			Method:    request.Method,
			Path:      request.Path,
			Status:    request.Status,
			LatencyMS: request.Latency.Milliseconds(),
			UserAgent: request.UserAgent,
		})
		evidence.Timeline = append(evidence.Timeline, model.CheatTimelineItem{
			At:        request.Time,
			Source:    model.CheatSourceWebRequest,
			ActorKind: model.CheatActorKindUser,
			ActorID:   request.UserID,
			Method:    request.Method,
			Path:      request.Path,
			Status:    request.Status,
			LatencyMS: request.Latency.Milliseconds(),
			UserAgent: request.UserAgent,
		})
	}

	challenges := make(map[uint]*model.CheatChallengeRef)
	for _, submission := range submissions {
		if submission.IP != ip {
			continue
		}
		evidence.Web.Submissions = append(evidence.Web.Submissions, model.CheatFlagSubmission{
			SubmissionID:       submission.SubmissionID,
			TeamID:             submission.TeamID,
			UserID:             submission.UserID,
			ContestChallengeID: submission.ContestChallengeID,
			ChallengeName:      submission.ChallengeName,
			IP:                 submission.IP,
			At:                 submission.SubmittedAt,
			Value:              submission.Value,
			Solved:             submission.Solved,
		})
		evidence.Timeline = append(evidence.Timeline, model.CheatTimelineItem{
			At:        submission.SubmittedAt,
			Source:    model.CheatSourceFlagSubmission,
			ActorKind: model.CheatActorKindUser,
			ActorID:   submission.UserID,
			Target:    submission.ChallengeName,
			TargetID:  submission.ContestChallengeID,
			Value:     submission.Value,
			Detail:    fmt.Sprintf("Team-%d", submission.TeamID),
		})
		ref, ok := challenges[submission.ContestChallengeID]
		if !ok {
			ref = &model.CheatChallengeRef{ContestChallengeID: submission.ContestChallengeID, Name: submission.ChallengeName}
			challenges[submission.ContestChallengeID] = ref
		}
		if !slices.Contains(ref.ActorIDs, submission.UserID) {
			ref.ActorIDs = append(ref.ActorIDs, submission.UserID)
		}
	}
	for _, ref := range challenges {
		evidence.Web.Challenges = append(evidence.Web.Challenges, *ref)
	}
	slices.SortFunc(evidence.Web.Challenges, func(a, b model.CheatChallengeRef) int {
		return cmp.Compare(a.ContestChallengeID, b.ContestChallengeID)
	})

	sortCheatTimeline(evidence.Timeline)

	totalRequests := int64(0)
	firstSeen, lastSeen := time.Time{}, time.Time{}
	minGap, maxGap := 0.0, 0.0
	for i, user := range users {
		if i == 0 {
			totalRequests = user.TotalRequests
			minGap, maxGap = user.MinGapSeconds, user.MaxGapSeconds
		}
		if firstSeen.IsZero() || user.FirstTime.Before(firstSeen) {
			firstSeen = user.FirstTime
		}
		if user.LastTime.After(lastSeen) {
			lastSeen = user.LastTime
		}
	}
	truncated := totalRequests > int64(len(evidence.Web.Requests)) || len(evidence.Web.Submissions) >= cheatTimelineLimit
	evidence.Metrics = cheatMetrics(totalRequests+int64(len(evidence.Web.Submissions)), len(evidence.Actors),
		firstSeen, lastSeen, minGap, maxGap, truncated)
	return evidence
}

// buildSameVictimIPEvidence 组装共享靶机 IP 的证据链：
// 命中了哪些队伍的靶机、对应哪些题目、访问发生在什么时候。
func buildSameVictimIPEvidence(ip string, teams []db.TeamVictimIP, accesses []db.ContestVictimAccess) model.CheatEvidence {
	evidence := model.CheatEvidence{
		ReasonType: model.ReasonTypeSameVictimIPType,
		Subject:    ip,
		Victim:     &model.CheatVictimEvidence{},
	}
	firstSeen, lastSeen := time.Time{}, time.Time{}
	minGap, maxGap := 0.0, 0.0
	totalAccesses := int64(0)
	for i, team := range teams {
		actor := model.CheatActor{
			Kind:      model.CheatActorKindTeam,
			ID:        team.TeamID,
			Label:     fmt.Sprintf("Team-%d", team.TeamID),
			FirstSeen: team.FirstTime,
			LastSeen:  team.LastTime,
			Count:     team.AccessCount,
		}
		evidence.Actors = append(evidence.Actors, actor)
		totalAccesses += team.AccessCount
		if i == 0 {
			minGap, maxGap = team.MinGapSeconds, team.MaxGapSeconds
		}
		if firstSeen.IsZero() || team.FirstTime.Before(firstSeen) {
			firstSeen = team.FirstTime
		}
		if team.LastTime.After(lastSeen) {
			lastSeen = team.LastTime
		}
	}

	victims := make(map[uint]*model.CheatVictimRef)
	for _, access := range accesses {
		if access.SrcIP != ip {
			continue
		}
		evidence.Victim.Accesses = append(evidence.Victim.Accesses, model.CheatVictimAccess{
			At:       access.ObservedAt,
			TeamID:   access.TeamID,
			VictimID: access.VictimID,
			Source:   access.Source,
			Capture:  access.Capture,
		})
		evidence.Timeline = append(evidence.Timeline, model.CheatTimelineItem{
			At:        access.ObservedAt,
			Source:    model.CheatSourceVictimAccess,
			ActorKind: model.CheatActorKindTeam,
			ActorID:   access.TeamID,
			Target:    access.ChallengeName,
			TargetID:  access.ContestChallengeID,
			Detail:    fmt.Sprintf("Victim-%d", access.VictimID),
		})
		ref, ok := victims[access.VictimID]
		if !ok {
			ref = &model.CheatVictimRef{
				VictimID:           access.VictimID,
				TeamID:             access.TeamID,
				ContestChallengeID: access.ContestChallengeID,
				ChallengeName:      access.ChallengeName,
				FirstSeen:          access.ObservedAt,
				LastSeen:           access.ObservedAt,
			}
			victims[access.VictimID] = ref
		}
		ref.AccessCount++
		if access.ObservedAt.Before(ref.FirstSeen) {
			ref.FirstSeen = access.ObservedAt
		}
		if access.ObservedAt.After(ref.LastSeen) {
			ref.LastSeen = access.ObservedAt
		}
	}
	for _, ref := range victims {
		evidence.Victim.Victims = append(evidence.Victim.Victims, *ref)
	}
	slices.SortFunc(evidence.Victim.Victims, func(a, b model.CheatVictimRef) int {
		return cmp.Compare(a.VictimID, b.VictimID)
	})

	sortCheatTimeline(evidence.Timeline)

	truncated := totalAccesses > int64(len(evidence.Victim.Accesses))
	evidence.Metrics = cheatMetrics(totalAccesses, len(evidence.Actors), firstSeen, lastSeen, minGap, maxGap, truncated)
	return evidence
}

// buildWrongFlagEvidence 组装误用他人 flag 的证据链：
// flag 的合法持有者与解题时刻、冒用者与提交时刻，以及涉及的题目。
func buildWrongFlagEvidence(group *wrongFlagGroup) model.CheatEvidence {
	evidence := model.CheatEvidence{
		ReasonType: model.ReasonTypeWrongFlagType,
		Subject:    group.Value,
		Flag: &model.CheatFlagEvidence{
			ContestChallengeID: group.ContestChallengeID,
			ChallengeName:      group.ChallengeName,
			SubmittingTeamID:   group.SubmittingTeamID,
			MatchedTeamIDs:     group.MatchedTeamIDs,
			MatchedValues:      []string{group.Value},
			Owners:             []model.CheatFlagOwner{},
		},
	}
	owners := make(map[uint]*model.CheatFlagOwner)
	for _, row := range group.Submissions {
		owner, ok := owners[row.MatchedTeamID]
		if !ok {
			owner = &model.CheatFlagOwner{
				TeamID: row.MatchedTeamID,
				UserID: row.MatchedUserID,
				Value:  group.Value,
			}
			owners[row.MatchedTeamID] = owner
		}
		if owner.ObtainedAt.IsZero() || (!row.MatchedSolvedAt.IsZero() && row.MatchedSolvedAt.Before(owner.ObtainedAt)) {
			owner.ObtainedAt = row.MatchedSolvedAt
		}
	}
	ownerIDs := make([]uint, 0, len(owners))
	for teamID := range owners {
		ownerIDs = append(ownerIDs, teamID)
	}
	slices.Sort(ownerIDs)
	for _, teamID := range ownerIDs {
		owner := owners[teamID]
		evidence.Flag.Owners = append(evidence.Flag.Owners, *owner)
		evidence.Actors = append(evidence.Actors, model.CheatActor{
			Kind:      model.CheatActorKindTeam,
			ID:        owner.TeamID,
			Label:     fmt.Sprintf("Team-%d", owner.TeamID),
			FirstSeen: owner.ObtainedAt,
			LastSeen:  owner.ObtainedAt,
			Count:     1,
		})
	}

	firstSeen, lastSeen := time.Time{}, time.Time{}
	for _, row := range group.Submissions {
		evidence.Flag.Submissions = append(evidence.Flag.Submissions, model.CheatFlagSubmission{
			SubmissionID:       row.SubmissionID,
			TeamID:             row.TeamID,
			UserID:             row.UserID,
			ContestChallengeID: row.ContestChallengeID,
			ChallengeName:      row.ChallengeName,
			IP:                 row.IP,
			At:                 row.CreatedAt,
			Value:              row.Value,
			Solved:             false,
		})
		evidence.Timeline = append(evidence.Timeline, model.CheatTimelineItem{
			At:        row.CreatedAt,
			Source:    model.CheatSourceFlagSubmission,
			ActorKind: model.CheatActorKindUser,
			ActorID:   row.UserID,
			Target:    row.ChallengeName,
			TargetID:  row.ContestChallengeID,
			Value:     row.Value,
			Detail:    fmt.Sprintf("Team-%d", row.TeamID),
		})
		if firstSeen.IsZero() || row.CreatedAt.Before(firstSeen) {
			firstSeen = row.CreatedAt
		}
		if row.CreatedAt.After(lastSeen) {
			lastSeen = row.CreatedAt
		}
	}
	for _, owner := range evidence.Flag.Owners {
		if owner.ObtainedAt.IsZero() {
			continue
		}
		evidence.Timeline = append(evidence.Timeline, model.CheatTimelineItem{
			At:        owner.ObtainedAt,
			Source:    model.CheatSourceFlagSolve,
			ActorKind: model.CheatActorKindTeam,
			ActorID:   owner.TeamID,
			Target:    group.ChallengeName,
			TargetID:  group.ContestChallengeID,
			Value:     owner.Value,
			Detail:    fmt.Sprintf("Team-%d", owner.TeamID),
		})
		if owner.ObtainedAt.Before(firstSeen) {
			firstSeen = owner.ObtainedAt
		}
		if owner.ObtainedAt.After(lastSeen) {
			lastSeen = owner.ObtainedAt
		}
	}

	for _, teamID := range group.MatchedTeamIDs {
		if _, ok := owners[teamID]; ok {
			continue
		}
		evidence.Actors = append(evidence.Actors, model.CheatActor{
			Kind:  model.CheatActorKindTeam,
			ID:    teamID,
			Label: fmt.Sprintf("Team-%d", teamID),
		})
	}
	evidence.Actors = append(evidence.Actors, model.CheatActor{
		Kind:      model.CheatActorKindTeam,
		ID:        group.SubmittingTeamID,
		Label:     fmt.Sprintf("Team-%d", group.SubmittingTeamID),
		FirstSeen: firstSeen,
		LastSeen:  lastSeen,
		Count:     int64(len(group.Submissions)),
	})

	sortCheatTimeline(evidence.Timeline)

	minGap, maxGap := cheatTimelineGap(evidence.Timeline)
	evidence.Metrics = cheatMetrics(int64(len(group.Submissions)), len(evidence.Actors),
		firstSeen, lastSeen, minGap, maxGap, false)
	return evidence
}

// cheatMetrics 汇总观测窗口的规模与时间跨度。
func cheatMetrics(events int64, actors int, firstSeen, lastSeen time.Time, minGap, maxGap float64, truncated bool) model.CheatMetrics {
	metrics := model.CheatMetrics{
		Events:        events,
		ActorCount:    actors,
		FirstSeen:     firstSeen,
		LastSeen:      lastSeen,
		MinGapSeconds: minGap,
		MaxGapSeconds: maxGap,
		Truncated:     truncated,
	}
	if !firstSeen.IsZero() && !lastSeen.IsZero() && lastSeen.After(firstSeen) {
		metrics.SpanSeconds = lastSeen.Sub(firstSeen).Seconds()
	}
	return metrics
}

// cheatTimelineGap 计算时间线上不同角色相邻事件之间的最小与最大间隔。
func cheatTimelineGap(timeline []model.CheatTimelineItem) (float64, float64) {
	minGap, maxGap := 0.0, 0.0
	for i := 1; i < len(timeline); i++ {
		if timeline[i].ActorKind == timeline[i-1].ActorKind && timeline[i].ActorID == timeline[i-1].ActorID {
			continue
		}
		gap := timeline[i].At.Sub(timeline[i-1].At).Seconds()
		if gap < 0 {
			continue
		}
		if minGap == 0 || gap < minGap {
			minGap = gap
		}
		if gap > maxGap {
			maxGap = gap
		}
	}
	return minGap, maxGap
}

// sortCheatTimeline 按观测时间排列证据事件。
func sortCheatTimeline(timeline []model.CheatTimelineItem) {
	slices.SortStableFunc(timeline, func(a, b model.CheatTimelineItem) int {
		return a.At.Compare(b.At)
	})
}

// formatGap 把秒数格式化成便于阅读的间隔。
func formatGap(seconds float64) string {
	if seconds <= 0 {
		return "0s"
	}
	duration := time.Duration(seconds * float64(time.Second))
	if duration >= time.Hour {
		return duration.Truncate(time.Minute).String()
	}
	if duration >= time.Minute {
		return duration.Truncate(time.Second).String()
	}
	return duration.Round(time.Millisecond).String()
}
