package log

import (
	"strings"

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
	Logger.SetLevel(logrus.InfoLevel)
}

// SetLevel applies the configured application log level, keeping Info for unknown values.
func SetLevel(level string) {
	parsed, err := logrus.ParseLevel(strings.ToLower(strings.TrimSpace(level)))
	if err != nil {
		parsed = logrus.InfoLevel
	}
	Logger.SetLevel(parsed)
}
