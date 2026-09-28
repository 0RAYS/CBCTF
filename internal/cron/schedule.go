package cron

import (
	"time"

	"github.com/robfig/cron/v3"
)

// startupSchedule makes the first execution due immediately. Keeping it inside
// the scheduler preserves job wrappers, execution tracking and shutdown waiting.
// Next is called only by the scheduler goroutine.
type startupSchedule struct {
	cron.Schedule
	pending bool
}

func (s *startupSchedule) Next(now time.Time) time.Time {
	if s.pending {
		s.pending = false
		return now
	}
	return s.Schedule.Next(now)
}
