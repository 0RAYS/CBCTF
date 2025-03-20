package model

import (
	"CBCTF/internal/config"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/plugin/optimisticlock"
	"time"
)

type Docker struct {
	ID                uint                   `json:"id" gorm:"primaryKey"`
	Port              int32                  `json:"port"`
	ContestID         uint                   `json:"contest_id"`
	TeamID            uint                   `json:"team_id"`
	ChallengeID       string                 `json:"challenge_id"`
	Start             time.Time              `json:"start"`
	PodName           string                 `json:"pod"`
	ContainerName     string                 `json:"container"`
	ServiceName       string                 `json:"service"`
	NetworkPolicyName string                 `json:"network_policy"`
	IP                string                 `json:"ip"`
	Duration          time.Duration          `json:"-"`
	CreatorID         uint                   `json:"creator_id"`
	CreatedAt         time.Time              `json:"-"`
	UpdatedAt         time.Time              `json:"-"`
	DeletedAt         gorm.DeletedAt         `json:"-" gorm:"index"`
	Version           optimisticlock.Version `json:"-" gorm:"default:1"`
}

// MarshalJSON Duration 转为秒
func (d *Docker) MarshalJSON() ([]byte, error) {
	type Tmp Docker // 定义一个别名以避免递归调用
	return json.Marshal(&struct {
		*Tmp
		Duration int64 `json:"duration"`
	}{
		Tmp:      (*Tmp)(d),
		Duration: int64(d.Duration.Seconds()),
	})
}

// TrafficPath 流量文件路径
func (d *Docker) TrafficPath() string {
	return fmt.Sprintf("%s/traffics/%s/%d/%d.pcap", config.Env.Path, d.ChallengeID, d.TeamID, d.ID)
}

// RemoteAddr 返回远程地址
func (d *Docker) RemoteAddr() string {
	return fmt.Sprintf("%s:%d", d.IP, d.Port)
}

// Remaining 返回剩余时间
func (d *Docker) Remaining() time.Duration {
	return d.Start.Add(d.Duration).Sub(time.Now())
}

func InitDocker(flag Flag, usage Usage, creatorID uint) Docker {
	podName := fmt.Sprintf("victim-%s-%d-pod", usage.ChallengeID, flag.TeamID)
	serviceName := fmt.Sprintf("victim-%s-%d-svc", usage.ChallengeID, flag.TeamID)
	containerName := fmt.Sprintf("victim-%s-%d", usage.ChallengeID, flag.TeamID)
	networkPolicyName := fmt.Sprintf("victim-%s-%d-net", usage.ChallengeID, flag.TeamID)
	return Docker{
		ContestID:         flag.ContestID,
		ChallengeID:       flag.ChallengeID,
		TeamID:            flag.TeamID,
		Port:              usage.Port,
		CreatorID:         creatorID,
		Start:             time.Now(),
		Duration:          1 * time.Hour,
		PodName:           podName,
		ContainerName:     containerName,
		ServiceName:       serviceName,
		NetworkPolicyName: networkPolicyName,
	}
}
