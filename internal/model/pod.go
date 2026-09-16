package model

// Pod K8s Pod 实例
// BelongsTo Victim
type Pod struct {
	VictimID uint    `json:"victim_id"`
	Victim   Victim  `json:"-"`
	Name     string  `json:"name"`
	Spec     PodSpec `gorm:"type:jsonb" json:"-"`
	BaseModel
}
