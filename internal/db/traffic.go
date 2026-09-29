package db

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/traffic"
)

type TrafficRepo struct {
	BaseRepo[model.Traffic]
}

// ReplaceAnalysis replaces a complete capture snapshot, including empty results.
func (t *TrafficRepo) ReplaceAnalysis(victimID uint, ips []string, accesses []traffic.TrafficAccess, report *traffic.AnalysisReport, archived bool) model.RetVal {
	res := t.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "victim_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"ips", "accesses", "analysis", "archived", "updated_at"}),
		Where:     clause.Where{Exprs: []clause.Expression{clause.Expr{SQL: "NOT traffics.archived"}}},
	}).Create(&model.Traffic{VictimID: victimID, IPs: ips, Accesses: accesses, Analysis: report, Archived: archived})
	if res.Error != nil {
		return model.RetVal{Msg: i18n.Model.Traffic.GetError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}

func (t *TrafficRepo) GetAnalysis(victimID uint) (model.Traffic, model.RetVal) {
	var record model.Traffic
	if err := t.DB.Where("victim_id = ?", victimID).Limit(1).Find(&record).Error; err != nil {
		return record, model.RetVal{Msg: i18n.Model.Traffic.GetError, Attr: map[string]any{"Error": err.Error()}}
	}
	return record, model.SuccessRetVal()
}

func (t *TrafficRepo) KnownFlags(victim model.Victim) ([]string, model.RetVal) {
	var values []string
	var err error
	if victim.TeamID.Valid && victim.ContestChallengeID.Valid {
		err = t.DB.Model(&model.TeamFlag{}).
			Joins("JOIN contest_flags ON contest_flags.id = team_flags.contest_flag_id AND contest_flags.deleted_at IS NULL").
			Where("team_flags.team_id = ? AND contest_flags.contest_challenge_id = ?", victim.TeamID.V, victim.ContestChallengeID.V).
			Pluck("team_flags.value", &values).Error
	} else {
		err = t.DB.Model(&model.ChallengeFlag{}).Where("challenge_id = ?", victim.ChallengeID).Pluck("value", &values).Error
		static := make([]string, 0, len(values))
		for _, value := range values {
			if match := model.StaticFlagTmpl.FindStringSubmatch(value); len(match) > 1 {
				static = append(static, match[1])
			} else if !model.DynamicFlagTmpl.MatchString(value) && !model.UUIDFlagTmpl.MatchString(value) {
				static = append(static, value)
			}
		}
		values = static
	}
	if err != nil {
		return nil, model.RetVal{Msg: i18n.Model.Traffic.GetError, Attr: map[string]any{"Error": err.Error()}}
	}
	return values, model.SuccessRetVal()
}

func InitTrafficRepo(tx *gorm.DB) *TrafficRepo {
	return &TrafficRepo{
		DB: tx,
	}
}

// GetVictimIPs 返回靶机的 IP 列表。
func (t *TrafficRepo) GetVictimIPs(victimID uint) ([]string, model.RetVal) {
	var record model.Traffic
	res := t.DB.Where("victim_id = ? AND deleted_at IS NULL", victimID).Limit(1).Find(&record)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get traffic IPs: victim_id=%d error=%s", victimID, res.Error)
		return nil, model.RetVal{Msg: i18n.Model.Traffic.GetError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return record.IPs, model.SuccessRetVal()
}

// TeamVictimIP 作弊检测查询结果
type TeamVictimIP struct {
	FirstTime time.Time
	SrcIP     string
	TeamID    uint
	VictimIDs []uint `gorm:"serializer:json"`
}

// ListSharedContestVictimIPs 返回同一比赛中出现在多支队伍靶机里的 IP。
func (t *TrafficRepo) ListSharedContestVictimIPs(contestID uint, start, end time.Time) ([]TeamVictimIP, model.RetVal) {
	if contestID == 0 {
		return nil, model.SuccessRetVal()
	}

	var results []TeamVictimIP
	res := t.DB.Raw(`
		WITH expanded AS (
			SELECT
				victims.team_id,
				access->>'ip' AS ip_val,
				(access->>'time')::timestamptz AS observed_at,
				victims.id AS victim_id
			FROM traffics
			CROSS JOIN LATERAL jsonb_array_elements(traffics.accesses) AS access
			INNER JOIN victims ON victims.id = traffics.victim_id
			INNER JOIN teams   ON teams.id   = victims.team_id AND teams.deleted_at IS NULL
			WHERE
				victims.contest_id     = ?
				AND victims.team_id    IS NOT NULL
				AND traffics.deleted_at IS NULL
				AND (access->>'time')::timestamptz >= ?
				AND (access->>'time')::timestamptz < ?
		),
		shared_ips AS (
			SELECT ip_val
			FROM expanded
			GROUP BY ip_val
			HAVING COUNT(DISTINCT team_id) > 1
		)
		SELECT
			e.team_id,
			e.ip_val AS src_ip,
			MIN(e.observed_at) AS first_time,
			jsonb_agg(DISTINCT e.victim_id) AS victim_ids
		FROM expanded e
		INNER JOIN shared_ips s ON s.ip_val = e.ip_val
		GROUP BY e.team_id, e.ip_val
		ORDER BY e.ip_val ASC, first_time ASC, e.team_id ASC
	`, contestID, start, end).Scan(&results)

	if res.Error != nil {
		log.Logger.Warningf("Failed to list shared victim IPs: %s", res.Error)
		return nil, model.RetVal{Msg: i18n.Model.Traffic.GetError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return results, model.SuccessRetVal()
}
