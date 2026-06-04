package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TrafficRepo struct {
	BaseRepo[model.Traffic]
}

func InitTrafficRepo(tx *gorm.DB) *TrafficRepo {
	return &TrafficRepo{
		BaseRepo: BaseRepo[model.Traffic]{
			DB: tx,
		},
	}
}

// UpsertIPs 写入或合并靶机的 IP 列表（ON CONFLICT 时合并去重）
func (t *TrafficRepo) UpsertIPs(victimID uint, ips []string) model.RetVal {
	if len(ips) == 0 {
		return model.SuccessRetVal()
	}
	record := model.Traffic{
		VictimID: victimID,
		IPs:      ips,
	}
	res := t.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "victim_id"}},
		DoUpdates: clause.Assignments(map[string]any{
			"ips":        gorm.Expr("(SELECT jsonb_agg(DISTINCT val) FROM jsonb_array_elements_text(traffics.ips || EXCLUDED.ips) AS val)"),
			"updated_at": gorm.Expr("NOW()"),
		}),
	}).Create(&record)
	if res.Error != nil {
		log.Logger.Warningf("Failed to upsert traffic IPs: victim_id=%d error=%s", victimID, res.Error)
		return model.RetVal{Msg: i18n.Model.CreateError, Attr: map[string]any{"Model": "Traffic", "Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
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
	TeamID    uint
	SrcIP     string
	FirstTime time.Time
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
				ip_val,
				victims.created_at
			FROM traffics
			CROSS JOIN LATERAL jsonb_array_elements_text(traffics.ips) AS ip_val
			INNER JOIN victims ON victims.id = traffics.victim_id
			INNER JOIN teams   ON teams.id   = victims.team_id AND teams.deleted_at IS NULL
			WHERE
				victims.contest_id     = ?
				AND victims.team_id    IS NOT NULL
				AND traffics.deleted_at IS NULL
				AND victims.created_at BETWEEN ? AND ?
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
			MIN(e.created_at) AS first_time
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
