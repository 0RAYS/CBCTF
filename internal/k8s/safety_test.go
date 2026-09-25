package k8s

import (
	"context"
	"strings"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	ktesting "k8s.io/client-go/testing"

	"CBCTF/internal/model"
)

func TestCollectionDeletionRequiresOwner(t *testing.T) {
	operations := map[string]func(context.Context, map[string]string) model.RetVal{
		"Pod":           DeletePodCollection,
		"Service":       DeleteServiceCollection,
		"EndpointSlice": DeleteEndpointCollection,
		"ConfigMap":     DeleteConfigMapCollection,
		"NetworkPolicy": DeleteNetworkPolicyCollection,
		"Subnet":        DeleteSubnetCollection,
		"VPC":           DeleteVPCCollection,
		"IP":            DeleteIPCollection,
		"VM":            DeleteVMCollection,
		"NAD": func(ctx context.Context, labels map[string]string) model.RetVal {
			return DeleteNetAttachDefCollection(ctx, "test", labels)
		},
	}
	for name, operation := range operations {
		t.Run(name, func(t *testing.T) {
			// These must return before touching a client. Several clients are nil in tests.
			for _, filters := range []map[string]string{
				nil,
				{},
				{"role": "victim"},
				{"victim_id": ""},
				{"victim_id": "0"},
				{"victim_id": "-1"},
				{"victim_id": "7", "bad key": "x"},
				{"victim_id": "7,role=other"},
			} {
				if ret := operation(context.Background(), filters); ret.OK {
					t.Fatalf("unsafe selector accepted: %v", filters)
				}
			}
		})
	}
	for _, tc := range []struct {
		resource string
		labels   map[string]string
		want     string
	}{
		{"Pod", map[string]string{"victim_id": "7", "role": "victim"}, "role=victim,victim_id=7"},
		{"Service", map[string]string{"generator_id": "8"}, "generator_id=8"},
		{"IP", map[string]string{"ovn.kubernetes.io/subnet": "victim-network"}, "ovn.kubernetes.io/subnet=victim-network"},
	} {
		options, ret := deleteCollectionOptions(tc.resource, tc.labels)
		if !ret.OK || options.LabelSelector != tc.want {
			t.Fatalf("valid selector rejected: %+v %+v", options, ret)
		}
	}
	if ret := DeleteNetAttachDefCollection(context.Background(), "", map[string]string{"victim_id": "7"}); ret.OK {
		t.Fatal("empty namespace accepted")
	}
}

func TestGeneratorCleanupRejectsForeignPod(t *testing.T) {
	pod := &corev1.Pod{
		Name:      "collision",
		Namespace: "test",
		UID:       "other",
		Labels:    map[string]string{"challenge_id": "9", "generator_id": "8", RoleLabel: GeneratorPodTag},
	}
	client := useFakePods(t, pod)
	ret := StopGenerator(context.Background(), model.Generator{ID: 7, ChallengeID: 9, Name: pod.Name})
	if ret.OK {
		t.Fatal("foreign workload accepted")
	}
	for _, action := range client.Actions() {
		if action.GetVerb() != "get" {
			t.Fatalf("foreign workload mutated: %s", action.GetVerb())
		}
	}
}

func TestGeneratorCleanupRejectsMissingOwnerLabel(t *testing.T) {
	pod := &corev1.Pod{
		Name:      "unlabelled",
		Namespace: "test",
		UID:       "old",
		Labels:    map[string]string{"challenge_id": "9", RoleLabel: GeneratorPodTag},
	}
	client := useFakePods(t, pod)
	generator := model.Generator{ID: 7, ChallengeID: 9, Name: pod.Name}
	if ret := StopGenerator(context.Background(), generator); ret.OK {
		t.Fatal("generator without owner label accepted for cleanup")
	}
	for _, action := range client.Actions() {
		if action.GetVerb() != "get" {
			t.Fatalf("unowned workload mutated: %s", action.GetVerb())
		}
	}
}

func TestPodDeletionRequiresObservedIdentity(t *testing.T) {
	client := useFakePods(t)
	for _, deletePod := range []func(context.Context, string, types.UID) model.RetVal{DeletePod, DeletePodAndWait} {
		for _, identity := range [][2]string{{"pod", ""}, {"", "uid"}, {"", ""}} {
			if ret := deletePod(context.Background(), identity[0], types.UID(identity[1])); ret.OK {
				t.Fatalf("missing identity accepted: %v", identity)
			}
		}
	}
	if len(client.Actions()) != 0 {
		t.Fatal("invalid deletion reached Kubernetes API")
	}
}

func TestPodDeletionPinsUIDAndRejectsReplacement(t *testing.T) {
	pod := &corev1.Pod{Name: "owned", Namespace: "test", UID: "old"}
	client := useFakePods(t, pod)
	checked := false
	client.PrependReactor("delete", "pods", func(action ktesting.Action) (bool, runtime.Object, error) {
		options := action.(ktesting.DeleteAction).GetDeleteOptions()
		if options.Preconditions == nil || options.Preconditions.UID == nil || *options.Preconditions.UID != pod.UID {
			t.Fatal("missing UID precondition")
		}
		checked = true
		return false, nil, nil
	})
	client.PrependReactor("get", "pods", func(ktesting.Action) (bool, runtime.Object, error) {
		replacement := pod.DeepCopy()
		replacement.UID = "new"
		return true, replacement, nil
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if ret := DeletePodAndWait(ctx, pod.Name, pod.UID); ret.OK || !checked {
		t.Fatalf("replacement accepted: %+v", ret)
	}
}

func TestImagePullJobHasBoundedLifetimeAndNoCredentials(t *testing.T) {
	useFakePods(t)
	job, ret := CreateJob(context.Background(), CreateJobOptions{
		Name:       "pull-test",
		Images:     []string{"test-image"},
		PullPolicy: "Always",
	})
	if !ret.OK {
		t.Fatal(ret)
	}
	pod := job.Spec.Template.Spec
	if pod.AutomountServiceAccountToken == nil || *pod.AutomountServiceAccountToken ||
		pod.EnableServiceLinks == nil || *pod.EnableServiceLinks {
		t.Fatal("image-pull Pod exposes credentials/service env")
	}
	if job.Spec.ActiveDeadlineSeconds == nil || *job.Spec.ActiveDeadlineSeconds != 600 ||
		job.Spec.TTLSecondsAfterFinished == nil || *job.Spec.TTLSecondsAfterFinished != 300 {
		t.Fatal("job lifetime is unbounded")
	}
}

func TestClientConfigurationSharesRateLimiter(t *testing.T) {
	config := &rest.Config{}
	configureClientRateLimit(config)
	if config.RateLimiter == nil || rest.CopyConfig(config).RateLimiter != config.RateLimiter ||
		config.RateLimiter.QPS() != 100 || config.Burst != 150 {
		t.Fatal("missing shared API request budget")
	}
	if !strings.Contains(config.UserAgent, "cbctf") {
		t.Fatal("missing identifiable client user-agent")
	}
}
