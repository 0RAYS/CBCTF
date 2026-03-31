package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

const (
	TaskSuccessStatus = "success"
	TaskFailedStatus  = "failed"
)

type TaskPayload struct {
	V any
}

func (p TaskPayload) Value() (driver.Value, error) {
	if p.V == nil {
		return nil, nil
	}
	return json.Marshal(p.V)
}

func (p *TaskPayload) Scan(value any) error {
	if value == nil {
		p.V = nil
		return nil
	}
	if err := scanJSON(value, &p.V); err != nil {
		return fmt.Errorf("failed to scan TaskPayload value")
	}
	return nil
}

// Task stores terminal task execution records only.
// A record is written when a task succeeds or when it finally fails.
type Task struct {
	TaskID      string      `gorm:"type:varchar(255);index;not null" json:"task_id"`
	Type        string      `gorm:"type:varchar(255);index;not null" json:"type"`
	Queue       string      `gorm:"type:varchar(255);index;not null" json:"queue"`
	Status      string      `gorm:"type:varchar(32);index;not null" json:"status"`
	Payload     TaskPayload `gorm:"type:jsonb" json:"payload"`
	Result      TaskPayload `gorm:"type:jsonb" json:"result"`
	Error       string      `gorm:"type:text" json:"error"`
	RetryCount  int         `gorm:"not null;default:0" json:"retry_count"`
	MaxRetry    int         `gorm:"not null;default:0" json:"max_retry"`
	ProcessedAt time.Time   `gorm:"index;not null" json:"processed_at"`
	BaseModel
}
