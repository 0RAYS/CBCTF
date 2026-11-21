package model

import (
	"slices"

	netv1 "k8s.io/api/networking/v1"
)

// ResolvedNetworkPolicy 汇总作用于单个容器的网络策略规则。
type ResolvedNetworkPolicy struct {
	Internet        bool
	CIDRs           []string
	Peers           []string
	AllowExternal   bool
	RestrictIngress bool
	From            []*netv1.IPBlock
	To              []*netv1.IPBlock
}

// ResolveNetworkPolicy 汇总作用于 service 的策略条目。
// service 为空表示汇总全部条目，用于共享 Pod 网络下唯一的一个 Pod。
// 没有匹配条目时返回平台默认策略：拒绝全部出站，允许集群外部入站。
func ResolveNetworkPolicy(policies NetworkPolicies, service string) ResolvedNetworkPolicy {
	matched := make([]NetworkPolicy, 0, len(policies))
	for _, policy := range policies {
		if service == "" || policy.Service == "" || policy.Service == service {
			matched = append(matched, policy)
		}
	}
	if len(matched) == 0 {
		return ResolvedNetworkPolicy{AllowExternal: true, RestrictIngress: true}
	}
	resolved := ResolvedNetworkPolicy{}
	legacyOpen := false
	for _, policy := range matched {
		if policy.AllowInternet != nil && *policy.AllowInternet {
			resolved.Internet = true
		}
		if policy.AllowExternal != nil {
			resolved.AllowExternal = resolved.AllowExternal || *policy.AllowExternal
			resolved.RestrictIngress = true
		}
		if len(policy.From) > 0 || len(policy.Peers) > 0 {
			resolved.RestrictIngress = true
		} else if policy.AllowExternal == nil {
			// 旧条目没有声明入站限制，保持历史上的不限入站行为。
			legacyOpen = true
		}
		resolved.CIDRs = append(resolved.CIDRs, policy.EgressCIDRs...)
		resolved.From = append(resolved.From, policy.From...)
		resolved.To = append(resolved.To, policy.To...)
	}
	resolved.Peers = resolvePeers(policies, service, matched)
	if legacyOpen {
		resolved.AllowExternal = true
		resolved.RestrictIngress = false
	}
	return resolved
}

// resolvePeers 计算容器之间的互通关系，取对称闭包：任一侧声明的互通对两端同时生效。
func resolvePeers(policies NetworkPolicies, service string, matched []NetworkPolicy) []string {
	peers := make(map[string]struct{})
	for _, policy := range matched {
		for _, peer := range policy.Peers {
			if peer != "" && peer != service {
				peers[peer] = struct{}{}
			}
		}
	}
	if service != "" {
		for _, policy := range policies {
			if policy.Service == "" || policy.Service == service {
				continue
			}
			if slices.Contains(policy.Peers, service) {
				peers[policy.Service] = struct{}{}
			}
		}
	}
	result := make([]string, 0, len(peers))
	for peer := range peers {
		result = append(result, peer)
	}
	slices.Sort(result)
	return result
}
