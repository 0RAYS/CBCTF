// Package log is copied from gorm.Logger and add TraceID
package log

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"strconv"
	"time"
)

type Level uint

const (
	// Silent silent log level
	Silent Level = iota + 1
	// Error error log level
	Error
	// Warn warn log level
	Warn
	// Info info log level
	Info
)

// ErrRecordNotFound record not found error
var ErrRecordNotFound = errors.New("record not found")

// Config gormTraceLogger config
type Config struct {
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
	LogLevel                  Level
}

// NewGormLogger initialize gormTraceLogger
func NewGormLogger(level Level) gormLogger.Interface {
	const (
		infoStr = "%s\n[info] "
		warnStr = "%s\n[warn] "
		errStr  = "%s\n[error] "
	)

	config := Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  level,
		IgnoreRecordNotFoundError: false,
	}

	return &gormTraceLogger{
		Entry:   Logger.WithField("type", "GORM"),
		Config:  config,
		infoStr: infoStr,
		warnStr: warnStr,
		errStr:  errStr,
	}
}

type gormTraceLogger struct {
	*logrus.Entry
	Config
	infoStr, warnStr, errStr string
}

// LogMode log mode
func (l *gormTraceLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	l.LogLevel = Level(level)
	return l
}

// Info print info
func (l *gormTraceLogger) Info(ctx context.Context, msg string, data ...any) {
	l.WithField("TraceID", ctx.Value("TraceID")).Infof(l.infoStr+msg, append([]any{utils.FileWithLineNum()}, data...)...)
}

// Warn print warn messages
func (l *gormTraceLogger) Warn(ctx context.Context, msg string, data ...any) {
	l.WithField("TraceID", ctx.Value("TraceID")).Warnf(l.warnStr+msg, append([]any{utils.FileWithLineNum()}, data...)...)
}

// Error print error messages
func (l *gormTraceLogger) Error(ctx context.Context, msg string, data ...any) {
	l.WithField("TraceID", ctx.Value("TraceID")).Errorf(l.errStr+msg, append([]any{utils.FileWithLineNum()}, data...)...)
}

// Trace print sql message
func (l *gormTraceLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= Silent {
		return
	}
	var e *logrus.Entry
	traceID := ctx.Value("TraceID")
	if traceID == nil {
		traceID = "00000000-0000-0000-0000-000000000000"
	}
	e = l.WithField("TraceID", ctx.Value("TraceID"))
	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= Error && (!errors.Is(err, ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			e.WithFields(logrus.Fields{
				"fileWithLineNum": utils.FileWithLineNum(),
				"duration":        fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6),
				"rows":            "-",
				"sql":             sql,
			}).Error(err)
		} else {
			e.WithFields(logrus.Fields{
				"fileWithLineNum": utils.FileWithLineNum(),
				"duration":        fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6),
				"rows":            strconv.FormatInt(rows, 10),
				"sql":             sql,
			}).Error(err)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			e.WithFields(logrus.Fields{
				"fileWithLineNum": utils.FileWithLineNum(),
				"duration":        fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6),
				"rows":            "-",
				"sql":             sql,
			}).Warn(slowLog)
		} else {
			e.WithFields(logrus.Fields{
				"fileWithLineNum": utils.FileWithLineNum(),
				"duration":        fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6),
				"rows":            strconv.FormatInt(rows, 10),
				"sql":             sql,
			}).Warn(slowLog)
		}

	case l.LogLevel == Info:
		sql, rows := fc()
		if rows == -1 {
			e.WithFields(logrus.Fields{
				"fileWithLineNum": utils.FileWithLineNum(),
				"duration":        fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6),
				"rows":            "-",
				"sql":             sql,
			}).Info()
		} else {
			e.WithFields(logrus.Fields{
				"fileWithLineNum": utils.FileWithLineNum(),
				"duration":        fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6),
				"rows":            strconv.FormatInt(rows, 10),
				"sql":             sql,
			}).Info()
		}
	}
}
