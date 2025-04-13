package model

import (
	"gorm.io/gorm"
	"time"
)

type Victim struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UsageID   uint           `json:"usage_id"`
	Usage     Usage          `json:"-"`
	TeamID    uint           `json:"team_id"`
	Team      Team           `json:"-"`
	UserID    uint           `json:"user_id"`
	User      User           `json:"-"`
	Pods      []Pod          `json:"-"`
	IPBlock   string         `json:"ip_block"`
	Start     time.Time      `json:"start"`
	Duration  time.Duration  `json:"duration"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Version   uint           `gorm:"default:1" json:"-"`
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
