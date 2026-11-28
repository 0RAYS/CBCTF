package log

import (
	"os"
	"strings"

	"github.com/mattn/go-isatty"
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
	Logger.SetFormatter(NewFormatter(isTerminal(os.Stderr)))
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

// isTerminal keeps ANSI colors out of redirected output such as files and log collectors.
func isTerminal(f *os.File) bool {
	return isatty.IsTerminal(f.Fd()) || isatty.IsCygwinTerminal(f.Fd())
}
