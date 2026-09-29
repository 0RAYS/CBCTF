package prometheus

import (
	"testing"

	"CBCTF/internal/config"
	"CBCTF/internal/db"
	prom "github.com/prometheus/client_golang/prometheus"
)

func TestUnavailableDatabaseEmitsFailureNotZeroBusinessMetrics(t *testing.T) {
	oldDB, oldConfig := db.DB, config.Env
	db.DB = nil
	config.Env = &config.Config{}
	t.Cleanup(func() { db.DB, config.Env = oldDB, oldConfig })
	registry := prom.NewRegistry()
	registry.MustRegister(NewCTFCollector(), NewPostgresCollector())
	metrics, err := registry.Gather()
	if err != nil {
		t.Fatal(err)
	}
	if len(metrics) != 2 {
		t.Fatalf("unexpected fabricated database metrics: %v", metrics)
	}
	for _, family := range metrics {
		if family.GetName() != "cbctf_ctf_collection_success" && family.GetName() != "postgres_collection_success" {
			t.Fatalf("fake business zero: %s", family.GetName())
		}
		for _, metric := range family.Metric {
			if metric.GetGauge().GetValue() != 0 {
				t.Fatal("database absence reported as successful collection")
			}
		}
	}
}
