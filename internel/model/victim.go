package model

import (
	"CBCTF/internel/config"
	"fmt"
	"time"
)

type Victim struct {
	UsageID  uint          `json:"usage_id"`
	Usage    Usage         `json:"-"`
	TeamID   uint          `json:"team_id"`
	Team     Team          `json:"-"`
	UserID   uint          `json:"user_id"`
	User     User          `json:"-"`
	Pods     []Pod         `json:"-"`
	Traffics []Traffic     `json:"-"`
	IPBlock  string        `json:"ip_block"`
	Start    time.Time     `json:"start"`
	Duration time.Duration `json:"duration"`
	BaseModel
}

func (v Victim) TrafficZipPath() string {
	return fmt.Sprintf("%s/traffics/%d/traffics.zip", config.Env.Path, v.ID)
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
