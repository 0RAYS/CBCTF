package model

import (
	"CBCTF/internal/config"
	"fmt"
	"gorm.io/gorm"
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
	Duration      time.Duration  `json:"duration"`
	CreatorID     uint           `json:"creator_id"`
	CreatedAt     time.Time      `json:"-"`
	UpdatedAt     time.Time      `json:"-"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

func (d *Docker) RemoteAddr() string {
	return fmt.Sprintf("%s:%d", config.Env.K8S.Master, d.Port)
}

func InitDocker(flag Flag, challenge Challenge, creatorID uint) Docker {
	podName := fmt.Sprintf("docker-%s-%d-pod", challenge.ID, flag.TeamID)
	serviceName := fmt.Sprintf("docker-%s-%d-service", challenge.ID, flag.TeamID)
	containerName := fmt.Sprintf("%s-%d-docker", challenge.ID, flag.TeamID)
	return Docker{
		ContestID:     flag.ContestID,
		ChallengeID:   flag.ChallengeID,
		TeamID:        flag.TeamID,
		Port:          challenge.Port,
		CreatorID:     creatorID,
		Start:         time.Now(),
		Duration:      1 * time.Hour,
		PodName:       podName,
		ContainerName: containerName,
		ServiceName:   serviceName,
	}
}
