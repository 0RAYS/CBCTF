package repo

import (
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
	SrcPort  uint16
	DstPort  uint16
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
		SrcPort:  c.SrcPort,
		DstPort:  c.DstPort,
		Type:     c.Type,
		Subtype:  c.Subtype,
		Count:    c.Count,
		Size:     c.Size,
	}
}

type UpdateTrafficOptions struct {
}

func (u UpdateTrafficOptions) Convert2Map() map[string]any {
	return make(map[string]any)
}

func InitTrafficRepo(tx *gorm.DB) *TrafficRepo {
	return &TrafficRepo{
		BasicRepo: BasicRepo[model.Traffic]{
			DB: tx,
		},
	}
}
