package model

import (
	"CBCTF/internel/i18n"
	"database/sql/driver"
	"encoding/json"
	"fmt"
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
	Challenge      Challenge       `json:"-"`
	ChallengeFlags []ChallengeFlag `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Name           string          `json:"name"`
	Image          string          `json:"image"`
	CPU            float32         `json:"cpu"`
	Memory         int64           `json:"memory"`
	WorkingDir     string          `json:"working_dir"`
	Command        StringList      `gorm:"default:null;type:json" json:"command"`
	Exposes        Exposes         `gorm:"default:null;type:json" json:"exposes"`
	Environment    StringMap       `gorm:"default:null;type:json" json:"environment"`
	BasicModel
}

func (d Docker) GetModelName() string {
	return "Docker"
}

func (d Docker) GetVersion() uint {
	return d.Version
}

func (d Docker) CreateErrorString() string {
	return i18n.CreateDockerError
}

func (d Docker) DeleteErrorString() string {
	return i18n.DeleteDockerError
}

func (d Docker) GetErrorString() string {
	return i18n.GetDockerError
}

func (d Docker) NotFoundErrorString() string {
	return i18n.DockerNotFound
}

func (d Docker) UpdateErrorString() string {
	return i18n.UpdateDockerError
}

func (d Docker) GetUniqueKey() []string {
	return []string{"id"}
}

func (d Docker) GetForeignKeys() []string {
	return []string{"id", "challenge_id"}
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
