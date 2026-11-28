package view

import "time"

// RangePoolView 是一个靶场题目的实例池状态。
type RangePoolView struct {
	ChallengeID   uint
	ChallengeName string
	PoolSize      int
	MaxLeased     int
	Leased        int
	Instances     []RangeInstanceView
}

// RangeInstanceView 是池中一个靶场实例的运行状态。
type RangeInstanceView struct {
	ID            uint
	Status        string
	TeamID        uint
	TeamName      string
	LeaseExpireAt time.Time
	TemplateKey   string
	Nodes         int
	Start         time.Time
}
