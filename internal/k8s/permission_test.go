package k8s

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func TestPermissionChecksUseScopedSubresources(t *testing.T) {
	old := globalNamespace
	globalNamespace = "review-ns"
	t.Cleanup(func() { globalNamespace = old })
	found := map[string]bool{}
	for _, check := range buildPermissionChecks() {
		if strings.Contains(check.Resource, "/") {
			t.Errorf("subresource incorrectly embedded: %+v", check)
		}
		if check.Resource == "namespaces" && (check.Name != globalNamespace || check.Namespace != "") {
			t.Errorf("namespace GET not constrained: %+v", check)
		}
		if check.Subresource != "" {
			found[check.Subresource] = true
			if check.Resource != "pods" || check.Namespace != globalNamespace {
				t.Errorf("unscoped subresource: %+v", check)
			}
		}
	}
	if !found["exec"] || !found["log"] {
		t.Fatal("missing pod subresource checks")
	}
}

// No cluster is needed. When Helm is installed, ensure rendered permissions track
// the actual backend self-check contract instead of testing two hand-written lists.
func TestHelmRuntimeRBAC(t *testing.T) {
	helm, err := exec.LookPath("helm")
	if err != nil {
		t.Skip("helm is not installed")
	}
	old := globalNamespace
	globalNamespace = "review-ns"
	t.Cleanup(func() { globalNamespace = old })
	for _, custom := range []bool{false, true} {
		t.Run(map[bool]string{false: "default", true: "existing-service-account"}[custom], func(t *testing.T) {
			args := []string{"template", "review", filepath.Join("..", "..", "chart"), "--namespace", globalNamespace}
			account := "review-cbctf"
			if custom {
				args = append(args, "--set", "serviceAccount.create=false,serviceAccount.name=existing-sa,startupProbe.failureThreshold=90,terminationGracePeriodSeconds=180")
				account = "existing-sa"
			}
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			output, err := exec.CommandContext(ctx, helm, args...).CombinedOutput()
			if err != nil {
				t.Fatalf("helm: %v\n%s", err, output)
			}
			decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(output), 4096)
			var role rbacv1.Role
			var cluster rbacv1.ClusterRole
			var binding rbacv1.RoleBinding
			var clusterBinding rbacv1.ClusterRoleBinding
			var deployment appsv1.Deployment
			for {
				var raw json.RawMessage
				if err := decoder.Decode(&raw); err == io.EOF {
					break
				} else if err != nil {
					t.Fatal(err)
				}
				var meta struct {
					Kind string `json:"kind"`
				}
				if err := json.Unmarshal(raw, &meta); err != nil {
					t.Fatal(err)
				}
				var target any
				switch meta.Kind {
				case "Role":
					target = &role
				case "ClusterRole":
					target = &cluster
				case "RoleBinding":
					target = &binding
				case "ClusterRoleBinding":
					target = &clusterBinding
				case "Deployment":
					var dep appsv1.Deployment
					if err := json.Unmarshal(raw, &dep); err != nil {
						t.Fatal(err)
					}
					if dep.Name == "review-cbctf" {
						deployment = dep
					}
					continue
				default:
					continue
				}
				if err := json.Unmarshal(raw, target); err != nil {
					t.Fatal(err)
				}
			}
			if role.Namespace != globalNamespace || binding.Namespace != globalNamespace || binding.RoleRef.Kind != "Role" || binding.RoleRef.Name != role.Name {
				t.Fatalf("bad Role binding: %+v", binding)
			}
			if clusterBinding.RoleRef.Kind != "ClusterRole" || clusterBinding.RoleRef.Name != cluster.Name {
				t.Fatal("bad cluster binding")
			}
			for _, subjects := range [][]rbacv1.Subject{binding.Subjects, clusterBinding.Subjects} {
				if len(subjects) != 1 || subjects[0].Kind != "ServiceAccount" || subjects[0].Name != account || subjects[0].Namespace != globalNamespace {
					t.Fatalf("bad subjects: %+v", subjects)
				}
			}
			clusterResources := []string{"namespaces", "nodes", "selfsubjectaccessreviews", "subnets", "vpcs", "ips"}
			for _, rule := range cluster.Rules {
				for _, resource := range rule.Resources {
					if !slices.Contains(clusterResources, resource) {
						t.Fatalf("namespaced/broad resource in ClusterRole: %s", resource)
					}
				}
			}
			for _, rules := range [][]rbacv1.PolicyRule{role.Rules, cluster.Rules} {
				for _, rule := range rules {
					if slices.Contains(rule.Verbs, "*") || slices.Contains(rule.Resources, "*") || slices.Contains(rule.APIGroups, "*") {
						t.Fatalf("wildcard rule: %+v", rule)
					}
				}
			}
			for _, check := range buildPermissionChecks() {
				rules := cluster.Rules
				if check.Namespace != "" {
					rules = role.Rules
				}
				resource := check.Resource
				if check.Subresource != "" {
					resource += "/" + check.Subresource
				}
				found := false
				for _, rule := range rules {
					if slices.Contains(rule.APIGroups, check.Group) && slices.Contains(rule.Resources, resource) && slices.Contains(rule.Verbs, check.Verb) && (len(rule.ResourceNames) == 0 || slices.Contains(rule.ResourceNames, check.Name)) {
						found = true
					}
				}
				if !found {
					t.Errorf("missing rendered permission: %+v", check)
				}
			}
			for _, rule := range cluster.Rules {
				if slices.Contains(rule.Resources, "namespaces") && !slices.Equal(rule.ResourceNames, []string{globalNamespace}) {
					t.Fatal("namespace GET is not name-scoped")
				}
			}
			wantGrace, wantFailures := int64(120), int32(60)
			if custom {
				wantGrace, wantFailures = 180, 90
			}
			pod := deployment.Spec.Template.Spec
			if pod.TerminationGracePeriodSeconds == nil || *pod.TerminationGracePeriodSeconds != wantGrace {
				t.Fatalf("bad shutdown grace: %+v", pod.TerminationGracePeriodSeconds)
			}
			if len(pod.Containers) == 0 || pod.Containers[0].StartupProbe == nil || pod.Containers[0].StartupProbe.FailureThreshold != wantFailures {
				t.Fatal("startup probe is missing/not configurable")
			}
		})
	}
}
