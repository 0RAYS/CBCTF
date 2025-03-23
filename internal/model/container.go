package model

import (
	"CBCTF/internal/config"
	"CBCTF/internal/utils"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/plugin/optimisticlock"
	"time"
)

type Container struct {
	ID                 uint                   `json:"id" gorm:"primaryKey"`
	Port               int32                  `json:"port"`
	ContestID          uint                   `json:"contest_id"`
	TeamID             uint                   `json:"team_id"`
	ChallengeID        string                 `json:"challenge_id"`
	Start              time.Time              `json:"start"`
	PodName            string                 `json:"pod"`
	ContainerName      string                 `json:"container"`
	ServiceName        string                 `json:"service"`
	NetworkPolicyName  string                 `json:"network_policy"`
	PodNames           utils.Strings          `json:"pods" gorm:"type:json"`
	ContainerNames     utils.Strings          `json:"containers" gorm:"type:json"`
	ServiceNames       utils.Strings          `json:"services" gorm:"type:json"`
	NetworkPolicyNames utils.Strings          `json:"network_policies" gorm:"type:json"`
	IPPorts            utils.IPPorts          `json:"ip_ports" gorm:"type:json"`
	IP                 string                 `json:"ip"`
	Duration           time.Duration          `json:"-"`
	CreatorID          uint                   `json:"creator_id"`
	CreatedAt          time.Time              `json:"-"`
	UpdatedAt          time.Time              `json:"-"`
	DeletedAt          gorm.DeletedAt         `json:"-" gorm:"index"`
	Version            optimisticlock.Version `json:"-" gorm:"default:1"`
}

// MarshalJSON Duration 转为秒
func (c Container) MarshalJSON() ([]byte, error) {
	type Tmp Container
	return json.Marshal(struct {
		Tmp
		Duration int64 `json:"duration"`
	}{
		Tmp:      Tmp(c),
		Duration: int64(c.Duration.Seconds()),
	})
}

// TrafficPath 流量文件路径
func (c Container) TrafficPath() string {
	return fmt.Sprintf("%s/traffics/%s/%d/%d.pcap", config.Env.Path, c.ChallengeID, c.TeamID, c.ID)
}

// RemoteAddress 返回远程地址
func (c Container) RemoteAddress() string {
	return fmt.Sprintf("%s:%d", c.IP, c.Port)
}

// RemoteAddresses 返回多个远程地址
func (c Container) RemoteAddresses() []string {
	addresses := make([]string, 0, len(c.IPPorts))
	for _, ipPort := range c.IPPorts {
		addresses = append(addresses, fmt.Sprintf("%s:%d", ipPort.IP, ipPort.Port))
	}
	return addresses
}

// Remaining 返回剩余时间
func (c Container) Remaining() time.Duration {
	return c.Start.Add(c.Duration).Sub(time.Now())
}

func InitContainer(flag Flag, usage Usage, creatorID uint) Container {
	podName := fmt.Sprintf("victim-%s-%d-pod", usage.ChallengeID, flag.TeamID)
	serviceName := fmt.Sprintf("victim-%s-%d-svc", usage.ChallengeID, flag.TeamID)
	containerName := fmt.Sprintf("victim-%s-%d", usage.ChallengeID, flag.TeamID)
	networkPolicyName := fmt.Sprintf("victim-%s-%d-net", usage.ChallengeID, flag.TeamID)
	return Container{
		ContestID:          flag.ContestID,
		ChallengeID:        flag.ChallengeID,
		TeamID:             flag.TeamID,
		Port:               usage.Port,
		CreatorID:          creatorID,
		Start:              time.Now(),
		Duration:           1 * time.Hour,
		PodName:            podName,
		ContainerName:      containerName,
		ServiceName:        serviceName,
		NetworkPolicyName:  networkPolicyName,
		PodNames:           utils.Strings{podName},
		ContainerNames:     utils.Strings{containerName},
		ServiceNames:       utils.Strings{serviceName},
		NetworkPolicyNames: utils.Strings{networkPolicyName},
	}
}
