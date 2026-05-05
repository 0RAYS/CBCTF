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
	Logger.SetFormatter(Formatter{})
	Logger.SetLevel(logrus.DebugLevel)
}
