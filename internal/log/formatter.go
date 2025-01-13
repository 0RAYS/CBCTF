package log

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
	"time"
)

type Formatter struct{}

var defaultFormatter = logrus.TextFormatter{ForceColors: true}

func safeGetValue[T any](entry *logrus.Entry, key string) T {
	V, _ := entry.Data[key].(T)
	return V
}

var colors = map[string]func(a ...interface{}) string{
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

func MethodColor(method string) func(a ...interface{}) string {
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

func levelColor(level logrus.Level) func(a ...interface{}) string {
	switch level {
	case logrus.DebugLevel, logrus.TraceLevel:
		return colors["Debug"]
	case logrus.ErrorLevel:
		return colors["Error"]
	case logrus.WarnLevel:
		return colors["Warning"]
	case logrus.InfoLevel:
		return colors["Info"]
	default:
		return colors["Time"]
	}
}

func StatusCodeColor(code int) func(a ...interface{}) string {
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

func FormatTraceIdIfExist(entry *logrus.Entry) string {
	uu, ok := entry.Data["traceID"]
	if !ok {
		return ""
	}
	return fmt.Sprint(uu)
}

func (f Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	Type, ok := entry.Data["type"].(string)
	if !ok {
		return defaultFormatter.Format(entry)
	}
	LevelColor := levelColor(entry.Level)
	LevelText := strings.ToUpper(entry.Level.String())
	Type = strings.ToUpper(Type)
	ret := new(bytes.Buffer)
	base, _ := os.Getwd()
	base = strings.ReplaceAll(base, "\\", "/")
	switch Type {
	case "GORM":
		_, _ = fmt.Fprintf(ret, "%s  %s",
			LevelColor(Type+"-"+LevelText), entry.Time.Format("2006/01/02 15:04:05"),
		)
		_, _ = fmt.Fprintf(ret, "| %s rows %s %s | %s",
			colors["Debug"](safeGetValue[string](entry, "rows")),
			colors["Debug"](safeGetValue[string](entry, "duration")),
			strings.Replace(safeGetValue[string](entry, "fileWithLineNum"), base, "", 1),
			safeGetValue[string](entry, "sql"),
		)
	case "GIN":
		_, _ = fmt.Fprintf(ret, "%s  %s",
			LevelColor(Type+"-"+LevelText+" "), safeGetValue[string](entry, "timeStamp"),
		)
		method := safeGetValue[string](entry, "method")
		statueColor := StatusCodeColor(safeGetValue[int](entry, "statusCode"))
		methodColor := MethodColor(method)
		Latency := safeGetValue[time.Duration](entry, "Latency")
		if Latency > time.Minute {
			Latency = Latency.Truncate(time.Second)
		}
		_, _ = fmt.Fprintf(ret, "| %s |%s| %13v | %15s |%s %#v",
			FormatTraceIdIfExist(entry),
			statueColor(fmt.Sprintf(" %d ", safeGetValue[int](entry, "statusCode"))),
			Latency,
			safeGetValue[string](entry, "clientIP"),
			methodColor(fmt.Sprintf(" %s ", method)),
			safeGetValue[string](entry, "path"),
		)
		if entry.Message != "" {
			_, _ = fmt.Fprintf(ret, "  %s", entry.Message)
		}
	}
	ret.WriteByte('\n')
	return ret.Bytes(), nil
}
