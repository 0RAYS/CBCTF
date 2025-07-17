package model

import (
	"CBCTF/internel/config"
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
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
	VPC                string           `json:"vpc"`
	Subnets            Subnets          `gorm:"default:null;type:json" json:"subnets"`
	NetAttachDefs      StringMap        `gorm:"default:null;type:json" json:"net_attach_defs"`
	Gateways           Gateways         `gorm:"default:null;type:json" json:"gateways"`
	EIPs               EIPs             `gorm:"default:null;type:json" json:"eips"`
	DNats              DNats            `gorm:"default:null;type:json" json:"dnats"`
	SNats              SNats            `gorm:"default:null;type:json" json:"snats"`
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

type Subnet struct {
	Name     string `json:"name"`
	CIDR     string `json:"cidr"`
	Gateway  string `json:"gateway"`
	External bool   `json:"external"`
}

type Subnets []Subnet

func (s Subnets) Value() (driver.Value, error) {
	s = slices.DeleteFunc(s, func(s Subnet) bool {
		if s.CIDR == "" || s.Gateway == "" {
			return true
		}
		_, cidr, err := net.ParseCIDR(s.CIDR)
		if err != nil {
			return true
		}
		s.CIDR = cidr.String()
		if s.Gateway == "" {
			s.Gateway, err = utils.GetFirstIP(s.CIDR)
			if err != nil {
				log.Logger.Warningf("Get first IP fail: %v", err)
				return true
			}
		}
		gateway := net.ParseIP(s.Gateway)
		if gateway == nil || !cidr.Contains(gateway) {
			return true
		}
		return false
	})
	return json.Marshal(s)
}

func (s *Subnets) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan Subnets value")
	}
	return json.Unmarshal(bytes, s)
}

type Gateway struct {
	Name   string `json:"name"`
	VPC    string `json:"vpc"`
	Subnet string `json:"subnet"`
	LanIP  string `json:"lan_ip"`
}

type Gateways []Gateway

func (g Gateways) Value() (driver.Value, error) {
	g = slices.DeleteFunc(g, func(g Gateway) bool {
		if net.ParseIP(g.LanIP) == nil {
			return true
		}
		return false
	})
	return json.Marshal(g)
}

func (g *Gateways) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan Gateways value")
	}
	return json.Unmarshal(bytes, g)
}

type EIP struct {
	Name    string `json:"name"`
	Gateway string `json:"gateway"`
	IP      string `json:"ip"`
}

type EIPs []EIP

func (e EIPs) Value() (driver.Value, error) {
	e = slices.DeleteFunc(e, func(e EIP) bool {
		if net.ParseIP(e.IP) == nil {
			return true
		}
		return false
	})
	return json.Marshal(e)
}

func (e *EIPs) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan EIPs value")
	}
	return json.Unmarshal(bytes, e)
}

type DNat struct {
	Name         string `json:"name"`
	EIP          string `json:"eip"`
	ExternalPort int32  `json:"external_port"`
	InternalIP   string `json:"internal_ip"`
	InternalPort int32  `json:"internal_port"`
	Protocol     string `json:"protocol"`
}

type DNats []DNat

func (d DNats) Value() (driver.Value, error) {
	d = slices.DeleteFunc(d, func(d DNat) bool {
		if d.ExternalPort < 0 || d.ExternalPort > 65535 || d.InternalPort < 0 || d.InternalPort > 65535 {
			return true
		}
		if net.ParseIP(d.InternalIP) == nil {
			return true
		}
		return false
	})
	return json.Marshal(d)
}

func (d *DNats) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan DNats value")
	}
	return json.Unmarshal(bytes, d)
}

type SNat struct {
	Name         string `json:"name"`
	EIP          string `json:"eip"`
	InternalCIDR string `json:"internal_cidr"`
}

type SNats []SNat

func (s SNats) Value() (driver.Value, error) {
	s = slices.DeleteFunc(s, func(s SNat) bool {
		if net.ParseIP(s.InternalCIDR) == nil {
			return true
		}
		return false
	})
	return json.Marshal(s)
}

func (s *SNats) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan SNats value")
	}
	return json.Unmarshal(bytes, s)
}
