package log

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"
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

func safeGetValue[T any](entry *logrus.Entry, key string, defaultV ...any) T {
	V, err := entry.Data[key].(T)
	if !err {
		if len(defaultV) > 0 {
			var tmp T
			if reflect.TypeOf(tmp) != reflect.TypeOf(defaultV[0]) {
				Logger.Fatalf("type mismatch: want %v, got %v", reflect.TypeOf(tmp), reflect.TypeOf(defaultV[0]))
			}
			return defaultV[0].(T)
		}
		return V
	}
	return V
}

type Formatter struct{}

func (f Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	t, ok := entry.Data["type"].(string)
	if !ok {
		t = "LOG"
	}
	t = strings.ToUpper(t)
	LevelColor := levelColor(entry.Level)
	LevelText := fmt.Sprintf("%-12s", t+"-"+entry.Level.String())
	base, _ := os.Getwd()
	base = strings.ReplaceAll(base, "\\", "/") + "/"
	ret := new(bytes.Buffer)
	switch t {
	case "LOG":
		_, _ = fmt.Fprintf(ret, "%s %s | ",
			LevelColor(LevelText),
			entry.Time.Format("2006-01-02 15:04:05"),
		)
		caller := fmt.Sprintf("%s:%d", strings.Replace(entry.Caller.File, base, "", 1), entry.Caller.Line)
		caller = fmt.Sprintf("%-36s", caller)
		_, _ = fmt.Fprintf(ret, "%s | %s", caller, LevelColor(entry.Message))
	case "GIN":
		StatusCodeColor := statusCodeColor(safeGetValue[int](entry, "StatusCode"))
		MethodColor := methodColor(safeGetValue[string](entry, "Method"))
		Latency := safeGetValue[time.Duration](entry, "Latency")
		if Latency > time.Minute {
			Latency = Latency.Truncate(time.Second)
		}
		_, _ = fmt.Fprintf(ret, "%s %s | ",
			LevelColor(LevelText),
			entry.Time.Format("2006-01-02 15:04:05"),
		)
		_, _ = fmt.Fprintf(ret, "%s | %s | %13v | ",
			safeGetValue[string](entry, "TraceID"),
			StatusCodeColor(safeGetValue[int](entry, "StatusCode")),
			Latency,
		)
		_, _ = fmt.Fprintf(ret, "%s | %s | \"%s\"",
			fmt.Sprintf("%-15s", safeGetValue[string](entry, "ClientIP")),
			MethodColor(fmt.Sprintf("%-7s", safeGetValue[string](entry, "Method"))),
			safeGetValue[string](entry, "Path"),
		)
	case "GORM":
		_, _ = fmt.Fprintf(ret, "%s %s | ",
			LevelColor(LevelText),
			entry.Time.Format("2006-01-02 15:04:05"),
		)
		_, _ = fmt.Fprintf(ret, "%s | %s rows %s | %s",
			fmt.Sprintf("%-36s", strings.Replace(safeGetValue[string](entry, "fileWithLineNum"), base, "", 1)),
			colors["Debug"](safeGetValue[string](entry, "rows")),
			colors["Debug"](safeGetValue[string](entry, "duration")),
			safeGetValue[string](entry, "sql"),
		)
	}
	ret.WriteByte('\n')
	return ret.Bytes(), nil
}

type TextFormatter struct{}

func (f TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	t, ok := entry.Data["type"].(string)
	if !ok {
		t = "LOG"
	}
	t = strings.ToUpper(t)
	LevelText := fmt.Sprintf("%-12s", t+"-"+entry.Level.String())
	base, _ := os.Getwd()
	base = strings.ReplaceAll(base, "\\", "/") + "/"
	ret := new(bytes.Buffer)
	switch t {
	case "LOG":
		_, _ = fmt.Fprintf(ret, "%s %s | ",
			LevelText,
			entry.Time.Format("2006-01-02 15:04:05"),
		)
		caller := fmt.Sprintf("%s:%d", strings.Replace(entry.Caller.File, base, "", 1), entry.Caller.Line)
		caller = fmt.Sprintf("%-36s", caller)
		_, _ = fmt.Fprintf(ret, "%s | %s", caller, entry.Message)
	case "GIN":
		Latency := safeGetValue[time.Duration](entry, "Latency")
		if Latency > time.Minute {
			Latency = Latency.Truncate(time.Second)
		}
		_, _ = fmt.Fprintf(ret, "%s %s | ",
			LevelText,
			entry.Time.Format("2006-01-02 15:04:05"),
		)
		_, _ = fmt.Fprintf(ret, "%s | %d | %13v | ",
			safeGetValue[string](entry, "TraceID"),
			safeGetValue[int](entry, "StatusCode"),
			Latency,
		)
		_, _ = fmt.Fprintf(ret, "%s | %s | \"%s\"",
			fmt.Sprintf("%-15s", safeGetValue[string](entry, "ClientIP")),
			fmt.Sprintf("%-7s", safeGetValue[string](entry, "Method")),
			safeGetValue[string](entry, "Path"),
		)
	case "GORM":
		_, _ = fmt.Fprintf(ret, "%s %s | ", LevelText, entry.Time.Format("2006-01-02 15:04:05"))
		_, _ = fmt.Fprintf(ret, "%s | %s rows %s | %s",
			fmt.Sprintf("%-36s", strings.Replace(safeGetValue[string](entry, "fileWithLineNum"), base, "", 1)),
			safeGetValue[string](entry, "rows"),
			safeGetValue[string](entry, "duration"),
			safeGetValue[string](entry, "sql"),
		)
	}
	ret.WriteByte('\n')
	return ret.Bytes(), nil
}
