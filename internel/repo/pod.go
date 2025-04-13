package repo

import (
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"CBCTF/internel/utils"
	"gorm.io/gorm"
)

type PodRepo struct {
	Repo[model.Pod]
}

type CreatePodOptions struct {
	VictimID          uint
	Name              string
	ExposeIP          string
	PodIP             string
	ServiceName       string
	NetworkPolicyName string
	Exposes           model.Exposes
	NetworkPolicies   model.NetworkPolicies
}

type UpdatePodOptions struct {
	ExposeIP *string `json:"ip"`
	PodIP    *string `json:"pod_ip"`
}

func InitPodRepo(tx *gorm.DB) *PodRepo {
	return &PodRepo{Repo: Repo[model.Pod]{DB: tx, Model: "Pod"}}
}

func (p *PodRepo) GetByVictimID(victimID uint, deleted bool, preloadL ...string) ([]model.Pod, bool, string) {
	pods := make([]model.Pod, 0)
	res := p.DB.Model(&model.Pod{})
	if deleted {
		res = res.Unscoped()
	}
	res = res.Where("victim_id = ?", victimID)
	res = preload(res, preloadL...).Find(&pods)
	if res.RowsAffected == 0 {
		return pods, false, "GetPodError"
	}
	return pods, true, "Success"
}

func (p *PodRepo) GetByName(name string, deleted bool, preloadL ...string) ([]model.Pod, bool, string) {
	pods := make([]model.Pod, 0)
	res := p.DB.Model(&model.Pod{})
	if deleted {
		res = res.Unscoped()
	}
	res = res.Where("name = ?", name)
	res = preload(res, preloadL...).Find(&pods)
	if res.RowsAffected == 0 {
		return pods, false, "GetPodError"
	}
	return pods, true, "Success"
}

func (p *PodRepo) Update(id uint, options UpdateUserOptions) (bool, string) {
	var count int
	data := utils.UpdateOptions2Map(options)
	for {
		count++
		if count > 10 {
			log.Logger.Warningf("Failed to update Pod: too many times failed due to optimistic lock")
			return false, "DeadLock"
		}
		pod, ok, msg := p.GetByID(id)
		if !ok {
			return false, msg
		}
		data["version"] = pod.Version + 1
		res := p.DB.Model(&model.Pod{}).Omit("id", "created_at", "updated_at", "deleted_at").
			Where("id = ? AND version = ?", id, pod.Version).Updates(data)
		if res.Error != nil {
			log.Logger.Warningf("Failed to update Pod: %s", res.Error)
			return false, "UpdatePodError"
		}
		if res.RowsAffected == 0 {
			continue
		}
		break
	}
	return true, "Success"
}
