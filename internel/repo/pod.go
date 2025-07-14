package repo

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"gorm.io/gorm"
)

type PodRepo struct {
	BasicRepo[model.Pod]
}

type CreatePodOptions struct {
	VictimID     uint
	Name         string
	ExposedIP    string
	PodPorts     model.Exposes
	ExposedPorts model.Int32List
	IPs          model.IPs
}

func (c CreatePodOptions) Convert2Model() model.Model {
	return model.Pod{
		VictimID:     c.VictimID,
		Name:         c.Name,
		ExposedIP:    c.ExposedIP,
		PodPorts:     c.PodPorts,
		ExposedPorts: c.ExposedPorts,
		IPs:          c.IPs,
	}
}

type UpdatePodOptions struct {
	Name            *string
	ExposedIP       *string
	PodPorts        *model.Int32List
	ExposedPorts    *model.Int32List
	NetworkPolicies *model.NetworkPolicies
}

func (u UpdatePodOptions) Convert2Map() map[string]any {
	options := make(map[string]any)
	if u.Name != nil {
		options["name"] = *u.Name
	}
	if u.ExposedIP != nil {
		options["exposed_ip"] = *u.ExposedIP
	}
	if u.PodPorts != nil {
		options["pod_ports"] = *u.PodPorts
	}
	if u.ExposedPorts != nil {
		options["exposed_ports"] = *u.ExposedPorts
	}
	if u.NetworkPolicies != nil {
		options["network_policies"] = *u.NetworkPolicies
	}
	return options
}

func InitPodRepo(tx *gorm.DB) *PodRepo {
	return &PodRepo{
		BasicRepo: BasicRepo[model.Pod]{
			DB: tx,
		},
	}
}

func (p *PodRepo) Delete(idL ...uint) (bool, string) {
	podL, _, ok, msg := p.List(-1, -1, GetOptions{
		Conditions: map[string]any{"id": idL},
	})
	if !ok && msg != i18n.PodNotFound {
		return false, msg
	}
	containerIDL := make([]uint, 0)
	for _, pod := range podL {
		for _, container := range pod.Containers {
			containerIDL = append(containerIDL, container.ID)
		}
	}
	if ok, msg = InitContainerRepo(p.DB).Delete(containerIDL...); !ok {
		return false, msg
	}
	if res := p.DB.Model(&model.Pod{}).Where("id IN ?", idL).Delete(&model.Pod{}); res.Error != nil {
		log.Logger.Warningf("Failed to delete Pod: %s", res.Error)
		return false, i18n.DeletePodError
	}
	return true, i18n.Success
}
