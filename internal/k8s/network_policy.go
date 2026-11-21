package k8s

import (
	"context"
	"net/netip"

	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
)

// publicInternetBlock 出公网范围：0.0.0.0/0 扣除私有网段。
// 私有网段不在默认范围内，需要题目显式放行。
var publicInternetBlock = netv1.IPBlock{
	CIDR:   "0.0.0.0/0",
	Except: []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16", "100.64.0.0/10", "169.254.0.0/16"},
}

type CreateNetworkPolicyOptions struct {
	Name        string
	Labels      map[string]string
	Selector    map[string]string
	Service     string
	Policies    model.NetworkPolicies
	SubnetCIDRs []string
	PeerIPs     map[string][]string
	VPCSubnets  []string
}

func CreateNetworkPolicy(ctx context.Context, options CreateNetworkPolicyOptions) (*netv1.NetworkPolicy, model.RetVal) {
	resolved := model.ResolveNetworkPolicy(options.Policies, options.Service)
	selector := options.Selector
	if selector == nil {
		selector = options.Labels
	}
	networkPolicy := &netv1.NetworkPolicy{
		OwnerReferences: resourceOwners(ctx),
		Name:            options.Name,
		Namespace:       globalNamespace,
		Labels:          options.Labels,
		Spec:            buildNetworkPolicySpec(resolved, options.SubnetCIDRs, options.PeerIPs, options.VPCSubnets, selector),
	}
	networkPolicy, err := kubeClient.NetworkingV1().NetworkPolicies(globalNamespace).Create(ctx, networkPolicy, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create NetworkPolicy %s: %s", options.Name, err)
		return nil, model.FailRetVal(i18n.K8S.CreateError, "NetworkPolicy", options.Name, err)
	}
	return networkPolicy, model.SuccessRetVal()
}

// buildNetworkPolicySpec 把解析后的策略转换为 Kubernetes NetworkPolicy。
// 出站始终受策略类型约束：没有规则时表示拒绝全部出站。
// 容器所在网段与互通容器的地址始终放行，其余目标由策略决定。
// 入站默认接受集群外部的访问，填写自定义入站 CIDR 后只有这些来源可以访问。
func buildNetworkPolicySpec(
	resolved model.ResolvedNetworkPolicy,
	subnetCIDRs []string,
	peerIPs map[string][]string,
	vpcSubnets []string,
	selector map[string]string,
) netv1.NetworkPolicySpec {
	spec := netv1.NetworkPolicySpec{
		PodSelector: metav1.LabelSelector{MatchLabels: selector},
		PolicyTypes: []netv1.PolicyType{netv1.PolicyTypeEgress, netv1.PolicyTypeIngress},
	}
	subnetBlocks := make([]netv1.NetworkPolicyPeer, 0, len(subnetCIDRs))
	for _, cidr := range subnetCIDRs {
		block, err := parseIPBlock(cidr)
		if err != nil {
			log.Logger.Warningf("Skipped invalid subnet CIDR %s: %s", cidr, err)
			continue
		}
		subnetBlocks = append(subnetBlocks, netv1.NetworkPolicyPeer{IPBlock: block})
	}
	peerBlocks := make([]netv1.NetworkPolicyPeer, 0, len(resolved.Peers))
	for _, peer := range resolved.Peers {
		for _, ip := range peerIPs[peer] {
			peerBlocks = append(peerBlocks, netv1.NetworkPolicyPeer{IPBlock: &netv1.IPBlock{CIDR: ip + "/32"}})
		}
	}
	egressBlocks := make([]netv1.NetworkPolicyPeer, 0)
	for _, cidr := range resolved.CIDRs {
		block, err := parseIPBlock(cidr)
		if err != nil {
			log.Logger.Warningf("Skipped invalid egress CIDR %s: %s", cidr, err)
			continue
		}
		egressBlocks = append(egressBlocks, netv1.NetworkPolicyPeer{IPBlock: block})
	}
	for _, block := range resolved.To {
		egressBlocks = append(egressBlocks, netv1.NetworkPolicyPeer{IPBlock: block})
	}
	if resolved.Internet {
		egressBlocks = append(egressBlocks, netv1.NetworkPolicyPeer{IPBlock: internetBlock(vpcSubnets)})
	}
	spec.Egress = mergeEgressRules(concatPeers(subnetBlocks, peerBlocks, egressBlocks))

	ingressBlocks := make([]netv1.NetworkPolicyPeer, 0, len(resolved.From)+1)
	for _, block := range resolved.From {
		ingressBlocks = append(ingressBlocks, netv1.NetworkPolicyPeer{IPBlock: block})
	}
	if len(ingressBlocks) == 0 {
		ingressBlocks = append(ingressBlocks, netv1.NetworkPolicyPeer{IPBlock: externalBlock(vpcSubnets)})
	}
	spec.Ingress = mergeIngressRules(concatPeers(subnetBlocks, peerBlocks, ingressBlocks))
	return spec
}

// internetBlock 出公网范围：0.0.0.0/0 扣除私有网段与题目 VPC 网段。
func internetBlock(vpcSubnets []string) *netv1.IPBlock {
	block := publicInternetBlock.DeepCopy()
	block.Except = append(block.Except, vpcSubnets...)
	return block
}

// externalBlock 集群外部访问范围：0.0.0.0/0 扣除题目 VPC 网段，
// VPC 内的容器按互通规则放行。
func externalBlock(vpcSubnets []string) *netv1.IPBlock {
	return &netv1.IPBlock{CIDR: "0.0.0.0/0", Except: append([]string(nil), vpcSubnets...)}
}

func concatPeers(groups ...[]netv1.NetworkPolicyPeer) []netv1.NetworkPolicyPeer {
	result := make([]netv1.NetworkPolicyPeer, 0)
	for _, group := range groups {
		result = append(result, group...)
	}
	return result
}

// mergeEgressRules 把放行目标合并成一条出站规则。
func mergeEgressRules(peers []netv1.NetworkPolicyPeer) []netv1.NetworkPolicyEgressRule {
	if len(peers) == 0 {
		return nil
	}
	return []netv1.NetworkPolicyEgressRule{{To: peers}}
}

// mergeIngressRules 把放行来源合并成一条入站规则。
func mergeIngressRules(peers []netv1.NetworkPolicyPeer) []netv1.NetworkPolicyIngressRule {
	if len(peers) == 0 {
		return nil
	}
	return []netv1.NetworkPolicyIngressRule{{From: peers}}
}

func parseIPBlock(cidr string) (*netv1.IPBlock, error) {
	if _, err := netip.ParsePrefix(cidr); err != nil {
		return nil, err
	}
	return &netv1.IPBlock{CIDR: cidr}, nil
}

func DeleteNetworkPolicyCollection(ctx context.Context, labels map[string]string) model.RetVal {
	return deleteCollection("NetworkPolicy", "NetworkPolicies", labels, func(options metav1.ListOptions) error {
		return kubeClient.NetworkingV1().NetworkPolicies(globalNamespace).DeleteCollection(ctx, metav1.DeleteOptions{}, options)
	})
}
