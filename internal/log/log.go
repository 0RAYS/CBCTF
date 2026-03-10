package log

import (
	"CBCTF/internal/config"
	"strings"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
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
	if config.Env.Log.Save {
		writer, err := rotatelogs.New("logs/%Y%m%d.log")
		if err != nil {
			Logger.Fatal(err)
		}
		Logger.AddHook(lfshook.NewHook(
			lfshook.WriterMap{
				logrus.InfoLevel:  writer,
				logrus.ErrorLevel: writer,
				logrus.DebugLevel: writer,
				logrus.WarnLevel:  writer,
				logrus.TraceLevel: writer,
				logrus.FatalLevel: writer,
			}, Formatter{},
		))
	}
	switch strings.ToUpper(config.Env.Log.Level) {
	case "DEBUG":
		Logger.SetLevel(logrus.DebugLevel)
	case "ERROR":
		Logger.SetLevel(logrus.ErrorLevel)
	case "INFO":
		Logger.SetLevel(logrus.InfoLevel)
	case "WARNING":
		Logger.SetLevel(logrus.WarnLevel)
	default:
		Logger.SetLevel(logrus.InfoLevel)
	}
}
