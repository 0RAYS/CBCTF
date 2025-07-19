package model

import (
	"CBCTF/internel/config"
	"CBCTF/internel/i18n"
	"CBCTF/internel/utils"
	"fmt"
	"time"
)

type Victim struct {
	ContestChallengeID uint             `json:"contest_challenge_id"`
	ContestChallenge   ContestChallenge `json:"-"`
	TeamID             uint             `json:"team_id"`
	Team               Team             `json:"-"`
	UserID             uint             `json:"user_id"`
	User               User             `json:"-"`
	Pods               []Pod            `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Traffics           []Traffic        `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Start              time.Time        `json:"start"`
	Duration           time.Duration    `json:"duration"`
	VPC                string           `json:"vpc"`
	Subnets            StringList       `gorm:"default:null;type:json" json:"subnets"`
	NetAttachDefs      StringList       `gorm:"default:null;type:json" json:"net_attach_defs"`
	Gateways           StringList       `gorm:"default:null;type:json" json:"gateways"`
	EIPs               StringList       `gorm:"default:null;type:json" json:"eips"`
	DNats              StringList       `gorm:"default:null;type:json" json:"dnats"`
	SNats              StringList       `gorm:"default:null;type:json" json:"snats"`
	NetworkPolicies    NetworkPolicies  `gorm:"default:null;type:json" json:"network_policies"`
	BasicModel
}

func (v Victim) GetModelName() string {
	return "Victim"
}

func (v Victim) GetVersion() uint {
	return v.Version
}

func (v Victim) CreateErrorString() string {
	return i18n.CreateVictimError
}

func (v Victim) DeleteErrorString() string {
	return i18n.DeleteVictimError
}

func (v Victim) GetErrorString() string {
	return i18n.GetVictimError
}

func (v Victim) NotFoundErrorString() string {
	return i18n.VictimNotFound
}

func (v Victim) UpdateErrorString() string {
	return i18n.UpdateVictimError
}

func (v Victim) GetUniqueKey() []string {
	return []string{"id"}
}

func (v Victim) GetForeignKeys() []string {
	return []string{"id", "contest_challenge_id", "team_id", "user_id"}
}

func (v Victim) GenPodName(challengeRandID string) string {
	return fmt.Sprintf("victim-%s-%s-pod", challengeRandID, utils.RandStr(5))
}

func (v Victim) TrafficZipPath() string {
	return fmt.Sprintf("%s/traffics/victim-%d/traffics.zip", config.Env.Path, v.ID)
}

// TrafficPaths Victim 需要预加载 Pod
func (v Victim) TrafficPaths() []string {
	data := make([]string, 0)
	for _, pod := range v.Pods {
		data = append(data, pod.TrafficPath())
	}
	return data
}

// RemoteAddr Victim 需要预加载 Pod
func (v Victim) RemoteAddr() []string {
	data := make([]string, 0)
	for _, pod := range v.Pods {
		data = append(data, pod.RemoteAddr()...)
	}
	return data
}

func (v Victim) Remaining() time.Duration {
	return v.Start.Add(v.Duration).Sub(time.Now())
}
