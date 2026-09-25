package k8s

import (
	"context"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"

	"CBCTF/internal/model"
)

func TestSharedPodCacheAndVictimReadiness(t *testing.T) {
	pod := &corev1.Pod{
		Name:      "web",
		Namespace: "test",
		UID:       "first",
		Status: corev1.PodStatus{
			Phase:      corev1.PodRunning,
			HostIP:     "192.0.2.10",
			Conditions: []corev1.PodCondition{{Type: corev1.PodReady, Status: corev1.ConditionTrue}},
		},
	}
	client := useFakePods(t, pod)
	t.Cleanup(Stop)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	victim := model.Victim{
		Pods: []model.Pod{{Name: "web"}},
		Resources: model.VictimResources{
			Submitted: true,
			UIDs:      model.StringMap{"web": "first"},
			NodePorts: []model.NodePortEndpoint{
				{
					PodName:  "web",
					Endpoint: model.Endpoint{Port: 30001, Protocol: "TCP"},
				},
			},
		},
	}
	for range 10 {
		ready, err := VictimReady(ctx, &victim)
		if err != nil || !ready {
			t.Fatalf("ready=%t err=%v", ready, err)
		}
		if _, ret := ListPods(ctx); !ret.OK {
			t.Fatal(ret)
		}
	}
	if len(victim.ExposedEndpoints) != 1 || victim.ExposedEndpoints[0].IP != "192.0.2.10" {
		t.Fatalf("bad endpoints: %+v", victim.ExposedEndpoints)
	}
	lists, watches := 0, 0
	for _, action := range client.Actions() {
		switch action.GetVerb() {
		case "list":
			lists++
		case "watch":
			watches++
		}
	}
	if lists != 1 || watches != 1 {
		t.Fatalf("expected one shared list/watch, got %d/%d", lists, watches)
	}
	victim.Resources.UIDs["web"] = "replaced"
	if ready, err := VictimReady(ctx, &victim); ready || err == nil {
		t.Fatal("replacement accepted as ready")
	}
}
