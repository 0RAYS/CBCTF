package k8s

import (
	"context"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestVictimOwnersSeparateNetworkLifetime(t *testing.T) {
	useFakePods(t)
	workloads := withResourceOwner(context.Background(), metav1.OwnerReference{APIVersion: "v1", Kind: "ConfigMap", Name: "victim-7-workloads", UID: "workload-root"})
	network := withResourceOwner(context.Background(), metav1.OwnerReference{APIVersion: "v1", Kind: "ConfigMap", Name: "victim-7-network", UID: "network-root"})
	cm, ret := CreateConfigMap(workloads, CreateConfigMapOptions{Name: "files"})
	if !ret.OK {
		t.Fatal(ret)
	}
	pod, ret := CreatePod(workloads, CreatePodOptions{Name: "web", SubmitOnly: true})
	if !ret.OK {
		t.Fatal(ret)
	}
	policy, ret := CreateNetworkPolicy(network, CreateNetworkPolicyOptions{Name: "policy"})
	if !ret.OK {
		t.Fatal(ret)
	}
	for _, obj := range []metav1.Object{cm, pod, policy} {
		refs := obj.GetOwnerReferences()
		if len(refs) != 1 || refs[0].BlockOwnerDeletion == nil || !*refs[0].BlockOwnerDeletion {
			t.Fatalf("non-blocking owner on %s", obj.GetName())
		}
		want := "workload-root"
		if obj.GetName() == "policy" {
			want = "network-root"
		}
		if string(refs[0].UID) != want {
			t.Fatalf("wrong lifetime on %s: %+v", obj.GetName(), refs)
		}
	}
	if len(resourceOwners(context.Background())) != 0 {
		t.Fatal("owner leaked outside resource tree")
	}
}
