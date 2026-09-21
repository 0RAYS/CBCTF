package k8s

import (
	"CBCTF/internal/model"
	"strconv"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
)

// PodDiagnostics is deliberately an allowlist. In particular, Kubernetes status
// messages/termination messages can contain user-controlled output or secrets.
type PodDiagnostics struct {
	Name              string                 `json:"name"`
	UID               string                 `json:"uid"`
	Status            corev1.PodPhase        `json:"status"`
	Ready             bool                   `json:"ready"`
	Terminating       bool                   `json:"terminating"`
	Reason            string                 `json:"reason,omitempty"`
	CreatedAt         time.Time              `json:"created_at"`
	ContainerStatuses []ContainerDiagnostics `json:"container_statuses"`
	Conditions        []ConditionDiagnostics `json:"conditions"`
}

type ContainerDiagnostics struct {
	Name         string `json:"name"`
	Init         bool   `json:"init"`
	Ready        bool   `json:"ready"`
	State        string `json:"state"`
	Reason       string `json:"reason,omitempty"`
	Restarts     int32  `json:"restarts"`
	ExitCode     *int32 `json:"exit_code,omitempty"`
	LastExitCode *int32 `json:"last_exit_code,omitempty"`
	LastReason   string `json:"last_reason,omitempty"`
}

type ConditionDiagnostics struct {
	Type   corev1.PodConditionType `json:"type"`
	Status corev1.ConditionStatus  `json:"status"`
	Reason string                  `json:"reason,omitempty"`
}

func DescribePod(pod *corev1.Pod) PodDiagnostics {
	result := PodDiagnostics{
		Name: pod.Name, UID: string(pod.UID), Status: pod.Status.Phase,
		Ready: podReady(pod), Terminating: pod.DeletionTimestamp != nil,
		Reason: diagnosticReason(pod.Status.Reason), CreatedAt: pod.CreationTimestamp.Time,
		ContainerStatuses: []ContainerDiagnostics{}, Conditions: []ConditionDiagnostics{},
	}
	appendContainers := func(containers []corev1.Container, statuses []corev1.ContainerStatus, init bool) {
		byName := make(map[string]corev1.ContainerStatus, len(statuses))
		for _, status := range statuses {
			byName[status.Name] = status
		}
		for _, container := range containers {
			if container.Name == CaptureContainerName {
				continue
			}
			status := byName[container.Name]
			item := ContainerDiagnostics{Name: container.Name, Init: init, Ready: status.Ready, Restarts: status.RestartCount, State: "unknown"}
			switch {
			case status.State.Waiting != nil:
				item.State, item.Reason = "waiting", diagnosticReason(status.State.Waiting.Reason)
			case status.State.Running != nil:
				item.State = "running"
			case status.State.Terminated != nil:
				item.State, item.Reason = "terminated", diagnosticReason(status.State.Terminated.Reason)
				item.ExitCode = new(status.State.Terminated.ExitCode)
			}
			if previous := status.LastTerminationState.Terminated; previous != nil {
				item.LastExitCode, item.LastReason = new(previous.ExitCode), diagnosticReason(previous.Reason)
			}
			result.ContainerStatuses = append(result.ContainerStatuses, item)
		}
	}
	appendContainers(pod.Spec.InitContainers, pod.Status.InitContainerStatuses, true)
	appendContainers(pod.Spec.Containers, pod.Status.ContainerStatuses, false)
	for _, condition := range pod.Status.Conditions {
		result.Conditions = append(result.Conditions, ConditionDiagnostics{Type: condition.Type, Status: condition.Status, Reason: diagnosticReason(condition.Reason)})
	}
	return result
}

func diagnosticReason(reason string) string {
	// Reasons are machine-readable codes, not free-form messages.
	return strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '_' || r == '-' {
			return r
		}
		return -1
	}, string([]rune(reason)[:min(len([]rune(reason)), 128)]))
}

func PodMatchesLabels(pod *corev1.Pod, expected map[string]string) bool {
	if pod == nil || len(expected) == 0 {
		return false
	}
	for key, value := range expected {
		if pod.Labels[key] != value {
			return false
		}
	}
	return true
}

func GeneratorOwnsPod(generator model.Generator, pod *corev1.Pod) bool {
	if pod == nil || pod.Name != generator.Name {
		return false
	}
	return PodMatchesLabels(pod, map[string]string{
		"challenge_id": strconv.FormatUint(uint64(generator.ChallengeID), 10),
		"generator_id": strconv.FormatUint(uint64(generator.ID), 10),
		RoleLabel:      GeneratorPodTag,
	})
}
