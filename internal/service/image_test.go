package service

import (
	"slices"
	"testing"

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
