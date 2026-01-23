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

func (t *TrafficRepo) GetVictimReqIP(id uint) ([]string, model.RetVal) {
	var ipL []string
	res := t.DB.Model(&model.Traffic{}).Where("victim_id = ?", id).Distinct("src_ip").Find(&ipL)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Traffic: %s", res.Error)
		return nil, model.RetVal{Msg: i18n.Model.GetError, Attr: map[string]any{"Model": model.Traffic{}.GetModelName(), "Error": res.Error.Error()}}
	}
	return ipL, model.SuccessRetVal()
}
