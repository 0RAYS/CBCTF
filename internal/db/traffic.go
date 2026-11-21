package db

import (
	"fmt"
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
func (t *TrafficRepo) ReplaceAnalysis(victimID uint, ips []string, accesses []traffic.Access, report *traffic.AnalysisReport, archived bool) model.RetVal {
	updates := clause.AssignmentColumns([]string{"ips", "accesses", "analysis", "archived", "updated_at"})
	if report != nil && report.Partial {
		// A failed source cannot retract observations already persisted in the
		// authoritative database. Merge positive evidence until a complete scan.
		updates = clause.Assignments(map[string]any{
			"ips":      gorm.Expr("(SELECT COALESCE(jsonb_agg(DISTINCT value), '[]'::jsonb) FROM jsonb_array_elements(COALESCE(traffics.ips, '[]'::jsonb) || EXCLUDED.ips))"),
			"accesses": gorm.Expr("(SELECT COALESCE(jsonb_agg(DISTINCT value), '[]'::jsonb) FROM jsonb_array_elements(COALESCE(traffics.accesses, '[]'::jsonb) || EXCLUDED.accesses))"),
			"analysis": gorm.Expr("EXCLUDED.analysis"), "archived": gorm.Expr("EXCLUDED.archived"), "updated_at": gorm.Expr("EXCLUDED.updated_at"),
		})
	}
	res := t.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "victim_id"}},
		DoUpdates: updates,
		Where:     clause.Where{Exprs: []clause.Expression{clause.Expr{SQL: "NOT traffics.archived OR (EXCLUDED.archived AND traffics.analysis->>'partial' = 'true')"}}},
	}).Create(&model.Traffic{VictimID: victimID, IPs: ips, Accesses: accesses, Analysis: report, Archived: archived})
	if res.Error != nil {
		target := fmt.Sprintf("victim_id=%d", victimID)
		return model.FailRetVal(i18n.Model.Traffic.GetError, "Traffic", target, res.Error)
	}
	return model.SuccessRetVal()
}

func (t *TrafficRepo) GetAnalysis(victimID uint) (model.Traffic, model.RetVal) {
	var record model.Traffic
	target := fmt.Sprintf("victim_id=%d", victimID)
	if err := t.DB.Where("victim_id = ?", victimID).Limit(1).Find(&record).Error; err != nil {
		return record, model.FailRetVal(i18n.Model.Traffic.GetError, "Traffic", target, err)
	}
	return record, model.SuccessRetVal()
}

func (t *TrafficRepo) KnownFlags(victim model.Victim) ([]string, model.RetVal) {
	var values []string
	var err error
	var target string
	if victim.TeamID.Valid {
		target = fmt.Sprintf("challenge_id=%d,team_id=%d", victim.ChallengeID, victim.TeamID.V)
		err = t.DB.Model(&model.TeamFlag{}).
			Joins("JOIN challenge_flags ON challenge_flags.id = team_flags.challenge_flag_id AND challenge_flags.deleted_at IS NULL").
			Where("team_flags.team_id = ? AND challenge_flags.challenge_id = ?", victim.TeamID.V, victim.ChallengeID).
			Pluck("team_flags.value", &values).Error
	} else {
		target = fmt.Sprintf("challenge_id=%d", victim.ChallengeID)
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
		return nil, model.FailRetVal(i18n.Model.Traffic.GetError, "Traffic", target, err)
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
	target := fmt.Sprintf("victim_id=%d", victimID)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get traffic IPs: %s (target=%s)", res.Error, target)
		return nil, model.FailRetVal(i18n.Model.Traffic.GetError, "Traffic", target, res.Error)
	}
	return record.IPs, model.SuccessRetVal()
}

