package service

import (
	"context"
	"time"

	"CBCTF/internal/db"
	"CBCTF/internal/model"
	"CBCTF/internal/traffic"
)

type TrafficOverlapTeam struct {
	TeamID    uint      `json:"team_id"`
	FirstSeen time.Time `json:"first_seen"`
	VictimIDs []uint    `json:"victim_ids"`
}

type TrafficOverlap struct {
	IP             string               `json:"ip"`
	Classification string               `json:"classification"`
	Teams          []TrafficOverlapTeam `json:"teams"`
}

func GetContestTrafficOverlaps(ctx context.Context, contest model.Contest) ([]TrafficOverlap, model.RetVal) {
	rows, ret := db.InitTrafficRepo(db.DB.WithContext(ctx)).ListSharedContestVictimIPs(contest.ID, contest.Start, contest.Start.Add(contest.Duration))
	if !ret.OK {
		return nil, ret
	}
	return groupTrafficOverlaps(rows), model.SuccessRetVal()
}

func groupTrafficOverlaps(rows []db.TeamVictimIP) []TrafficOverlap {
	result := make([]TrafficOverlap, 0)
	index := make(map[string]int)
	for _, row := range rows {
		if !traffic.IsPublicTrafficIP(row.SrcIP) || checkWhitelistIP(row.SrcIP) {
			continue
		}
		i, ok := index[row.SrcIP]
		if !ok {
			i = len(result)
			index[row.SrcIP] = i
			result = append(result, TrafficOverlap{IP: row.SrcIP, Classification: "suspicious", Teams: []TrafficOverlapTeam{}})
		}
		result[i].Teams = append(result[i].Teams, TrafficOverlapTeam{TeamID: row.TeamID, FirstSeen: row.FirstTime, VictimIDs: row.VictimIDs})
	}
	return result
}
