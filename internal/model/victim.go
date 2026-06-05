package model

import (
	"CBCTF/internal/config"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net/netip"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"
)

const (
	WaitingVictimStatus     = "waiting"
	PendingVictimStatus     = "pending"
	RunningVictimStatus     = "running"
	TerminatingVictimStatus = "terminating"
	StoppedVictimStatus     = "stopped"
)

// Victim 靶机实例
// BelongsTo Challenge
// BelongsTo Contest (nullable)
// BelongsTo ContestChallenge (nullable)
// BelongsTo Team (nullable)
// BelongsTo User
// HasMany Pod
type Victim struct {
	ChallengeID        uint             `gorm:"index" json:"challenge_id"`
	Challenge          Challenge        `json:"-"`
	ContestID          sql.Null[uint]   `gorm:"index" json:"contest_id"`
	Contest            Contest          `json:"-"`
	ContestChallengeID sql.Null[uint]   `gorm:"index" json:"contest_challenge_id"`
	ContestChallenge   ContestChallenge `json:"-"`
	TeamID             sql.Null[uint]   `gorm:"index" json:"team_id"`
	Team               Team             `json:"-"`
	UserID             uint             `gorm:"index" json:"user_id"`
	User               User             `json:"-"`
	Pods               []Pod            `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Start              time.Time        `gorm:"default:null" json:"start"`
	Duration           time.Duration    `json:"duration"`
	Spec               VictimSpec       `gorm:"default:null;type:jsonb" json:"-"`
	Resources          VictimResources  `gorm:"default:null;type:jsonb" json:"-"`
	Endpoints          Endpoints        `gorm:"default:null;type:jsonb" json:"-"`
	ExposedEndpoints   Endpoints        `gorm:"default:null;type:jsonb" json:"-"`
	Status             string           `gorm:"index" json:"status"`
	BaseModel
}

func (v Victim) TrafficBasePath() string {
	return filepath.Join(config.Env.Path, "traffics", "victim-"+strconv.FormatUint(uint64(v.ID), 10))
}

func (v Victim) TrafficZipPath() string {
	return filepath.Join(v.TrafficBasePath(), "traffics.zip")
}

func (v Victim) RemoteAddr() []string {
	data := make([]string, 0)
	for _, endpoint := range v.ExposedEndpoints {
		addr := fmt.Sprintf("%s://%s:%d", strings.ToLower(endpoint.Protocol), endpoint.IP, endpoint.Port)
		if strings.TrimSpace(endpoint.Name) != "" {
			addr = fmt.Sprintf("%s: %s", endpoint.Name, addr)
		}
		data = append(data, addr)
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
	NetAttachDef *NetAttachDef `json:"net_attach_def"`
}

type NetAttachDef struct {
	Name string `json:"name"`
}

func (v VPC) Value() (driver.Value, error) {
	v.Subnets = slices.DeleteFunc(v.Subnets, func(s *Subnet) bool {
		cidr, err := netip.ParsePrefix(s.CIDRBlock)
		if err != nil {
			return true
		}
		if gateway, err := netip.ParseAddr(s.Gateway); err != nil || !cidr.Contains(gateway) {
			return true
		}
		s.ExcludeIps = slices.DeleteFunc(s.ExcludeIps, func(ip string) bool {
			if i, err := netip.ParseAddr(ip); err != nil || !cidr.Contains(i) {
				return true
			}
			return false
		})
		return false
	})
	return json.Marshal(v)
}

func (v *VPC) Scan(value any) error {
	if err := scanJSON(value, v); err != nil {
		return fmt.Errorf("failed to scan VPC value: %v", value)
	}
	return nil
}

type Endpoint struct {
	Name     string `json:"name"`
	IP       string `json:"ip"`
	Port     int32  `json:"port"`
	Protocol string `json:"protocol"`
}

type Endpoints []Endpoint

func (e Endpoints) Value() (driver.Value, error) {
	e = slices.DeleteFunc(e, func(e Endpoint) bool {
		if _, err := netip.ParseAddr(e.IP); err != nil {
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
	if err := scanJSON(value, e); err != nil {
		return fmt.Errorf("failed to scan Endpoints value: %v", value)
	}
	return nil
}
