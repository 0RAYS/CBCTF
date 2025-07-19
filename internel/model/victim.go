package model

import (
	"CBCTF/internel/config"
	"CBCTF/internel/i18n"
	"CBCTF/internel/utils"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net"
	"slices"
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
	VPC                VPC              `gorm:"default:null;type:json" json:"-"`
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

type VPC struct {
	Name    string    `json:"name"`
	Subnets []*Subnet `json:"subnets"`
}

type Subnet struct {
	DefName      string        `json:"def_name"`
	Name         string        `json:"name"`
	CIDRBlock    string        `json:"cidr_block"`
	Gateway      string        `json:"gateway"`
	ExcludeIps   []string      `json:"exclude_ips"`
	NatGateway   *NatGateway   `json:"nat_gateway"`
	NetAttachDef *NetAttachDef `json:"net_attach_def"`
}

type NatGateway struct {
	Name  string `json:"name"`
	LanIP string `json:"lan_ip"`
	EIPs  []*EIP `json:"eips"`
}

type EIP struct {
	Name  string  `json:"name"`
	DNats []*DNat `json:"dnats"`
	SNats []*SNat `json:"snats"`
}

type DNat struct {
	Name         string `json:"name"`
	ExternalPort string `json:"external_port"`
	InternalIP   string `json:"internal_ip"`
	InternalPort string `json:"internal_port"`
	Protocol     string `json:"protocol"`
}

type SNat struct {
	Name string `json:"name"`
}

type NetAttachDef struct {
	Name string `json:"name"`
}

func (v VPC) Value() (driver.Value, error) {
	v.Subnets = slices.DeleteFunc(v.Subnets, func(s *Subnet) bool {
		_, cidr, err := net.ParseCIDR(s.CIDRBlock)
		if err != nil {
			return true
		}
		if gateway := net.ParseIP(s.Gateway); gateway == nil || !cidr.Contains(gateway) {
			return true
		}
		s.ExcludeIps = slices.DeleteFunc(s.ExcludeIps, func(ip string) bool {
			if i := net.ParseIP(ip); i == nil || !cidr.Contains(i) {
				return true
			}
			return false
		})
		if lanIP := net.ParseIP(s.NatGateway.LanIP); lanIP == nil || !cidr.Contains(lanIP) {
			return true
		}
		for _, eip := range s.NatGateway.EIPs {
			eip.DNats = slices.DeleteFunc(eip.DNats, func(d *DNat) bool {
				if i := net.ParseIP(d.InternalIP); i == nil || !cidr.Contains(i) {
					return true
				}
				return false
			})
		}
		return false
	})
	return json.Marshal(v)
}

func (v *VPC) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan VPC value: %v", value)
	}
	return json.Unmarshal(bytes, v)
}
