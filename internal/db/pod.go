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

func (p *PodRepo) Delete(idL ...uint) model.RetVal {
	if len(idL) == 0 {
		return model.SuccessRetVal()
	}
	if res := p.DB.Model(&model.Pod{}).Where("id IN ?", idL).Delete(&model.Pod{}); res.Error != nil {
		log.Logger.Warningf("Failed to delete Pod: %s", res.Error)
		return model.RetVal{Msg: i18n.Model.Pod.DeleteError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}

func (p *PodRepo) DeleteByVictimID(victimIDL ...uint) model.RetVal {
	return p.DeleteByFieldID("victim_id", victimIDL...)
}
