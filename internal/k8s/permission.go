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
		{Resource: "pods/exec", Verb: "create", Namespace: ns},
		// Core: Services
		{Resource: "services", Verb: "create", Namespace: ns},
		{Resource: "services", Verb: "list", Namespace: ns},
		{Resource: "services", Verb: "delete", Namespace: ns},
		// Core: ConfigMaps
		{Resource: "configmaps", Verb: "create", Namespace: ns},
		{Resource: "configmaps", Verb: "get", Namespace: ns},
		{Resource: "configmaps", Verb: "delete", Namespace: ns},
		// Core: PersistentVolumes (集群级别)
		{Resource: "persistentvolumes", Verb: "create"},
		{Resource: "persistentvolumes", Verb: "get"},
		// Core: PersistentVolumeClaims
		{Resource: "persistentvolumeclaims", Verb: "create", Namespace: ns},
		{Resource: "persistentvolumeclaims", Verb: "get", Namespace: ns},
		// Core: Namespaces (集群级别)
		{Resource: "namespaces", Verb: "create"},
		// Core: Nodes (集群级别)
		{Resource: "nodes", Verb: "list"},
		// networking.k8s.io: NetworkPolicies
		{Group: "networking.k8s.io", Resource: "networkpolicies", Verb: "create", Namespace: ns},
		{Group: "networking.k8s.io", Resource: "networkpolicies", Verb: "get", Namespace: ns},
		{Group: "networking.k8s.io", Resource: "networkpolicies", Verb: "delete", Namespace: ns},
		// discovery.k8s.io: EndpointSlices
		{Group: "discovery.k8s.io", Resource: "endpointslices", Verb: "create", Namespace: ns},
		{Group: "discovery.k8s.io", Resource: "endpointslices", Verb: "get", Namespace: ns},
		{Group: "discovery.k8s.io", Resource: "endpointslices", Verb: "delete", Namespace: ns},
		// Multus: NetworkAttachmentDefinitions
		{Group: "k8s.cni.cncf.io", Resource: "network-attachment-definitions", Verb: "create", Namespace: ns},
		{Group: "k8s.cni.cncf.io", Resource: "network-attachment-definitions", Verb: "get", Namespace: ns},
		{Group: "k8s.cni.cncf.io", Resource: "network-attachment-definitions", Verb: "delete", Namespace: ns},
		// KubeOVN: Subnets (集群级别)
		{Group: "kubeovn.io", Resource: "subnets", Verb: "create"},
		{Group: "kubeovn.io", Resource: "subnets", Verb: "get"},
		{Group: "kubeovn.io", Resource: "subnets", Verb: "delete"},
		// KubeOVN: VPCs (集群级别)
		{Group: "kubeovn.io", Resource: "vpcs", Verb: "create"},
		{Group: "kubeovn.io", Resource: "vpcs", Verb: "get"},
		{Group: "kubeovn.io", Resource: "vpcs", Verb: "delete"},
		// KubeOVN: IPs (集群级别)
		{Group: "kubeovn.io", Resource: "ips", Verb: "delete"},
		// KubeOVN: IptablesEIPs (集群级别)
		{Group: "kubeovn.io", Resource: "iptables-eips", Verb: "create"},
		{Group: "kubeovn.io", Resource: "iptables-eips", Verb: "get"},
		{Group: "kubeovn.io", Resource: "iptables-eips", Verb: "delete"},
		// KubeOVN: IptablesDnatRules (集群级别)
		{Group: "kubeovn.io", Resource: "iptables-dnat-rules", Verb: "create"},
		{Group: "kubeovn.io", Resource: "iptables-dnat-rules", Verb: "get"},
		{Group: "kubeovn.io", Resource: "iptables-dnat-rules", Verb: "delete"},
		// KubeOVN: IptablesSnatRules (集群级别)
		{Group: "kubeovn.io", Resource: "iptables-snat-rules", Verb: "create"},
		{Group: "kubeovn.io", Resource: "iptables-snat-rules", Verb: "get"},
		{Group: "kubeovn.io", Resource: "iptables-snat-rules", Verb: "delete"},
		// KubeOVN: VpcNatGateways (集群级别)
		{Group: "kubeovn.io", Resource: "vpc-nat-gateways", Verb: "create"},
		{Group: "kubeovn.io", Resource: "vpc-nat-gateways", Verb: "get"},
		{Group: "kubeovn.io", Resource: "vpc-nat-gateways", Verb: "delete"},
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