// TeamVictimIP 作弊检测查询结果。
// MinGapSeconds 与 MaxGapSeconds 是两支队伍访问之间的最小与最大间隔。
type TeamVictimIP struct {
	FirstTime     time.Time
	LastTime      time.Time
	AccessCount   int64
	SrcIP         string
	TeamID        uint
	VictimIDs     []uint `gorm:"serializer:json"`
	MinGapSeconds float64
	MaxGapSeconds float64
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
		),
		ordered AS (
			SELECT e.team_id, e.ip_val, e.observed_at,
				LAG(e.observed_at) OVER (PARTITION BY e.ip_val ORDER BY e.observed_at) AS prev_time,
				LAG(e.team_id) OVER (PARTITION BY e.ip_val ORDER BY e.observed_at) AS prev_team
			FROM expanded e
			INNER JOIN shared_ips s ON s.ip_val = e.ip_val
		),
		gaps AS (
			SELECT ip_val,
				MIN(EXTRACT(EPOCH FROM (observed_at - prev_time))) FILTER (WHERE prev_team IS DISTINCT FROM team_id) AS min_gap,
				MAX(EXTRACT(EPOCH FROM (observed_at - prev_time))) FILTER (WHERE prev_team IS DISTINCT FROM team_id) AS max_gap
			FROM ordered
			GROUP BY ip_val
		)
		SELECT
			e.team_id,
			e.ip_val AS src_ip,
			MIN(e.observed_at) AS first_time,
			MAX(e.observed_at) AS last_time,
			COUNT(*) AS access_count,
			jsonb_agg(DISTINCT e.victim_id) AS victim_ids,
			MIN(g.min_gap)::float8 AS min_gap_seconds,
			MAX(g.max_gap)::float8 AS max_gap_seconds
		FROM expanded e
		INNER JOIN shared_ips s ON s.ip_val = e.ip_val
		INNER JOIN gaps g ON g.ip_val = e.ip_val
		GROUP BY e.team_id, e.ip_val
		ORDER BY e.ip_val ASC, first_time ASC, e.team_id ASC
	`, contestID, start, end).Scan(&results)

	target := fmt.Sprintf("contest_id=%d", contestID)
	if res.Error != nil {
		log.Logger.Warningf("Failed to list shared victim IPs: %s (target=%s)", res.Error, target)
		return nil, model.FailRetVal(i18n.Model.Traffic.GetError, "Traffic", target, res.Error)
	}
	return results, model.SuccessRetVal()
}

// ContestVictimAccess 是共享 IP 命中靶机的一次访问观测。
type ContestVictimAccess struct {
	SrcIP         string
	TeamID        uint
	VictimID      uint
	ChallengeID   uint
	ChallengeName string
	ObservedAt    time.Time
	Source        string
	Capture       string
}

// ListContestVictimAccesses 返回给定 IP 命中靶机的访问明细，每个 IP 最多保留 limit 条最早记录。
func (t *TrafficRepo) ListContestVictimAccesses(contestID uint, ips []string, start, end time.Time, limit int) ([]ContestVictimAccess, model.RetVal) {
	if contestID == 0 || len(ips) == 0 {
		return nil, model.SuccessRetVal()
	}

	var results []ContestVictimAccess
	res := t.DB.Raw(`
		WITH expanded AS (
			SELECT
				victims.team_id,
				victims.id AS victim_id,
				COALESCE(victims.challenge_id, 0) AS challenge_id,
				access->>'ip' AS ip_val,
				(access->>'time')::timestamptz AS observed_at,
				access->>'source' AS source,
				access->>'capture' AS capture
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
				AND access->>'ip' IN ?
		),
		ranked AS (
			SELECT e.*, ROW_NUMBER() OVER (PARTITION BY e.ip_val ORDER BY e.observed_at ASC, e.victim_id ASC) AS event_rank
			FROM expanded e
		)
		SELECT
			r.ip_val AS src_ip,
			r.team_id,
			r.victim_id,
			r.challenge_id,
			r.observed_at,
			r.source,
			r.capture,
			challenges.name AS challenge_name
		FROM ranked r
		LEFT JOIN challenges ON challenges.id = r.challenge_id AND challenges.deleted_at IS NULL
		WHERE r.event_rank <= ?
		ORDER BY r.src_ip ASC, r.observed_at ASC
	`, contestID, start, end, ips, limit).Scan(&results)

	target := fmt.Sprintf("contest_id=%d,ips=%d", contestID, len(ips))
	if res.Error != nil {
		log.Logger.Warningf("Failed to list contest victim accesses: %s (target=%s)", res.Error, target)
		return nil, model.FailRetVal(i18n.Model.Traffic.GetError, "Traffic", target, res.Error)
	}
	return results, model.SuccessRetVal()
}
