package service

import (
	"context"
	"sort"
	"time"

	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/traffic"
)

type TrafficAnalysisResult struct {
	Report   *traffic.AnalysisReport `json:"report"`
	Accesses []traffic.Access        `json:"accesses"`
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
	if record.Archived && record.Analysis != nil && !record.Analysis.Partial {
		return TrafficAnalysisResult{Report: record.Analysis, Accesses: record.Accesses, Archived: true}, model.SuccessRetVal()
	}
	if record.ID != 0 && time.Since(record.UpdatedAt) < 15*time.Second {
		return TrafficAnalysisResult{Report: record.Analysis, Accesses: record.Accesses, Archived: record.Archived}, model.SuccessRetVal()
	}
	flags, ret := repo.KnownFlags(victim)
	if !ret.OK {
		return TrafficAnalysisResult{}, ret
	}
	report, err := traffic.AnalyzeDir(ctx, victim.TrafficBasePath(), traffic.AnalysisOptions{KnownFlags: flags, InternalIPs: victim.TrafficInternalIPs()})
	if err != nil {
		return TrafficAnalysisResult{}, model.RetVal{Msg: i18n.Model.File.ReadPcapError, Attr: map[string]any{"Error": err.Error()}}
	}
	result, err := traffic.ReadPcapDir(ctx, victim.TrafficBasePath(), victim.TrafficProxyPorts())
	if err != nil {
		return TrafficAnalysisResult{}, model.RetVal{Msg: i18n.Model.File.ReadPcapError, Attr: map[string]any{"Error": err.Error()}}
	}
	accesses := traffic.CollectTrafficAccesses(result, victim.TrafficInternalIPs())
	report.AddSourceIssues(result.SourceIssues...)
	ipSet := make(map[string]bool)
	for _, c := range result.Connections {
		ipSet[c.SrcIP] = true
		ipSet[c.DstIP] = true
	}
	for _, access := range accesses {
		ipSet[access.IP] = true
	}
	ips := make([]string, 0, len(ipSet))
	for ip := range ipSet {
		ips = append(ips, ip)
	}
	sort.Strings(ips)
	if ret = repo.ReplaceAnalysis(victim.ID, ips, accesses, &report, record.Archived); !ret.OK {
		return TrafficAnalysisResult{}, ret
	}
	stored, ret := repo.GetAnalysis(victim.ID)
	if !ret.OK {
		return TrafficAnalysisResult{}, ret
	}
	return TrafficAnalysisResult{Report: stored.Analysis, Accesses: stored.Accesses, Archived: stored.Archived}, model.SuccessRetVal()
}
