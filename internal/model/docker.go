package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net"
	"slices"
)

const (
	EnvFlagPrefix      = "FLAG"
	VolumeFlagPrefix   = "FLAG"
	VolumeFlagLabelKey = "value"
)

// Docker
// HasMany ChallengeFlag
type Docker struct {
	ChallengeID    uint            `json:"challenge_id"`
	ChallengeFlags []ChallengeFlag `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Name           string          `json:"name"`
	Image          string          `json:"image"`
	CPU            float32         `json:"cpu"`
	Memory         int64           `json:"memory"`
	WorkingDir     string          `json:"working_dir"`
	Command        StringList      `gorm:"default:null;type:json" json:"command"`
	Exposes        Exposes         `gorm:"default:null;type:json" json:"exposes"`
	Environment    StringMap       `gorm:"default:null;type:json" json:"environment"`
	Networks       Networks        `gorm:"default:null;type:json" json:"networks"`
	BaseModel
}

func (d Docker) ModelName() string {
	return "Docker"
}

func (d Docker) GetBaseModel() BaseModel {
	return d.BaseModel
}

func (d Docker) UniqueFields() []string {
	return []string{"id"}
}

func (d Docker) QueryFields() []string {
	return []string{}
}

type Network struct {
	Name     string `json:"name"`
	CIDR     string `json:"cidr"`
	Gateway  string `json:"gateway"`
	IP       string `json:"ip"`
	External bool   `json:"external"`
}

type Networks []Network

func (n Networks) Value() (driver.Value, error) {
	n = slices.DeleteFunc(n, func(n Network) bool {
		if n.CIDR == "" || n.Gateway == "" || n.IP == "" {
			return true
		}
		_, cidr, err := net.ParseCIDR(n.CIDR)
		if err != nil {
			return true
		}
		gateway := net.ParseIP(n.Gateway)
		if gateway == nil || !cidr.Contains(gateway) {
			return true
		}
		if n.IP == "" {
			return false
		}
		ip := net.ParseIP(n.IP)
		if ip == nil {
			return true
		}
		return !cidr.Contains(ip)
	})
	return json.Marshal(n)
}

func (n *Networks) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan Networks value")
	}
	return json.Unmarshal(bytes, n)
}

type Expose struct {
	Port     int32  `json:"port"`
	Protocol string `json:"protocol"`
}

type Exposes []Expose

func (e Exposes) Value() (driver.Value, error) {
	e = slices.DeleteFunc(e, func(n Expose) bool {
		if n.Port < 0 || n.Port > 65535 {
			return true
		}
		return false
	})
	return json.Marshal(e)
}

func (e *Exposes) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan Exposes value")
	}
	return json.Unmarshal(bytes, e)
}
