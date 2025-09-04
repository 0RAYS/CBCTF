package log

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

var colors = map[string]func(a ...any) string{
	"Time":    color.New(color.FgGreen).Add(color.Bold).SprintFunc(),
	"Warning": color.New(color.FgYellow).Add(color.Bold).SprintFunc(),
	"Panic":   color.New(color.FgRed, color.BgWhite).SprintFunc(),
	"Error":   color.New(color.FgRed).Add(color.Bold).SprintFunc(),
	"Info":    color.New(color.FgGreen).Add(color.Bold).SprintFunc(),
	"Debug":   color.New(color.FgBlue).Add(color.Bold).SprintFunc(),
	"GET":     color.New(color.BgBlue).SprintFunc(),
	"POST":    color.New(color.BgCyan).SprintFunc(),
	"PUT":     color.New(color.BgYellow).SprintFunc(),
	"DELETE":  color.New(color.BgRed).SprintFunc(),
	"PATCH":   color.New(color.BgGreen).SprintFunc(),
	"HEAD":    color.New(color.BgMagenta).SprintFunc(),
	"OPTIONS": color.New(color.BgWhite, color.FgBlack).SprintFunc(),
	"Default": color.New().SprintFunc(),
}

func methodColor(method string) func(a ...any) string {
	switch method {
	case http.MethodGet:
		return colors["GET"]
	case http.MethodPost:
		return colors["POST"]
	case http.MethodPut:
		return colors["PUT"]
	case http.MethodDelete:
		return colors["DELETE"]
	case http.MethodPatch:
		return colors["PATCH"]
	case http.MethodHead:
		return colors["HEAD"]
	case http.MethodOptions:
		return colors["OPTIONS"]
	default:
		return colors["Default"]
	}
}

func levelColor(level logrus.Level) func(a ...any) string {
	switch level {
	case logrus.DebugLevel, logrus.TraceLevel:
		return colors["Debug"]
	case logrus.ErrorLevel:
		return colors["Error"]
	case logrus.WarnLevel:
		return colors["Warning"]
	case logrus.InfoLevel:
		return colors["Info"]
	case logrus.PanicLevel:
		return colors["Panic"]
	case logrus.FatalLevel:
		return colors["Error"]
	default:
		return colors["Time"]
	}
}

func statusCodeColor(code int) func(a ...any) string {
	switch {
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		return colors["Info"]
	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
		return colors["Debug"]
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		return colors["Warning"]
	default:
		return colors["Error"]
	}
}

func safeGetValue[T any](entry *logrus.Entry, key string, defaultV T) T {
	V, ok := entry.Data[key].(T)
	if !ok {
		return defaultV
	}
	return V
}

type Formatter struct{}

func (f Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	color.NoColor = false
	t, ok := entry.Data["Type"].(string)
	if !ok {
		t = DefaultLogType
	}
	LevelColor := levelColor(entry.Level)
	LevelText := fmt.Sprintf("%-12s", t+"-"+entry.Level.String())
	ret := new(bytes.Buffer)
	_, _ = fmt.Fprintf(ret, "%s %s | ", LevelColor(LevelText), entry.Time.Format("2006-01-02 15:04:05"))
	switch t {
	case DefaultLogType:
		filepath := strings.SplitN(entry.Caller.File, "/CBCTF/", 2)
		caller := fmt.Sprintf("%s:%d", filepath[len(filepath)-1], entry.Caller.Line)
		caller = fmt.Sprintf("%-36s", caller)
		_, _ = fmt.Fprintf(ret, "%s | %s", caller, LevelColor(entry.Message))
	case TaskLogType:
		_, _ = fmt.Fprintf(ret, "%-36s | %s", "Async Queue Task", LevelColor(entry.Message))
	case GinLogType:
		StatusCodeColor := statusCodeColor(safeGetValue(entry, "StatusCode", -1))
		MethodColor := methodColor(safeGetValue(entry, "Method", "ERROR"))
		Latency := safeGetValue(entry, "Latency", time.Duration(0))
		_, _ = fmt.Fprintf(ret, "%s | %s | %13v | ",
			safeGetValue(entry, "TraceID", "00000000-0000-0000-0000-000000000000"),
			StatusCodeColor(safeGetValue(entry, "StatusCode", -1)),
			Latency,
		)
		_, _ = fmt.Fprintf(ret, "%s | %s | \"%s\"",
			fmt.Sprintf("%-15s", safeGetValue(entry, "ClientIP", "0.0.0.0")),
			MethodColor(fmt.Sprintf("%-7s", safeGetValue(entry, "Method", "ERROR"))),
			safeGetValue(entry, "Path", "/error"),
		)
	case GormLogType:
		filepath := strings.SplitN(safeGetValue(entry, "FileWithLineNum", "/unknown/path:-1"), "/CBCTF/", 2)
		_, _ = fmt.Fprintf(ret, "%s | %s rows %s | %s",
			fmt.Sprintf("%-36s", filepath[len(filepath)-1]),
			colors["Debug"](safeGetValue(entry, "Rows", "-1")),
			colors["Debug"](safeGetValue(entry, "Duration", "0.000ms")),
			safeGetValue(entry, "SQL", "SELECT 'unknown'"),
		)
	}
	ret.WriteByte('\n')
	return ret.Bytes(), nil
}

type TextFormatter struct{}

func (f TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	t, ok := entry.Data["Type"].(string)
	if !ok {
		t = DefaultLogType
	}
	LevelText := fmt.Sprintf("%-12s", t+"-"+entry.Level.String())
	ret := new(bytes.Buffer)
	_, _ = fmt.Fprintf(ret, "%s %s | ", LevelText, entry.Time.Format("2006-01-02 15:04:05"))
	switch t {
	case DefaultLogType:
		filepath := strings.SplitN(entry.Caller.File, "/CBCTF/", 2)
		caller := fmt.Sprintf("%s:%d", filepath[len(filepath)-1], entry.Caller.Line)
		caller = fmt.Sprintf("%-36s", caller)
		_, _ = fmt.Fprintf(ret, "%s | %s", caller, entry.Message)
	case TaskLogType:
		_, _ = fmt.Fprintf(ret, "%-36s | %s", "Async Queue Task", entry.Message)
	case GinLogType:
		Latency := safeGetValue(entry, "Latency", time.Duration(0))
		_, _ = fmt.Fprintf(ret, "%s | %d | %13v | ",
			safeGetValue(entry, "TraceID", "00000000-0000-0000-0000-000000000000"),
			safeGetValue(entry, "StatusCode", -1),
			Latency,
		)
		_, _ = fmt.Fprintf(ret, "%s | %s | \"%s\"",
			fmt.Sprintf("%-15s", safeGetValue(entry, "ClientIP", "0.0.0.0")),
			fmt.Sprintf("%-7s", safeGetValue(entry, "Method", "ERROR")),
			safeGetValue(entry, "Path", "/error"),
		)
	case GormLogType:
		filepath := strings.SplitN(safeGetValue(entry, "FileWithLineNum", "/unknown/path:-1"), "/CBCTF/", 2)
		_, _ = fmt.Fprintf(ret, "%s | %s rows %s | %s",
			fmt.Sprintf("%-36s", filepath[len(filepath)-1]),
			safeGetValue(entry, "Rows", "-1"),
			safeGetValue(entry, "Duration", "0.000ms"),
			safeGetValue(entry, "SQL", "SELECT 'unknown'"),
		)
	}
	ret.WriteByte('\n')
	return ret.Bytes(), nil
}
