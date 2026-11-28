package log

import (
	"os"

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
	// The admin log page filters by level, so every level has to reach the mirror.
	Logger.SetLevel(logrus.DebugLevel)
}

// isTerminal keeps ANSI colors out of redirected output such as files and log collectors.
func isTerminal(f *os.File) bool {
	return isatty.IsTerminal(f.Fd()) || isatty.IsCygwinTerminal(f.Fd())
}
