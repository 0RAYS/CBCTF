package db

import (
	"CBCTF/internal/model"

	"gorm.io/gorm"
)

type PodRepo struct {
	BaseRepo[model.Pod]
}

type CreatePodOptions struct {
	VictimID uint
	Name     string
	Spec     model.PodSpec
}

func (c CreatePodOptions) Convert2Model() model.Model {
	return model.Pod{
		VictimID: c.VictimID,
		Name:     c.Name,
		Spec:     c.Spec,
	}
}

func InitPodRepo(tx *gorm.DB) *PodRepo {
	return &PodRepo{
		BaseRepo: BaseRepo[model.Pod]{
			DB: tx,
		},
	}
}

func (p *PodRepo) DeleteByVictimID(victimIDL ...uint) model.RetVal {
	return p.DeleteByFieldID("victim_id", victimIDL...)
}
