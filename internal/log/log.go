package log

import (
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

const (
	DefaultLogType = "LOG"
	GormLogType    = "GORM"
	GinLogType     = "GIN"
	TaskLogType    = "TASK"
)

func Init() {
	Logger = logrus.New()
	Logger.SetReportCaller(true)
	Logger.SetFormatter(NewFormatter())
	// The admin log page filters by level, so every level has to reach the mirror.
	Logger.SetLevel(logrus.DebugLevel)
}
