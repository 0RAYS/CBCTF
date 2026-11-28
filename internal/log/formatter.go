package log

import (
	"bytes"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

// timeLayout keeps the local offset so lines from different hosts stay comparable.
const timeLayout = "2006-01-02 15:04:05 -07:00"

// Formatter renders one log entry as a single line.
type Formatter struct {
	theme theme
}

// NewFormatter builds a formatter that always emits ANSI colors: the admin log
// page and the container output both render them.
func NewFormatter() Formatter {
	return Formatter{theme: newTheme()}
}

type theme struct {
	warning func(a ...any) string
	panic   func(a ...any) string
	fatal   func(a ...any) string
	error   func(a ...any) string
	info    func(a ...any) string
	debug   func(a ...any) string
	methods map[string]func(a ...any) string
}

func newTheme() theme {
	paint := func(attrs ...color.Attribute) func(a ...any) string {
		c := color.New(attrs...)
		c.EnableColor()
		return c.SprintFunc()
	}
	return theme{
		warning: paint(color.FgYellow, color.Bold),
		panic:   paint(color.FgRed, color.BgWhite),
		fatal:   paint(color.FgWhite, color.BgRed, color.Bold),
		error:   paint(color.FgRed, color.Bold),
		info:    paint(color.FgGreen, color.Bold),
		debug:   paint(color.FgBlue, color.Bold),
		methods: map[string]func(a ...any) string{
			http.MethodGet:     paint(color.BgBlue),
			http.MethodPost:    paint(color.BgCyan),
			http.MethodPut:     paint(color.BgYellow),
			http.MethodDelete:  paint(color.BgRed),
			http.MethodPatch:   paint(color.BgGreen),
			http.MethodHead:    paint(color.BgMagenta),
			http.MethodOptions: paint(color.BgWhite, color.FgBlack),
		},
	}
}

func (t theme) levelColor(level logrus.Level) func(a ...any) string {
	switch level {
	case logrus.DebugLevel, logrus.TraceLevel:
		return t.debug
	case logrus.ErrorLevel:
		return t.error
	case logrus.WarnLevel:
		return t.warning
	case logrus.InfoLevel:
		return t.info
	case logrus.PanicLevel:
		return t.panic
	case logrus.FatalLevel:
		return t.fatal
	default:
		return t.info
	}
}

func (t theme) methodColor(method string) func(a ...any) string {
	if c, ok := t.methods[method]; ok {
		return c
	}
	return func(a ...any) string { return fmt.Sprint(a...) }
}

func (t theme) statusCodeColor(code int) func(a ...any) string {
	switch {
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		return t.info
	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
		return t.debug
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		return t.warning
	default:
		return t.error
	}
}

func safeGetValue[T any](entry *logrus.Entry, key string, defaultV T) T {
	V, ok := entry.Data[key].(T)
	if !ok {
		return defaultV
	}
	return V
}

// extraFields renders the structured fields the layout of the entry type does not consume.
func extraFields(entry *logrus.Entry) []string {
	keys := make([]string, 0, len(entry.Data))
	for key := range entry.Data {
		if key == "Type" {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)
	fields := make([]string, 0, len(keys))
	for _, key := range keys {
		fields = append(fields, fmt.Sprintf("%s=%v", key, entry.Data[key]))
	}
	return fields
}

// shortFile keeps the last two path segments, e.g. "internal/db/base.go:119".
func shortFile(fileWithLine string) string {
	path := fileWithLine
	line := ""
	if idx := strings.LastIndex(fileWithLine, ":"); idx >= 0 {
		path = fileWithLine[:idx]
		line = fileWithLine[idx+1:]
	}
	segments := strings.FieldsFunc(path, func(r rune) bool { return r == '/' || r == '\\' })
	if len(segments) > 2 {
		segments = segments[len(segments)-2:]
	}
	short := strings.Join(segments, "/")
	if line == "" {
		return short
	}
	return short + ":" + line
}

func (f Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	t, ok := entry.Data["Type"].(string)
	if !ok {
		t = DefaultLogType
	}
	LevelColor := f.theme.levelColor(entry.Level)
	LevelText := fmt.Sprintf("%-12s", t+"-"+entry.Level.String())
	ret := new(bytes.Buffer)
	_, _ = fmt.Fprintf(ret, "%s %s | ", LevelColor(LevelText), entry.Time.Format(timeLayout))
	switch t {
	case DefaultLogType:
		var callerStr string
		if entry.Caller != nil {
			callerStr = shortFile(fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line))
		} else {
			callerStr = "unknown:0"
		}
		caller := fmt.Sprintf("%-36s", callerStr)
		msg := entry.Message
		if entry.Level <= logrus.ErrorLevel {
			msg = fmt.Sprintf("[%s] %s", strings.ToUpper(entry.Level.String()), entry.Message)
		}
		_, _ = fmt.Fprintf(ret, "%s | %s", caller, LevelColor(msg))
		if fields := extraFields(entry); len(fields) > 0 {
			_, _ = fmt.Fprintf(ret, " | %s", strings.Join(fields, " "))
		}
	case TaskLogType:
		_, _ = fmt.Fprintf(ret, "%-36s | %s", fmt.Sprintf("Async %s Task", entry.Data["Queue"]), LevelColor(entry.Message))
	case GinLogType:
		StatusCodeColor := f.theme.statusCodeColor(safeGetValue(entry, "StatusCode", -1))
		MethodColor := f.theme.methodColor(safeGetValue(entry, "Method", "ERROR"))
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
		if entry.Message != "" {
			_, _ = fmt.Fprintf(ret, " | %s", LevelColor(entry.Message))
		}
		if stack, ok := entry.Data["Stack"].(string); ok && stack != "" {
			_, _ = fmt.Fprintf(ret, "\n%s", stack)
		}
	case GormLogType:
		_, _ = fmt.Fprintf(ret, "%s | %s rows %s | %s",
			fmt.Sprintf("%-36s", shortFile(safeGetValue(entry, "FileWithLineNum", "/unknown/path:-1"))),
			f.theme.debug(safeGetValue(entry, "Rows", "-1")),
			f.theme.debug(safeGetValue(entry, "Duration", "0.000ms")),
			safeGetValue(entry, "SQL", "SELECT 'unknown'"),
		)
	}
	ret.WriteByte('\n')
	return ret.Bytes(), nil
}
