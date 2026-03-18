package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"

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
	TeamID   uint
	VictimID uint
	SrcIP    string
	StopTime gorm.DeletedAt
}

func (t *TrafficRepo) GetTeamVictimIP(teamIDL ...uint) ([]TeamVictimIP, model.RetVal) {
	if len(teamIDL) == 0 {
		return nil, model.SuccessRetVal()
	}
	var teamVictimIPL []TeamVictimIP
	res := t.DB.Table("traffics").
		Select("victims.team_id, victims.id AS victim_id, traffics.src_ip, victims.deleted_at AS stop_time").
		Joins("INNER JOIN victims ON traffics.victim_id = victims.id").
		Where("victims.team_id IN ? AND traffics.deleted_at IS NULL", teamIDL).
		Group("victims.team_id, victims.id, traffics.src_ip, victims.deleted_at").
		Scan(&teamVictimIPL)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Traffic: %s", res.Error)
		return nil, model.RetVal{Msg: i18n.Model.Traffic.GetError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return teamVictimIPL, model.SuccessRetVal()
}
