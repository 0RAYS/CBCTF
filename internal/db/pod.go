package db

import (
	"fmt"

	"gorm.io/gorm"

	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
)

type PodRepo struct {
	BaseRepo[model.Pod]
}

func InitPodRepo(tx *gorm.DB) *PodRepo {
	return &PodRepo{
		DB: tx,
	}
}

// CreateBatch persists generated pod records in one insert.
func (p *PodRepo) CreateBatch(pods []model.Pod) ([]model.Pod, model.RetVal) {
	if len(pods) == 0 {
		return pods, model.SuccessRetVal()
	}
	if res := p.DB.Model(&model.Pod{}).Create(&pods); res.Error != nil {
		log.Logger.Warningf("Failed to create Pods: %s", res.Error)
		return nil, model.FailRetVal(i18n.Model.CreateError, model.Name(model.Pod{}), fmt.Sprintf("count=%d", len(pods)), res.Error)
	}
	return pods, model.SuccessRetVal()
}

func (p *PodRepo) DeleteByVictimID(victimIDL ...uint) model.RetVal {
	return p.DeleteByFieldID("victim_id", victimIDL...)
}
