package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net/netip"
	"slices"
)

const (
	EnvFlagPrefix      = "FLAG"
	VolumeFlagPrefix   = "FLAG"
	VolumeFlagLabelKey = "value"
)

// Docker
// BelongsTo Challenge
// HasMany ChallengeFlag
type Docker struct {
	ChallengeID    uint            `gorm:"index" json:"challenge_id"`
	Challenge      Challenge       `json:"-"`
	ChallengeFlags []ChallengeFlag `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Name           string          `json:"name"`
	Image          string          `json:"image"`
	CPU            float32         `json:"cpu"`
	Memory         int64           `json:"memory"`
	WorkingDir     string          `json:"working_dir"`
	Command        StringList      `gorm:"default:null;type:jsonb" json:"command"`
	Exposes        Exposes         `gorm:"default:null;type:jsonb" json:"exposes"`
	Environment    StringMap       `gorm:"default:null;type:jsonb" json:"environment"`
	Networks       Networks        `gorm:"default:null;type:jsonb" json:"networks"`
	BaseModel
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
		cidr, err := netip.ParsePrefix(n.CIDR)
		if err != nil {
			return true
		}
		gateway, err := netip.ParseAddr(n.Gateway)
		if err != nil || !cidr.Contains(gateway) {
			return true
		}
		if n.IP == "" {
			return false
		}
		ip, err := netip.ParseAddr(n.IP)
		if err != nil {
			return true
		}
		return !cidr.Contains(ip)
	})
	return json.Marshal(n)
}

func (n *Networks) Scan(value any) error {
	if err := scanJSON(value, n); err != nil {
		return fmt.Errorf("failed to scan Networks value")
	}
	return nil
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
	if err := scanJSON(value, e); err != nil {
		return fmt.Errorf("failed to scan Exposes value")
	}
	return nil
}
