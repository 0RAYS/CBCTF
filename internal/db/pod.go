package db

import (
	"CBCTF/internal/model"

	"gorm.io/gorm"
)

type PodRepo struct {
	BaseRepo[model.Pod]
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
