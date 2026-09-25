package k8s

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"

	"CBCTF/internal/config"
	"CBCTF/internal/model"
	"CBCTF/internal/worker"
)

type workerRoundTrip func(*http.Request) (*http.Response, error)

func (f workerRoundTrip) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestWorkerPublishesOnlyCompletedZipWithoutExec(t *testing.T) {
	previousConfig, previousHTTP := config.Env, workerHTTP
	config.Env = &config.Config{Path: t.TempDir()}
	t.Cleanup(func() { config.Env = previousConfig; workerHTTP = previousHTTP })
	generator := model.Generator{ID: 9, Name: "generator", ChallengeID: 3, WorkerToken: "test-token"}
	pod := &corev1.Pod{
		Name:      generator.Name,
		Namespace: "test",
		Labels:    GeneratorLabels(generator, map[string]string{RoleLabel: GeneratorPodTag}),
		Status: corev1.PodStatus{
			PodIP: "192.0.2.1",
			Phase: corev1.PodRunning,
			Conditions: []corev1.PodCondition{
				{
					Type:   corev1.PodReady,
					Status: corev1.ConditionTrue,
				},
			},
		},
	}
	client := useFakePods(t, pod)
	t.Cleanup(Stop)
	var output bytes.Buffer
	archive := zip.NewWriter(&output)
	file, err := archive.Create("data.txt")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = file.Write([]byte("generated data"))
	if err := archive.Close(); err != nil {
		t.Fatal(err)
	}
	valid := false
	workerHTTP = &http.Client{Transport: workerRoundTrip(func(r *http.Request) (*http.Response, error) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Fatal("worker request has no instance authentication")
		}
		var request worker.Request
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatal(err)
		}
		if request.TeamID != 17 || len(request.Flags) != 1 || request.Flags[0] != "flag{test}" {
			t.Fatal("generator contract changed")
		}
		data := output.Bytes()
		if !valid {
			data = data[:len(data)/2]
		}
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewReader(data))}, nil
	})}
	challenge := model.Challenge{ID: 3}
	flags := []string{"flag{test}"}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := generateAttachment(ctx, challenge, generator, 17, flags); err == nil {
		t.Fatal("incomplete ZIP was accepted")
	}
	path := challenge.AttachmentCachePath(17, flags)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("incomplete artifact became downloadable")
	}
	valid = true
	if err := generateAttachment(ctx, challenge, generator, 17, flags); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil || !bytes.Equal(data, output.Bytes()) {
		t.Fatalf("artifact changed: %v", err)
	}
	for _, action := range client.Actions() {
		if action.GetVerb() != "list" && action.GetVerb() != "watch" {
			t.Fatalf("per-attachment API request: %s %s", action.GetVerb(), action.GetResource())
		}
	}
}
