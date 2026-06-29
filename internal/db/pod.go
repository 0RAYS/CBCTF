package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
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

// Create skips generic uniqueness preflight queries. Pods use a database
// generated primary key and have no natural unique key to validate here.
func (p *PodRepo) Create(pod model.Pod) (model.Pod, model.RetVal) {
	if res := p.DB.Model(&model.Pod{}).Create(&pod); res.Error != nil {
		log.Logger.Warningf("Failed to create Pod: %s", res.Error)
		return model.Pod{}, model.RetVal{Msg: i18n.Model.CreateError, Attr: map[string]any{"Model": model.Name(model.Pod{}), "Error": res.Error.Error()}}
	}
	return pod, model.SuccessRetVal()
}

// CreateBatch persists generated pod records in one insert.
func (p *PodRepo) CreateBatch(pods []model.Pod) ([]model.Pod, model.RetVal) {
	if len(pods) == 0 {
		return pods, model.SuccessRetVal()
	}
	if res := p.DB.Model(&model.Pod{}).Create(&pods); res.Error != nil {
		log.Logger.Warningf("Failed to create Pods: %s", res.Error)
		return nil, model.RetVal{Msg: i18n.Model.CreateError, Attr: map[string]any{"Model": model.Name(model.Pod{}), "Error": res.Error.Error()}}
	}
	return pods, model.SuccessRetVal()
}

func (p *PodRepo) DeleteByVictimID(victimIDL ...uint) model.RetVal {
	return p.DeleteByFieldID("victim_id", victimIDL...)
}
