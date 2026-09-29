package model

import (
	"context"
	"strings"

	"CBCTF/internal/i18n"
)

type BatchItem struct {
	Details *BatchResult `json:"details,omitempty"`
	ID      string       `json:"id"`
	Phase   string       `json:"phase"`
	Status  string       `json:"status"`
	Code    string       `json:"code,omitempty"`
}

type BatchResult struct {
	Status       string      `json:"status"`
	Requested    int         `json:"requested"`
	Succeeded    int         `json:"succeeded"`
	Skipped      int         `json:"skipped"`
	Failed       int         `json:"failed"`
	NotAttempted int         `json:"not_attempted"`
	Items        []BatchItem `json:"items"`
}

func NewBatch(size int) *BatchResult { return &BatchResult{Requested: size, Items: []BatchItem{}} }
func (b *BatchResult) Success(id, phase string) {
	b.Succeeded++
	b.Items = append(b.Items, BatchItem{ID: id, Phase: phase, Status: "success"})
}
func (b *BatchResult) Skip(id, reason string) {
	b.Skipped++
	b.Items = append(b.Items, BatchItem{ID: id, Phase: reason, Status: "skipped"})
}
func (b *BatchResult) Fail(id, phase string, ret RetVal) {
	b.Failed++
	details, _ := ret.Data.(*BatchResult)
	b.Items = append(b.Items, BatchItem{ID: id, Phase: phase, Status: "failed", Code: ret.Msg, Details: details})
}

func (b *BatchResult) Record(id, phase string, ret RetVal) {
	if ret.OK {
		b.Success(id, phase)
		b.Items[len(b.Items)-1].Details, _ = ret.Data.(*BatchResult)
	} else {
		b.Fail(id, phase, ret)
	}
}
func (b *BatchResult) Result(ctx context.Context) RetVal {
	b.NotAttempted = max(0, b.Requested-b.Succeeded-b.Skipped-b.Failed)
	b.Status = "success"
	if b.Failed > 0 || b.NotAttempted > 0 {
		b.Status = "failed"
		if b.Succeeded > 0 {
			b.Status = "partial"
		}
	}
	if ctx.Err() != nil {
		b.Status = "cancelled"
	}
	if b.Status == "success" {
		return SuccessRetVal(b)
	}
	return RetVal{Data: b, Msg: i18n.Common.BatchIncomplete, Attr: map[string]any{"Succeeded": b.Succeeded, "Failed": b.Failed, "Skipped": b.Skipped, "NotAttempted": b.NotAttempted}}
}

// Database/queue read-write failures invalidate later decisions. Domain errors
// (not found, already running, quota, validation) remain per-object outcomes.
func BatchDependencyFailed(ret RetVal) bool {
	if batch, ok := ret.Data.(*BatchResult); ok {
		if batch.Status == "cancelled" {
			return true
		}
		for _, item := range batch.Items {
			if item.Status == "failed" {
				nested := RetVal{Msg: item.Code}
				if item.Details != nil {
					nested.Data = item.Details
				}
				if BatchDependencyFailed(nested) {
					return true
				}
			}
		}
	}
	if strings.HasPrefix(ret.Msg, "model.") && (strings.HasSuffix(ret.Msg, ".createError") || strings.HasSuffix(ret.Msg, ".updateError") || strings.HasSuffix(ret.Msg, ".deleteError")) {
		return true
	}
	return strings.HasSuffix(ret.Msg, ".getError") || ret.Msg == i18n.Common.UnknownError || ret.Msg == i18n.Model.CreateError || ret.Msg == i18n.Model.UpdateError || ret.Msg == i18n.Model.DeleteError || ret.Msg == i18n.Task.EnqueueError ||
		(strings.HasPrefix(ret.Msg, "redis.") && ret.Msg != i18n.Redis.NotFound && ret.Msg != i18n.Redis.NoAvailablePort)
}
