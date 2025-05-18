package repo

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"CBCTF/internel/utils"
	"gorm.io/gorm"
)

type PodRepo struct {
	Repo[model.Pod]
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

type UpdatePodOptions struct {
	ExposedIP    *string      `json:"exposed_ip"`
	ExposedPorts *model.Ports `json:"exposed_ports"`
}

func InitPodRepo(tx *gorm.DB) *PodRepo {
	return &PodRepo{Repo: Repo[model.Pod]{DB: tx, Model: "Pod"}}
}

func (p *PodRepo) GetByPodName(name string, deleted bool, preloadL ...string) ([]model.Pod, bool, string) {
	pods := make([]model.Pod, 0)
	res := p.DB.Model(&model.Pod{})
	if deleted {
		res = res.Unscoped()
	}
	res = res.Where("name = ?", name)
	res = preload(res, preloadL...).Find(&pods)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Pod: %s", res.Error)
		return pods, false, i18n.GetPodError
	}
	if res.RowsAffected == 0 {
		return pods, false, i18n.PodNotFound
	}
	return pods, true, i18n.Success
}

func (p *PodRepo) GetByVictimID(victimID uint, deleted bool, preloadL ...string) ([]model.Pod, bool, string) {
	pods := make([]model.Pod, 0)
	res := p.DB.Model(&model.Pod{})
	if deleted {
		res = res.Unscoped()
	}
	res = res.Where("victim_id = ?", victimID)
	res = preload(res, preloadL...).Find(&pods)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Pod: %s", res.Error)
		return pods, false, i18n.GetPodError
	}
	if res.RowsAffected == 0 {
		return pods, false, i18n.PodNotFound
	}
	return pods, true, i18n.Success
}

func (p *PodRepo) Update(id uint, options UpdatePodOptions) (bool, string) {
	var count int
	data := utils.UpdateOptions2Map(options)
	for {
		count++
		if count > 10 {
			log.Logger.Warningf("Failed to update Pod: too many times failed due to optimistic lock")
			return false, i18n.DeadLock
		}
		pod, ok, msg := p.GetByID(id)
		if !ok {
			return false, msg
		}
		data["version"] = pod.Version + 1
		res := p.DB.Model(&model.Pod{}).Where("id = ? AND version = ?", id, pod.Version).Updates(data)
		if res.Error != nil {
			log.Logger.Warningf("Failed to update Pod: %s", res.Error)
			return false, i18n.UpdatePodError
		}
		if res.RowsAffected == 0 {
			continue
		}
		break
	}
	return true, i18n.Success
}

func (p *PodRepo) Delete(idL ...uint) (bool, string) {
	containerIDL := make([]uint, 0)
	for _, id := range idL {
		pod, ok, msg := p.GetByID(id, "Containers")
		if !ok {
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
		return false, i18n.DeletePodError
	}
	return true, i18n.Success
}
