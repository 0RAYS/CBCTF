package model

import (
	"slices"

	netv1 "k8s.io/api/networking/v1"
)

// ResolvedNetworkPolicy 汇总作用于单个容器的网络策略规则。
type ResolvedNetworkPolicy struct {
	Internet bool
	CIDRs    []string
	Peers    []string
	From     []*netv1.IPBlock
	To       []*netv1.IPBlock
}

// ResolveNetworkPolicy 汇总作用于 service 的策略条目。
// service 为空表示汇总全部条目，用于共享 Pod 网络下唯一的一个 Pod。
// 没有匹配条目时出站全部拒绝，入站仍接受集群外部的访问。
func ResolveNetworkPolicy(policies NetworkPolicies, service string) ResolvedNetworkPolicy {
	resolved := ResolvedNetworkPolicy{}
	for _, policy := range policies {
		if service != "" && policy.Service != "" && policy.Service != service {
			continue
		}
		if policy.AllowInternet != nil && *policy.AllowInternet {
			resolved.Internet = true
		}
		resolved.CIDRs = append(resolved.CIDRs, policy.EgressCIDRs...)
		resolved.From = append(resolved.From, policy.From...)
		resolved.To = append(resolved.To, policy.To...)
	}
	resolved.Peers = resolvePeers(policies, service)
	return resolved
}

// resolvePeers 计算容器之间的互通关系，取对称闭包：任一侧声明的互通对两端同时生效。
func resolvePeers(policies NetworkPolicies, service string) []string {
	peers := make(map[string]struct{})
	for _, policy := range policies {
		if service != "" && policy.Service != "" && policy.Service != service {
			continue
		}
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
