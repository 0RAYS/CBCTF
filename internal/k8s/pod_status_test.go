package k8s

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"CBCTF/internal/model"
)

func TestPodDiagnosticsAllowlist(t *testing.T) {
	pod := &corev1.Pod{
		Name: "test",
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name: "web",
					Env:  []corev1.EnvVar{{Name: "FLAG", Value: "secret-marker"}},
					Args: []string{"secret-marker"},
				},
				{Name: CaptureContainerName},
			},
			InitContainers: []corev1.Container{{Name: "init"}},
		},
		Status: corev1.PodStatus{
			Phase:   corev1.PodPending,
			Message: "secret-marker",
			Conditions: []corev1.PodCondition{
				{
					Type:    corev1.PodScheduled,
					Status:  corev1.ConditionFalse,
					Reason:  "Unschedulable",
					Message: "secret-marker",
				},
			},
			ContainerStatuses: []corev1.ContainerStatus{
				{
					Name:         "web",
					RestartCount: 3,
					State: corev1.ContainerState{
						Waiting: &corev1.ContainerStateWaiting{
							Reason:  "CrashLoopBackOff",
							Message: "secret-marker",
						},
					},
					LastTerminationState: corev1.ContainerState{
						Terminated: &corev1.ContainerStateTerminated{
							ExitCode: 137,
							Reason:   "OOMKilled",
							Message:  "secret-marker",
						},
					},
				},
			},
		},
	}
	result := DescribePod(pod)
	if result.Ready || len(result.ContainerStatuses) != 2 || !result.ContainerStatuses[0].Init {
		t.Fatalf("bad status: %+v", result)
	}
	web := result.ContainerStatuses[1]
	if web.Reason != "CrashLoopBackOff" || web.Restarts != 3 ||
		web.LastExitCode == nil || *web.LastExitCode != 137 || web.LastReason != "OOMKilled" {
		t.Fatalf("missing detail: %+v", web)
	}
	raw, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(raw, &fields); err != nil {
		t.Fatal(err)
	}
	if _, exists := fields["containers"]; exists {
		t.Fatal("legacy container names field retained")
	}
	if strings.Contains(string(raw), "secret-marker") || strings.Contains(string(raw), CaptureContainerName) {
		t.Fatalf("diagnostic leaked data: %s", raw)
	}
}

func TestGeneratorPodOwnership(t *testing.T) {
	generator := model.Generator{ID: 7, ChallengeID: 9, Name: "generator-7"}
	pod := &corev1.Pod{Name: generator.Name, Labels: map[string]string{"challenge_id": "9", RoleLabel: GeneratorPodTag}}
	if GeneratorOwnsPod(generator, pod) {
		t.Fatal("pod without generator_id accepted")
	}
	pod.Labels["generator_id"] = ""
	if GeneratorOwnsPod(generator, pod) {
		t.Fatal("empty owner accepted")
	}
	pod.Labels["generator_id"] = "8"
	if GeneratorOwnsPod(generator, pod) {
		t.Fatal("wrong owner accepted")
	}
	pod.Labels["generator_id"] = "7"
	if !GeneratorOwnsPod(generator, pod) {
		t.Fatal("owned pod rejected")
	}
	pod.Name = "different"
	if GeneratorOwnsPod(generator, pod) || PodMatchesLabels(pod, nil) {
		t.Fatal("unscoped pod accepted")
	}
}

func TestPodLogsAreBounded(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("tailLines") != "10000" || r.URL.Query().Get("limitBytes") != "1048576" || r.URL.Query().Get("container") != "web" {
			t.Errorf("unbounded request: %s", r.URL)
		}
		// Defense in depth: even a server that ignores limitBytes is bounded locally.
		_, _ = io.WriteString(w, strings.Repeat("x", int(MaxPodLogBytes)+2048))
	}))
	defer server.Close()
	old := kubeClient
	t.Cleanup(func() { kubeClient = old })
	var err error
	kubeClient, err = kubernetes.NewForConfig(&rest.Config{Host: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	result, ret := GetPodLogs(context.Background(), "test", "web", MaxPodLogLines+1)
	if !ret.OK || len(result) != int(MaxPodLogBytes) {
		t.Fatalf("length=%d ret=%+v", len(result), ret)
	}
	for _, lines := range []int64{-1, 0} {
		if *podLogOptions("web", lines).TailLines != 1000 {
			t.Fatal("default line limit missing")
		}
	}
}
