package k8s

import (
	"CBCTF/internal/log"
	"context"
	"fmt"
	"time"

	authorizationv1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type permissionCheck struct {
	Group     string
	Resource  string
	Verb      string
	Namespace string
}

func buildPermissionChecks() []permissionCheck {
	ns := globalNamespace
	return []permissionCheck{
		// Core: Pods
		{Resource: "pods", Verb: "create", Namespace: ns},
		{Resource: "pods", Verb: "get", Namespace: ns},
		{Resource: "pods", Verb: "list", Namespace: ns},
		{Resource: "pods", Verb: "delete", Namespace: ns},
		{Resource: "pods", Verb: "deletecollection", Namespace: ns},
		{Resource: "pods/exec", Verb: "create", Namespace: ns},
		// Core: Services
		{Resource: "services", Verb: "create", Namespace: ns},
		{Resource: "services", Verb: "list", Namespace: ns},
		{Resource: "services", Verb: "delete", Namespace: ns},
		// Core: ConfigMaps
		{Resource: "configmaps", Verb: "create", Namespace: ns},
		{Resource: "configmaps", Verb: "deletecollection", Namespace: ns},
		// Core: PersistentVolumeClaims
		{Resource: "persistentvolumeclaims", Verb: "get", Namespace: ns},
		// Core: Namespaces (集群级别)
		{Resource: "namespaces", Verb: "get"},
		// Core: Nodes (集群级别)
		{Resource: "nodes", Verb: "list"},
		// batch: Jobs
		{Group: "batch", Resource: "jobs", Verb: "create", Namespace: ns},
		// networking.k8s.io: NetworkPolicies
		{Group: "networking.k8s.io", Resource: "networkpolicies", Verb: "create", Namespace: ns},
		{Group: "networking.k8s.io", Resource: "networkpolicies", Verb: "deletecollection", Namespace: ns},
		// discovery.k8s.io: EndpointSlices
		{Group: "discovery.k8s.io", Resource: "endpointslices", Verb: "deletecollection", Namespace: ns},
		// Multus: NetworkAttachmentDefinitions
		{Group: "k8s.cni.cncf.io", Resource: "network-attachment-definitions", Verb: "create", Namespace: ns},
		{Group: "k8s.cni.cncf.io", Resource: "network-attachment-definitions", Verb: "create", Namespace: "kube-system"},
		{Group: "k8s.cni.cncf.io", Resource: "network-attachment-definitions", Verb: "get", Namespace: "kube-system"},
		{Group: "k8s.cni.cncf.io", Resource: "network-attachment-definitions", Verb: "deletecollection", Namespace: ns},
		// KubeOVN: Subnets (集群级别)
		{Group: "kubeovn.io", Resource: "subnets", Verb: "create"},
		{Group: "kubeovn.io", Resource: "subnets", Verb: "get"},
		{Group: "kubeovn.io", Resource: "subnets", Verb: "deletecollection"},
		// KubeOVN: VPCs (集群级别)
		{Group: "kubeovn.io", Resource: "vpcs", Verb: "create"},
		{Group: "kubeovn.io", Resource: "vpcs", Verb: "deletecollection"},
		// KubeOVN: IPs (集群级别)
		{Group: "kubeovn.io", Resource: "ips", Verb: "deletecollection"},
		// KubeOVN: IptablesEIPs (集群级别)
		{Group: "kubeovn.io", Resource: "iptables-eips", Verb: "create"},
		{Group: "kubeovn.io", Resource: "iptables-eips", Verb: "get"},
		{Group: "kubeovn.io", Resource: "iptables-eips", Verb: "deletecollection"},
		// KubeOVN: IptablesDnatRules (集群级别)
		{Group: "kubeovn.io", Resource: "iptables-dnat-rules", Verb: "create"},
		{Group: "kubeovn.io", Resource: "iptables-dnat-rules", Verb: "deletecollection"},
		// KubeOVN: IptablesSnatRules (集群级别)
		{Group: "kubeovn.io", Resource: "iptables-snat-rules", Verb: "create"},
		{Group: "kubeovn.io", Resource: "iptables-snat-rules", Verb: "deletecollection"},
		// KubeOVN: VpcNatGateways (集群级别)
		{Group: "kubeovn.io", Resource: "vpc-nat-gateways", Verb: "create"},
		{Group: "kubeovn.io", Resource: "vpc-nat-gateways", Verb: "deletecollection"},
	}
}

func checkPermissions() {
	checks := buildPermissionChecks()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var missing []string
	for _, check := range checks {
		sar := &authorizationv1.SelfSubjectAccessReview{
			Spec: authorizationv1.SelfSubjectAccessReviewSpec{
				ResourceAttributes: &authorizationv1.ResourceAttributes{
					Namespace: check.Namespace,
					Verb:      check.Verb,
					Group:     check.Group,
					Resource:  check.Resource,
				},
			},
		}
		result, err := kubeClient.AuthorizationV1().SelfSubjectAccessReviews().Create(ctx, sar, metav1.CreateOptions{})
		if err != nil {
			log.Logger.Fatalf("Failed to check K8s permission (%s %s/%s): %s", check.Verb, check.Group, check.Resource, err)
		}
		if !result.Status.Allowed {
			ns := check.Namespace
			if ns == "" {
				ns = "(cluster-scoped)"
			}
			missing = append(missing, fmt.Sprintf("%s %s/%s [%s]", check.Verb, check.Group, check.Resource, ns))
		}
	}

	if len(missing) > 0 {
		log.Logger.Warningf("Missing %d K8s permissions:", len(missing))
		for _, m := range missing {
			log.Logger.Warningf("  - %s", m)
		}
		log.Logger.Fatal("Insufficient K8s permissions, please check RBAC configuration")
	}
	log.Logger.Infof("K8s permissions verified (%d checks passed)", len(checks))
}
