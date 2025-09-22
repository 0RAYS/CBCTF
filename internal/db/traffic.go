package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"

	"gorm.io/gorm"
)

type TrafficRepo struct {
	BasicRepo[model.Traffic]
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
		BasicRepo: BasicRepo[model.Traffic]{
			DB: tx,
		},
	}
}

func (t *TrafficRepo) GetVictimReqIP(id uint) ([]string, bool, string) {
	var ipL []string
	res := t.DB.Model(&model.Traffic{}).Where("victim_id = ?", id).Distinct("src_ip").Find(&ipL)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Traffic: %s", res.Error)
		return nil, false, i18n.GetTrafficError
	}
	return ipL, true, i18n.Success
}
