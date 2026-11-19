// Package log provides the platform logger and the GORM adapter that feeds it.
package log

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type Level uint

const (
	Silent Level = iota + 1
	Error
	Warn
	Info
)

var ErrRecordNotFound = gorm.ErrRecordNotFound

type Config struct {
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
	LogLevel                  Level
}

func NewGormLogger(level Level) gormLogger.Interface {
	const (
		infoStr = "%s\n[info] "
		warnStr = "%s\n[warn] "
		errStr  = "%s\n[error] "
	)

	config := Config{
		SlowThreshold:             500 * time.Millisecond,
		LogLevel:                  level,
		IgnoreRecordNotFoundError: true,
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

func (l *gormTraceLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	l.LogLevel = Level(level)
	return l
}

func (l *gormTraceLogger) Info(_ context.Context, msg string, data ...any) {
	l.Infof(l.infoStr+msg, append([]any{utils.FileWithLineNum()}, data...)...)
}

func (l *gormTraceLogger) Warn(_ context.Context, msg string, data ...any) {
	l.Warnf(l.warnStr+msg, append([]any{utils.FileWithLineNum()}, data...)...)
}

func (l *gormTraceLogger) Error(_ context.Context, msg string, data ...any) {
	l.Errorf(l.errStr+msg, append([]any{utils.FileWithLineNum()}, data...)...)
}

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
