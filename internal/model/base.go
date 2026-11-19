package model

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/plugin/optimisticlock"

	"CBCTF/internal/i18n"
)

type BaseModel struct {
	CreatedAt time.Time              `json:"-"`
	UpdatedAt time.Time              `json:"-"`
	DeletedAt gorm.DeletedAt         `gorm:"index" json:"-"`
	Version   optimisticlock.Version `json:"-"`
	ID        uint                   `gorm:"primaryKey" json:"id"`
}

type Model interface {
	GetBaseModel() BaseModel
}

func (b BaseModel) GetBaseModel() BaseModel {
	return b
}

// Shared RetVal attribute names. Every failure records the resource kind and,
// when known, the concrete target it looked up, so logs, wrapped errors and API
// responses describe the same context.
const (
	// AttrModel holds the resource kind an operation targeted, e.g. "Pod".
	AttrModel = "Model"
	// AttrTarget holds the concrete resource an operation looked up, e.g. a
	// Pod name or a database id.
	AttrTarget = "Target"
	// AttrError holds the underlying cause of a failure.
	AttrError = "Error"
)

type RetVal struct {
	Data any
	Attr map[string]any
	Msg  string
	OK   bool
	// Err keeps the underlying cause intact when the failure is wrapped, so
	// errors.Is and errors.As still work.
	Err error `json:"-"`
}

func SuccessRetVal(data ...any) RetVal {
	if len(data) > 0 {
		if len(data) == 1 {
			return RetVal{OK: true, Msg: i18n.Common.Success, Attr: nil, Data: data[0]}
		}
		return RetVal{OK: true, Msg: i18n.Common.Success, Attr: nil, Data: data}
	}
	return RetVal{OK: true, Msg: i18n.Common.Success, Attr: nil, Data: nil}
}

// FailRetVal builds a failed RetVal recording the resource kind, the looked-up
// target and the underlying cause under the shared attribute names.
func FailRetVal(msg, kind, target string, err error) RetVal {
	attr := make(map[string]any, 3)
	if kind != "" {
		attr[AttrModel] = kind
	}
	if target != "" {
		attr[AttrTarget] = target
	}
	// The error attribute is always defined so that a message template
	// referencing it renders an empty cause instead of a missing placeholder.
	if err != nil {
		attr[AttrError] = err.Error()
	} else {
		attr[AttrError] = ""
	}
	return RetVal{Msg: msg, Attr: attr, Err: err}
}

// WithAttrs returns a copy of the failure carrying extra template attributes
// on top of the shared context.
func (r RetVal) WithAttrs(extra map[string]any) RetVal {
	if len(extra) == 0 {
		return r
	}
	attr := make(map[string]any, len(r.Attr)+len(extra))
	for key, value := range r.Attr {
		attr[key] = value
	}
	for key, value := range extra {
		attr[key] = value
	}
	r.Attr = attr
	return r
}

// Error renders a failure as a single line for logs and wrapped errors: the
// message key followed by the resource kind, the target and the cause.
func (r RetVal) Error() string {
	if r.OK {
		return r.Msg
	}
	var sb strings.Builder
	sb.WriteString(r.Msg)
	if kind, ok := r.Attr[AttrModel].(string); ok && kind != "" {
		sb.WriteString(" model=")
		sb.WriteString(kind)
	}
	if target, ok := r.Attr[AttrTarget].(string); ok && target != "" {
		sb.WriteString(" target=")
		sb.WriteString(target)
	}
	if cause, ok := r.Attr[AttrError].(string); ok && cause != "" {
		sb.WriteString(" error=")
		sb.WriteString(cause)
	}
	return sb.String()
}

// Unwrap exposes the underlying cause to errors.Is and errors.As.
func (r RetVal) Unwrap() error {
	return r.Err
}

// Cause returns the underlying cause of the failure: the wrapped error, the
// recorded error detail or, as a last resort, the message key itself.
func (r RetVal) Cause() error {
	if r.Err != nil {
		return r.Err
	}
	if detail, ok := r.Attr[AttrError].(string); ok && detail != "" {
		return errors.New(detail)
	}
	return errors.New(r.Msg)
}
