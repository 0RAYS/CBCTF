package model

import (
	"testing"

	netv1 "k8s.io/api/networking/v1"
)

func boolPointer(value bool) *bool {
	return &value
}

func TestResolveNetworkPolicyDefaultsToDeniedEgress(t *testing.T) {
	resolved := ResolveNetworkPolicy(nil, "web")
	if resolved.Internet {
		t.Fatal("default policy must deny internet egress")
	}
	if !resolved.AllowExternal || !resolved.RestrictIngress {
		t.Fatal("default policy must keep external ingress open")
	}
	if len(resolved.Peers) != 0 {
		t.Fatalf("default policy must not allow peers: %v", resolved.Peers)
	}
}

func TestResolveNetworkPolicyKeepsLegacyCIDRSemantics(t *testing.T) {
	legacy := NetworkPolicies{{
		Service: "web",
		From:    []*netv1.IPBlock{{CIDR: "10.1.0.0/24"}},
		To:      []*netv1.IPBlock{{CIDR: "0.0.0.0/0"}},
	}}
	resolved := ResolveNetworkPolicy(legacy, "web")
	if resolved.Internet {
		t.Fatal("legacy entry without allow_internet must not enable internet egress")
	}
	if !resolved.RestrictIngress {
		t.Fatal("legacy entry with from blocks must restrict ingress")
	}
	if resolved.AllowExternal {
		t.Fatal("legacy entry with from blocks must not allow external ingress")
	}
	if len(resolved.To) != 1 || len(resolved.From) != 1 {
		t.Fatalf("legacy blocks must be preserved: %+v", resolved)
	}
	open := ResolveNetworkPolicy(NetworkPolicies{{Service: "web"}}, "web")
	if !open.AllowExternal || open.RestrictIngress {
		t.Fatal("legacy entry without from blocks must keep ingress unrestricted")
	}
}

func TestResolveNetworkPolicyScopesEntriesPerService(t *testing.T) {
	policies := NetworkPolicies{
		{Service: "web", AllowInternet: boolPointer(true), Peers: []string{"db"}},
		{Service: "db", AllowInternet: boolPointer(false), AllowExternal: boolPointer(false)},
	}
	web := ResolveNetworkPolicy(policies, "web")
	if !web.Internet || web.AllowExternal {
		t.Fatalf("unexpected web policy: %+v", web)
	}
	if len(web.Peers) != 1 || web.Peers[0] != "db" {
		t.Fatalf("unexpected web peers: %v", web.Peers)
	}
	db := ResolveNetworkPolicy(policies, "db")
	if db.Internet {
		t.Fatal("db entry must not inherit internet egress from web")
	}
	if !db.RestrictIngress || db.AllowExternal {
		t.Fatalf("unexpected db ingress policy: %+v", db)
	}
	if len(db.Peers) != 1 || db.Peers[0] != "web" {
		t.Fatalf("peer relation must be symmetric: %v", db.Peers)
	}
}

func TestResolveNetworkPolicyAggregatesSharedPodEntries(t *testing.T) {
	policies := NetworkPolicies{
		{Service: "web", AllowInternet: boolPointer(true)},
		{Service: "db", EgressCIDRs: []string{"10.1.0.0/24"}},
	}
	shared := ResolveNetworkPolicy(policies, "")
	if !shared.Internet {
		t.Fatal("shared pod must inherit internet egress from any entry")
	}
	if len(shared.CIDRs) != 1 || shared.CIDRs[0] != "10.1.0.0/24" {
		t.Fatalf("unexpected shared egress ranges: %v", shared.CIDRs)
	}
}
