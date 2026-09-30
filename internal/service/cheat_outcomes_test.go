package service

import (
	"context"
	"testing"

	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/prometheus"

	metricdto "github.com/prometheus/client_model/go"
)

func TestCheatMetricsCountOnlyNewPersistedEvidence(t *testing.T) {
	counter := prometheus.CheatDetectionsTotal.WithLabelValues("987654321", string(model.ReasonTypeSameVictimIPType))
	metric := &metricdto.Metric{}
	_ = counter.Write(metric)
	before := metric.GetCounter().GetValue()
	batch := model.NewBatch(3)
	if recordCheatOutcome(batch, "failed", model.RetVal{Msg: i18n.Model.Cheat.CreateError}, false, 987654321, model.ReasonTypeSameVictimIPType) {
		t.Fatal("database failure should stop scan")
	}
	recordCheatOutcome(batch, "existing", model.SuccessRetVal(), false, 987654321, model.ReasonTypeSameVictimIPType)
	_ = counter.Write(metric)
	if metric.GetCounter().GetValue() != before {
		t.Fatal("failed/duplicate evidence incremented metrics")
	}
	recordCheatOutcome(batch, "new", model.SuccessRetVal(), true, 987654321, model.ReasonTypeSameVictimIPType)
	_ = counter.Write(metric)
	if metric.GetCounter().GetValue() != before+1 || batch.Failed != 1 || batch.Skipped != 1 || batch.Succeeded != 1 {
		t.Fatalf("wrong evidence accounting: %+v", batch)
	}
	parent := model.NewBatch(1)
	parent.Record("victim_ip", "scan", batch.Result(context.Background()))
	ret := parent.Result(context.Background())
	if ret.OK || !model.BatchDependencyFailed(ret) || parent.Items[0].Details == nil {
		t.Fatal("nested scan lost database failure")
	}
}
