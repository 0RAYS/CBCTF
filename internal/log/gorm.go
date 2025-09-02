// Package log is copied from gorm.Logger and add TraceID
package log

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
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
		Entry:   Logger.WithField("Type", GormLogType),
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
func (l *gormTraceLogger) Info(_ context.Context, msg string, data ...any) {
	l.Infof(l.infoStr+msg, append([]any{utils.FileWithLineNum()}, data...)...)
}

// Warn print warn messages
func (l *gormTraceLogger) Warn(_ context.Context, msg string, data ...any) {
	l.Warnf(l.warnStr+msg, append([]any{utils.FileWithLineNum()}, data...)...)
}

// Error print error messages
func (l *gormTraceLogger) Error(_ context.Context, msg string, data ...any) {
	l.Errorf(l.errStr+msg, append([]any{utils.FileWithLineNum()}, data...)...)
}

// Trace print sql message
func (l *gormTraceLogger) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= Silent {
		return
	}
	elapsed := time.Since(begin)
	sql, rows := fc()
	fields := logrus.Fields{
		"FileWithLineNum": utils.FileWithLineNum(),
		"Duration":        fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6),
		"Rows":            "-",
		"SQL":             sql,
	}
	if rows != -1 {
		fields["Rows"] = strconv.Itoa(int(rows))
	}
	switch {
	case err != nil && l.LogLevel >= Error && (!errors.Is(err, ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		l.WithFields(fields).Error(err)
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= Warn:
		l.WithFields(fields).Warnf("SLOW SQL >= %v", l.SlowThreshold)
	case l.LogLevel == Info:
		l.WithFields(fields).Info()
	}
}
