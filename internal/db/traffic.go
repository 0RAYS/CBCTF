package db

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
