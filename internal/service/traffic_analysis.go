package service

import (
	"context"
	"time"

	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/traffic"
)

type TrafficAnalysisResult struct {
	Report   *traffic.AnalysisReport `json:"report"`
	Accesses []traffic.TrafficAccess `json:"accesses"`
	Archived bool                    `json:"archived"`
}

func GetTrafficAnalysis(ctx context.Context, victim model.Victim) (TrafficAnalysisResult, model.RetVal) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	repo := db.InitTrafficRepo(db.DB.WithContext(ctx))
	record, ret := repo.GetAnalysis(victim.ID)
	if !ret.OK {
		return TrafficAnalysisResult{}, ret
	}
	if record.ID != 0 {
		return TrafficAnalysisResult{Report: record.Analysis, Accesses: record.Accesses, Archived: true}, model.SuccessRetVal()
	}
	flags, ret := repo.KnownFlags(victim)
	if !ret.OK {
		return TrafficAnalysisResult{}, ret
	}
	report, err := traffic.AnalyzeDir(ctx, victim.TrafficBasePath(), flags)
	if err != nil {
		return TrafficAnalysisResult{}, model.RetVal{Msg: i18n.Model.File.ReadPcapError, Attr: map[string]any{"Error": err.Error()}}
	}
	result, err := traffic.ReadPcapDirWithContext(ctx, victim.TrafficBasePath())
	if err != nil {
		return TrafficAnalysisResult{}, model.RetVal{Msg: i18n.Model.File.ReadPcapError, Attr: map[string]any{"Error": err.Error()}}
	}
	return TrafficAnalysisResult{Report: &report, Accesses: traffic.CollectTrafficAccesses(result, victim.TrafficInternalIPs())}, model.SuccessRetVal()
}
