package traffic

import "slices"

type SourceIssue struct {
	File  string `json:"file"`
	Phase string `json:"phase"`
	Error string `json:"error"`
}

type CaptureFileResult struct {
	File    string `json:"file"`
	Status  string `json:"status"`
	Packets int64  `json:"packets"`
}

func (r *AnalysisReport) AddSourceIssues(issues ...SourceIssue) {
	if len(issues) == 0 {
		return
	}
	r.SourceIssues = append(r.SourceIssues, issues...)
	r.Partial = true
	r.Truncated = true
	if slices.Contains(r.Warnings, "capture_source_incomplete") {
		return
	}
	r.Warnings = append(r.Warnings, "capture_source_incomplete")
}
