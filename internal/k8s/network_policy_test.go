package k8s

import (
	"context"
	"slices"
	"testing"

	netv1 "k8s.io/api/networking/v1"

	"CBCTF/internal/model"
)

func TestCreateNetworkPolicyAllowsInternetAndPeers(t *testing.T) {
	useFakePods(t)
	allow := true
	policies := model.NetworkPolicies{{
		Service:       "web",
		AllowInternet: &allow,
		AllowExternal: &allow,
		Peers:         []string{"db"},
	}}
	peerIPs := map[string][]string{"db": {"10.1.0.3"}}
	created, ret := CreateNetworkPolicy(context.Background(), CreateNetworkPolicyOptions{
		Name:        "np-web",
		Labels:      map[string]string{"service": "web"},
		Selector:    map[string]string{"service": "web"},
		Service:     "web",
		Policies:    policies,
		SubnetCIDRs: []string{"10.1.0.0/24"},
		PeerIPs:     peerIPs,
	})
	if !ret.OK {
		t.Fatal(ret)
	}
	spec := created.Spec
	if len(spec.PolicyTypes) != 2 {
		t.Fatalf("unexpected policy types: %v", spec.PolicyTypes)
	}
	if len(spec.Egress) != 1 {
		t.Fatalf("unexpected egress rules: %+v", spec.Egress)
	}
	blocks := spec.Egress[0].To
	if len(blocks) != 3 {
		t.Fatalf("egress must allow the subnet, the peer and the internet: %+v", blocks)
	}
	if blocks[0].IPBlock.CIDR != "10.1.0.0/24" {
		t.Fatalf("egress must allow the local subnet: %+v", blocks[0])
	}
	if blocks[1].IPBlock.CIDR != "10.1.0.3/32" {
		t.Fatalf("egress must allow the peer address: %+v", blocks[1])
	}
	if blocks[2].IPBlock.CIDR != "0.0.0.0/0" || len(blocks[2].IPBlock.Except) == 0 {
		t.Fatalf("egress must allow public internet only: %+v", blocks[2])
	}
	if len(spec.Ingress) != 1 {
		t.Fatalf("unexpected ingress rules: %+v", spec.Ingress)
	}
	ingress := spec.Ingress[0].From
	if len(ingress) != 3 {
		t.Fatalf("ingress must allow the subnet, the peer and external traffic: %+v", ingress)
	}
	if ingress[0].IPBlock.CIDR != "10.1.0.0/24" || ingress[1].IPBlock.CIDR != "10.1.0.3/32" {
		t.Fatalf("unexpected ingress sources: %+v", ingress)
	}
	if ingress[2].IPBlock.CIDR != "0.0.0.0/0" {
		t.Fatalf("ingress must allow external traffic: %+v", ingress[2])
	}
}

func TestCreateNetworkPolicyDeniesEgressWithoutRules(t *testing.T) {
	useFakePods(t)
	deny := false
	policies := model.NetworkPolicies{{Service: "web", AllowInternet: &deny, AllowExternal: &deny}}
	created, ret := CreateNetworkPolicy(context.Background(), CreateNetworkPolicyOptions{
		Name:     "np-web",
		Labels:   map[string]string{"service": "web"},
		Selector: map[string]string{"service": "web"},
		Service:  "web",
		Policies: policies,
	})
	if !ret.OK {
		t.Fatal(ret)
	}
	spec := created.Spec
	if len(spec.Egress) != 0 {
		t.Fatalf("egress must be empty: %+v", spec.Egress)
	}
	if !slices.Contains(spec.PolicyTypes, netv1.PolicyTypeEgress) {
		t.Fatalf("egress policy type is required to deny all egress: %v", spec.PolicyTypes)
	}
	if len(spec.Ingress) != 0 {
		t.Fatalf("ingress must be denied: %+v", spec.Ingress)
	}
}

func TestCreateNetworkPolicyKeepsSubnetConnectivity(t *testing.T) {
	useFakePods(t)
	created, ret := CreateNetworkPolicy(context.Background(), CreateNetworkPolicyOptions{
		Name:        "np-web",
		Labels:      map[string]string{"service": "web"},
		Selector:    map[string]string{"service": "web"},
		Service:     "web",
		SubnetCIDRs: []string{"10.1.0.0/24"},
	})
	if !ret.OK {
		t.Fatal(ret)
	}
	spec := created.Spec
	if len(spec.Egress) != 1 || spec.Egress[0].To[0].IPBlock.CIDR != "10.1.0.0/24" {
		t.Fatalf("egress must keep the local subnet: %+v", spec.Egress)
	}
	if len(spec.Ingress) != 1 || spec.Ingress[0].From[0].IPBlock.CIDR != "10.1.0.0/24" {
		t.Fatalf("ingress must keep the local subnet: %+v", spec.Ingress)
	}
}

func TestCreateNetworkPolicyExcludesVPCSubnetsFromExternalTraffic(t *testing.T) {
	useFakePods(t)
	allow := true
	created, ret := CreateNetworkPolicy(context.Background(), CreateNetworkPolicyOptions{
		Name:        "np-web",
		Labels:      map[string]string{"service": "web"},
		Selector:    map[string]string{"service": "web"},
		Service:     "web",
		Policies:    model.NetworkPolicies{{Service: "web", AllowInternet: &allow, AllowExternal: &allow}},
		SubnetCIDRs: []string{"10.1.0.0/24"},
		VPCSubnets:  []string{"10.1.0.0/24", "10.2.0.0/24"},
	})
	if !ret.OK {
		t.Fatal(ret)
	}
	spec := created.Spec
	internet := spec.Egress[0].To[len(spec.Egress[0].To)-1].IPBlock
	if !slices.Contains(internet.Except, "10.2.0.0/24") {
		t.Fatalf("internet egress must exclude the challenge subnets: %+v", internet)
	}
	external := spec.Ingress[0].From[len(spec.Ingress[0].From)-1].IPBlock
	if external.CIDR != "0.0.0.0/0" || !slices.Contains(external.Except, "10.2.0.0/24") {
		t.Fatalf("external ingress must exclude the challenge subnets: %+v", external)
	}
}
