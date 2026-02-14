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
	PodPorts model.Exposes
	Networks model.Networks
}

func (c CreatePodOptions) Convert2Model() model.Model {
	return model.Pod{
		VictimID: c.VictimID,
		Name:     c.Name,
		PodPorts: c.PodPorts,
		Networks: c.Networks,
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
	podL, _, ret := p.List(-1, -1, GetOptions{
		Conditions: map[string]any{"id": idL},
		Preloads:   map[string]GetOptions{"Containers": {}},
	})
	if !ret.OK {
		if ret.Msg != i18n.Model.NotFound {
			return ret
		}
		return model.SuccessRetVal()
	}
	containerIDL := make([]uint, 0)
	for _, pod := range podL {
		for _, container := range pod.Containers {
			containerIDL = append(containerIDL, container.ID)
		}
	}
	if ret = InitContainerRepo(p.DB).Delete(containerIDL...); !ret.OK {
		return ret
	}
	if res := p.DB.Model(&model.Pod{}).Where("id IN ?", idL).Delete(&model.Pod{}); res.Error != nil {
		log.Logger.Warningf("Failed to delete Pod: %s", res.Error)
		return model.RetVal{Msg: i18n.Model.DeleteError, Attr: map[string]interface{}{"Model": model.Pod{}.ModelName(), "Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}
