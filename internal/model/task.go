package model

import (
	"time"
)

const (
	TaskSuccessStatus = "success"
	TaskFailedStatus  = "failed"
)

// Task stores terminal task execution records only.
// A record is written when a task succeeds or when it finally fails.
type Task struct {
	TaskID      string    `gorm:"type:varchar(255);index;not null" json:"task_id"`
	Type        string    `gorm:"type:varchar(255);index;not null" json:"type"`
	Queue       string    `gorm:"type:varchar(255);index;not null" json:"queue"`
	Status      string    `gorm:"type:varchar(32);index;not null" json:"status"`
	Payload     any       `gorm:"serializer:json;type:jsonb" json:"payload"`
	Result      any       `gorm:"serializer:json;type:jsonb" json:"result"`
	Error       string    `gorm:"type:text" json:"error"`
	RetryCount  int       `gorm:"not null;default:0" json:"retry_count"`
	MaxRetry    int       `gorm:"not null;default:0" json:"max_retry"`
	ProcessedAt time.Time `gorm:"index;not null" json:"processed_at"`
	BaseModel
}
