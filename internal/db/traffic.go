package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"time"

	"gorm.io/gorm"
)

type TrafficRepo struct {
	BaseRepo[model.Traffic]
}

type CreateTrafficOptions struct {
	VictimID uint
	SrcIP    string
	DstIP    string
	Type     string
	Subtype  string
	Size     int
	Count    uint
}

func (c CreateTrafficOptions) Convert2Model() model.Model {
	return model.Traffic{
		VictimID: c.VictimID,
		SrcIP:    c.SrcIP,
		DstIP:    c.DstIP,
		Type:     c.Type,
		Subtype:  c.Subtype,
		Count:    c.Count,
		Size:     c.Size,
	}
}

func InitTrafficRepo(tx *gorm.DB) *TrafficRepo {
	return &TrafficRepo{
		BaseRepo: BaseRepo[model.Traffic]{
			DB: tx,
		},
	}
}

type TeamVictimIP struct {
	TeamID    uint
	SrcIP     string
	FirstTime time.Time
}

func (t *TrafficRepo) ListSharedContestVictimIPs(contestID uint, start, end time.Time) ([]TeamVictimIP, model.RetVal) {
	if contestID == 0 {
		return nil, model.SuccessRetVal()
	}

	sharedIPs := t.DB.Table("traffics").
		Select("traffics.src_ip").
		Joins("INNER JOIN victims ON victims.id = traffics.victim_id").
		Joins("INNER JOIN teams ON teams.id = victims.team_id AND teams.deleted_at IS NULL").
		Where("victims.contest_id = ? AND victims.team_id IS NOT NULL AND traffics.deleted_at IS NULL", contestID).
		Where("victims.created_at >= ? AND victims.created_at <= ?", start, end).
		Group("traffics.src_ip").
		Having("COUNT(DISTINCT victims.team_id) > 1")

	var teamVictimIPL []TeamVictimIP
	res := t.DB.Table("traffics").
		Select("victims.team_id, traffics.src_ip, MIN(victims.created_at) AS first_time").
		Joins("INNER JOIN victims ON victims.id = traffics.victim_id").
		Joins("INNER JOIN teams ON teams.id = victims.team_id AND teams.deleted_at IS NULL").
		Where("victims.contest_id = ? AND victims.team_id IS NOT NULL AND traffics.deleted_at IS NULL AND traffics.src_ip IN (?)", contestID, sharedIPs).
		Where("victims.created_at >= ? AND victims.created_at <= ?", start, end).
		Group("victims.team_id, traffics.src_ip").
		Order("traffics.src_ip ASC, first_time ASC, victims.team_id ASC").
		Scan(&teamVictimIPL)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Traffic: %s", res.Error)
		return nil, model.RetVal{Msg: i18n.Model.Traffic.GetError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return teamVictimIPL, model.SuccessRetVal()
}
