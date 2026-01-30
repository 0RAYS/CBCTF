package model

import (
	"CBCTF/internal/config"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net"
	"slices"
	"strings"
	"time"
)

type Victim struct {
	ChallengeID        uint             `json:"challenge_id"`
	ContestID          sql.Null[uint]   `json:"contest_id"`
	ContestChallengeID sql.Null[uint]   `json:"contest_challenge_id"`
	ContestChallenge   ContestChallenge `json:"-"`
	TeamID             sql.Null[uint]   `json:"team_id"`
	Team               Team             `json:"-"`
	UserID             sql.Null[uint]   `json:"user_id"`
	User               User             `json:"-"`
	Pods               []Pod            `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Start              time.Time        `json:"start"`
	Duration           time.Duration    `json:"duration"`
	VPC                VPC              `gorm:"default:null;type:json" json:"-"`
	Endpoints          Endpoints        `gorm:"default:null;type:json" json:"-"`
	ExposedEndpoints   Endpoints        `gorm:"default:null;type:json" json:"-"`
	NetworkPolicies    NetworkPolicies  `gorm:"default:null;type:json" json:"network_policies"`
	BaseModel
}

func (v Victim) GetModelName() string {
	return "Victim"
}

func (v Victim) GetBaseModel() BaseModel {
	return v.BaseModel
}

func (v Victim) GetUniqueKey() []string {
	return []string{"id"}
}

func (v Victim) GetAllowedQueryFields() []string {
	return []string{}
}

func (v Victim) TrafficBasePath() string {
	return fmt.Sprintf("%s/traffics/victim-%d", config.Env.Path, v.ID)
}

func (v Victim) TrafficZipPath() string {
	return fmt.Sprintf("%s/traffics.zip", v.TrafficBasePath())
}

// TrafficPaths Victim 需要预加载 Pod
func (v Victim) TrafficPaths() []string {
	data := make([]string, 0)
	for _, pod := range v.Pods {
		data = append(data, pod.TrafficPcapPath())
	}
	return data
}

func (v Victim) RemoteAddr() []string {
	data := make([]string, 0)
	for _, endpoint := range v.ExposedEndpoints {
		data = append(data, fmt.Sprintf("%s://%s:%d", strings.ToLower(endpoint.Protocol), endpoint.IP, endpoint.Port))
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
	IP    string  `json:"ip"`
	DNats []*DNat `json:"dnats"`
	SNats []*SNat `json:"snats"`
}

type DNat struct {
	Name         string `json:"name"`
	ExternalPort int32  `json:"external_port"`
	InternalIP   string `json:"internal_ip"`
	InternalPort int32  `json:"internal_port"`
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
		if s.NatGateway != nil {
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

type Endpoint struct {
	IP       string `json:"ip"`
	Port     int32  `json:"port"`
	Protocol string `json:"protocol"`
}

type Endpoints []Endpoint

func (e Endpoints) Value() (driver.Value, error) {
	e = slices.DeleteFunc(e, func(e Endpoint) bool {
		if net.ParseIP(e.IP) == nil {
			return true
		}
		if e.Port < 0 || e.Port > 65535 {
			return true
		}
		if strings.ToLower(e.Protocol) != "tcp" && strings.ToLower(e.Protocol) != "udp" {
			return true
		}
		return false
	})
	return json.Marshal(e)
}

func (e *Endpoints) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan Endpoints value: %v", value)
	}
	return json.Unmarshal(bytes, e)
}
