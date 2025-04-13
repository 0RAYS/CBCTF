package model

import (
	"gorm.io/gorm"
	"time"
)

type Pod struct {
	ID                uint            `gorm:"primaryKey" json:"id"`
	VictimID          uint            `json:"victim_id"`
	Victim            Victim          `json:"-"`
	Containers        []Container     `json:"-"`
	Traffics          []Traffic       `json:"-"`
	Name              string          `json:"name"`
	ExposeIP          string          `json:"ip"`
	PodIP             string          `json:"pod_ip"`
	ServiceName       string          `json:"service"`
	NetworkPolicyName string          `json:"network_policy"`
	Exposes           Exposes         `gorm:"type:json" json:"exposes"`
	NetworkPolicies   NetworkPolicies `json:"network_policies"`
	CreatedAt         time.Time       `json:"-"`
	UpdatedAt         time.Time       `json:"-"`
	DeletedAt         gorm.DeletedAt  `gorm:"index" json:"-"`
	Version           uint            `gorm:"default:1" json:"-"`
}
