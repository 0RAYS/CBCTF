package db

import (
	"time"

	"gorm.io/gorm"

	"CBCTF/internal/model"
)

// UpdateRunStatusOptions carries the run timestamps shared by the cron job,
// generator, SMTP and webhook records.
type UpdateRunStatusOptions struct {
	SuccessLast *time.Time
	FailureLast *time.Time
}

func (u UpdateRunStatusOptions) Convert2Map() map[string]any {
	options := make(map[string]any)
	if u.SuccessLast != nil {
		options["success_last"] = *u.SuccessLast
	}
	if u.FailureLast != nil {
		options["failure_last"] = *u.FailureLast
	}
	return options
}

// DiffUpdateRunStatusOptions increments the run counters of a record.
type DiffUpdateRunStatusOptions struct {
	Success int64
	Failure int64
}

func (d DiffUpdateRunStatusOptions) Convert2Expr() map[string]any {
	options := make(map[string]any)
	if d.Success != 0 {
		options["success"] = gorm.Expr("success + ?", d.Success)
	}
	if d.Failure != 0 {
		options["failure"] = gorm.Expr("failure + ?", d.Failure)
	}
	return options
}

// updateRunStatus counts a run of a scheduled record and stores when it happened.
func updateRunStatus[M model.Model](repo *BaseRepo[M], id uint, success bool, last time.Time) model.RetVal {
	var diffOptions DiffUpdateRunStatusOptions
	var options UpdateRunStatusOptions
	if success {
		diffOptions = DiffUpdateRunStatusOptions{Success: 1}
		options = UpdateRunStatusOptions{SuccessLast: &last}
	} else {
		diffOptions = DiffUpdateRunStatusOptions{Failure: 1}
		options = UpdateRunStatusOptions{FailureLast: &last}
	}
	if ret := repo.DiffUpdate(id, diffOptions); !ret.OK {
		return ret
	}
	return repo.Update(id, options)
}
