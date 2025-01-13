package log

import (
	"CBCTF/internal/config"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

func Init() {
	Logger.SetFormatter(Formatter{})
	writer, err := rotatelogs.New("logs/%Y%m%d.log")
	if err != nil {
		Logger.Fatalln(err)
	}
	Logger.AddHook(lfshook.NewHook(
		lfshook.WriterMap{
			logrus.InfoLevel:  writer,
			logrus.ErrorLevel: writer,
			logrus.DebugLevel: writer,
			logrus.WarnLevel:  writer,
			logrus.TraceLevel: writer,
			logrus.FatalLevel: writer,
		}, &logrus.TextFormatter{},
	))
	level := config.Env.GetString("gin.log.level")
	switch level {
	case "Debug":
		Logger.SetLevel(logrus.DebugLevel)
	case "Error":
		Logger.SetLevel(logrus.ErrorLevel)
	case "Info":
		Logger.SetLevel(logrus.InfoLevel)
	case "Trace":
		Logger.SetLevel(logrus.TraceLevel)
	case "Warning":
		Logger.SetLevel(logrus.WarnLevel)
	}
	Logger.Debug(level)
}
