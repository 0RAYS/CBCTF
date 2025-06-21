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
	VictimID        uint
	Name            string
	PodIP           string
	ExposedIP       string
	PodPorts        model.Ports
	ExposedPorts    model.Ports
	NetworkPolicies model.NetworkPolicies
}

func (c CreatePodOptions) Convert2Model() model.Model {
	return model.Pod{
		VictimID:        c.VictimID,
		Name:            c.Name,
		PodIP:           c.PodIP,
		ExposedIP:       c.ExposedIP,
		PodPorts:        c.PodPorts,
		ExposedPorts:    c.ExposedPorts,
		NetworkPolicies: c.NetworkPolicies,
	}
}

type UpdatePodOptions struct {
	Name            *string
	PodIP           *string
	ExposedIP       *string
	PodPorts        *model.Ports
	ExposedPorts    *model.Ports
	NetworkPolicies *model.NetworkPolicies
}

func (u UpdatePodOptions) Convert2Map() map[string]any {
	options := make(map[string]any)
	if u.Name != nil {
		options["name"] = *u.Name
	}
	if u.PodIP != nil {
		options["pod_ip"] = *u.PodIP
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
	containerIDL := make([]uint, 0)
	for _, id := range idL {
		pod, ok, msg := p.GetByID(id)
		if !ok && msg != i18n.PodNotFound {
			return false, msg
		}
		for _, container := range pod.Containers {
			containerIDL = append(containerIDL, container.ID)
		}
	}
	if ok, msg := InitContainerRepo(p.DB).Delete(containerIDL...); !ok {
		return false, msg
	}
	if res := p.DB.Model(&model.Pod{}).Where("id IN ?", idL).Delete(&model.Pod{}); res.Error != nil {
		log.Logger.Warningf("Failed to delete Pod: %s", res.Error)
		return false, model.Pod{}.DeleteErrorString()
	}
	return true, i18n.Success
}
