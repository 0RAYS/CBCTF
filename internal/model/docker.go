package model

import (
	"CBCTF/internal/config"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

type Docker struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	Port          int32          `json:"port"`
	ContestID     uint           `json:"contest_id"`
	TeamID        uint           `json:"team_id"`
	ChallengeID   string         `json:"challenge_id"`
	Start         time.Time      `json:"start"`
	PodName       string         `json:"pod"`
	ContainerName string         `json:"container"`
	ServiceName   string         `json:"service"`
	Duration      time.Duration  `json:"-"`
	CreatorID     uint           `json:"creator_id"`
	CreatedAt     time.Time      `json:"-"`
	UpdatedAt     time.Time      `json:"-"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

func (d Docker) MarshalJSON() ([]byte, error) {
	type Tmp Docker // 定义一个别名以避免递归调用
	return json.Marshal(&struct {
		Tmp
		Duration int64 `json:"duration"`
	}{
		Tmp:      Tmp(d),
		Duration: int64(d.Duration.Seconds()),
	})
}

func (d Docker) RemoteAddr() string {
	ip := config.Env.K8S.Nodes[rand.Intn(len(config.Env.K8S.Nodes))]
	return fmt.Sprintf("%s:%d", ip, d.Port)
}

func (d Docker) Remaining() time.Duration {
	return d.Start.Add(d.Duration).Sub(time.Now())
}

func InitDocker(flag Flag, challenge Challenge, creatorID uint, port int32) Docker {
	podName := fmt.Sprintf("docker-%s-%d-pod", challenge.ID, flag.TeamID)
	serviceName := fmt.Sprintf("docker-%s-%d-service", challenge.ID, flag.TeamID)
	containerName := fmt.Sprintf("%s-%d-docker", challenge.ID, flag.TeamID)
	return Docker{
		ContestID:     flag.ContestID,
		ChallengeID:   flag.ChallengeID,
		TeamID:        flag.TeamID,
		Port:          port,
		CreatorID:     creatorID,
		Start:         time.Now(),
		Duration:      1 * time.Hour,
		PodName:       podName,
		ContainerName: containerName,
		ServiceName:   serviceName,
	}
}
