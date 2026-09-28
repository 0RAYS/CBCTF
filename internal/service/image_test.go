package service

import (
	"slices"
	"strings"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"CBCTF/internal/config"
	"CBCTF/internal/model"
)

func TestChallengeImageInventoryIncludesWorkerAndEnabledSidecars(t *testing.T) {
	previous := config.Env
	config.Env = &config.Config{}
	t.Cleanup(func() { config.Env = previous })
	config.Env.K8S.WorkerImage = "ghcr.io/0rays/cbctf-worker:v1"
	config.Env.K8S.CaptureImage = "capture:v1"
	config.Env.K8S.Frp.FrpcImage = "frpc:v1"
	config.Env.K8S.Frp.NginxImage = "nginx"
	dynamic := model.Challenge{Type: model.DynamicChallengeType, GeneratorImage: "alpine:3"}
	want := []string{"docker.io/library/alpine:3", "ghcr.io/0rays/cbctf-worker:v1"}
	if images := ChallengeImages(dynamic); !slices.Equal(images, want) {
		t.Fatalf("dynamic images: %v", images)
	}
	pods := model.Challenge{
		Type: model.PodsChallengeType,
		Template: model.ChallengeTemplate{Pods: []model.ChallengePodTemplate{
			{Containers: []model.ChallengeContainerTemplate{{Image: "nginx"}, {Image: "docker.io/library/nginx:latest"}}},
		}},
	}
	if images := ChallengeImages(pods); !slices.Equal(images, []string{"docker.io/library/nginx:latest"}) {
		t.Fatalf("disabled sidecars or duplicate aliases: %v", images)
	}
	config.Env.K8S.CaptureEnabled, config.Env.K8S.Frp.On = true, true
	want = []string{"docker.io/library/capture:v1", "docker.io/library/frpc:v1", "docker.io/library/nginx:latest"}
	if images := ChallengeImages(pods); !slices.Equal(images, want) {
		t.Fatalf("enabled sidecars: %v", images)
	}
}

func TestListChallengeImagesUsesWholeLibraryAndRuntimeDependencies(t *testing.T) {
	previous := config.Env
	config.Env = &config.Config{}
	t.Cleanup(func() { config.Env = previous })
	config.Env.K8S.WorkerImage = "worker:v1"
	config.Env.K8S.CaptureEnabled = true
	config.Env.K8S.CaptureImage = "capture:v1"
	config.Env.K8S.Frp.On = true
	config.Env.K8S.Frp.FrpcImage = "frpc:v1"
	config.Env.K8S.Frp.NginxImage = "nginx"

	tx, err := gorm.Open(postgres.Open("host=127.0.0.1 user=test dbname=test"), &gorm.Config{
		DryRun: true, DisableAutomaticPing: true, Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatal(err)
	}
	pool, err := tx.DB()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = pool.Close() })
	challenges := []model.Challenge{
		{Type: model.DynamicChallengeType, GeneratorImage: "generator:v1"},
		{Type: model.PodsChallengeType, Template: model.ChallengeTemplate{Pods: []model.ChallengePodTemplate{
			{Containers: []model.ChallengeContainerTemplate{{Image: "app:v2"}, {Image: "docker.io/library/nginx:latest"}}},
		}}},
	}
	if err := tx.Callback().Query().After("gorm:query").Register("test:library-images", func(query *gorm.DB) {
		sql := query.Statement.SQL.String()
		if strings.Contains(sql, "contest") || strings.Contains(sql, "LIMIT") || !strings.Contains(sql, `"challenges"."deleted_at" IS NULL`) {
			t.Errorf("expected entire non-deleted challenge library: %s", sql)
		}
		*query.Statement.Dest.(*[]model.Challenge) = challenges
	}); err != nil {
		t.Fatal(err)
	}
	images, ret := ListChallengeImages(tx)
	want := []string{
		"docker.io/library/app:v2", "docker.io/library/capture:v1", "docker.io/library/frpc:v1",
		"docker.io/library/generator:v1", "docker.io/library/nginx:latest", "docker.io/library/worker:v1",
	}
	if !ret.OK || !slices.Equal(images, want) {
		t.Fatalf("images=%v, result=%+v; want %v", images, ret, want)
	}
	challenges = nil
	images, ret = ListChallengeImages(tx)
	want = []string{"docker.io/library/capture:v1", "docker.io/library/frpc:v1", "docker.io/library/nginx:latest", "docker.io/library/worker:v1"}
	if !ret.OK || !slices.Equal(images, want) {
		t.Fatalf("empty library lost runtime dependencies: %v, %+v", images, ret)
	}
	config.Env.K8S.CaptureEnabled, config.Env.K8S.Frp.On = false, false
	images, ret = ListChallengeImages(tx)
	if !ret.OK || !slices.Equal(images, []string{"docker.io/library/worker:v1"}) {
		t.Fatalf("disabled runtime dependencies included: %v, %+v", images, ret)
	}
}
