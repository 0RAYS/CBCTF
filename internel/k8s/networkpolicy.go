package k8s

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"CBCTF/internel/utils"
	"context"
	"fmt"
	netv1 "k8s.io/api/networking/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type CreateNetworkPolicyOptions struct {
	PodName string
	From    []*netv1.IPBlock
	To      []*netv1.IPBlock
}

func CreateNetworkPolicy(ctx context.Context, options CreateNetworkPolicyOptions) (*netv1.NetworkPolicy, bool, string) {
	var (
		networkPolicy *netv1.NetworkPolicy
		err           error
	)
	// 可能会连续创建多个, 不删除之前的, 交由定时任务处理
	//DeleteNetworkPolicyListByPodName(ctx, options.PodName)
	ingress, egress := func(from []*netv1.IPBlock, to []*netv1.IPBlock) ([]netv1.NetworkPolicyIngressRule, []netv1.NetworkPolicyEgressRule) {
		var ingress []netv1.NetworkPolicyIngressRule
		var egress []netv1.NetworkPolicyEgressRule
		if len(from) > 0 {
			var peers []netv1.NetworkPolicyPeer
			for _, f := range from {
				peers = append(peers, netv1.NetworkPolicyPeer{IPBlock: f})
			}
			ingress = append(ingress, netv1.NetworkPolicyIngressRule{From: peers})
		}
		if len(to) > 0 {
			var peers []netv1.NetworkPolicyPeer
			for _, t := range to {
				peers = append(peers, netv1.NetworkPolicyPeer{IPBlock: t})
			}
			egress = append(egress, netv1.NetworkPolicyEgressRule{To: peers})
		}
		return ingress, egress
	}(options.From, options.To)
	networkPolicy = &netv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("np-%s", strings.ToLower(utils.RandStr(10))),
			Namespace: namespaceName,
			Labels: map[string]string{
				"victim": options.PodName,
			},
		},
		Spec: netv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"victim": options.PodName,
				},
			},
			// 默认不允许出网
			PolicyTypes: []netv1.PolicyType{netv1.PolicyTypeEgress},
		},
	}
	if len(ingress) > 0 {
		networkPolicy.Spec.PolicyTypes = append(networkPolicy.Spec.PolicyTypes, netv1.PolicyTypeIngress)
		networkPolicy.Spec.Ingress = ingress
	}
	if len(egress) > 0 {
		networkPolicy.Spec.Egress = egress
	}
	networkPolicy, err = kubeClient.NetworkingV1().NetworkPolicies(namespaceName).Create(ctx, networkPolicy, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create pod %s NetworkPolicy: %v", options.PodName, err)
		return nil, false, i18n.CreateNetworkPolicyError
	}
	return networkPolicy, true, i18n.Success
}

func GetNetworkPolicyList(ctx context.Context) (*netv1.NetworkPolicyList, bool, string) {
	networkPolicyList, err := kubeClient.NetworkingV1().NetworkPolicies(namespaceName).List(ctx, metav1.ListOptions{})
	if err != nil {
		if apierror.IsNotFound(err) {
			return nil, false, i18n.NetworkPolicyNotFound
		}
		log.Logger.Warningf("Failed to list NetworkPolicy: %v", err)
		return nil, false, i18n.GetNetworkPolicyError
	}
	return networkPolicyList, true, i18n.Success
}

func GetNetworkPolicyListByPodName(ctx context.Context, podName string) (*netv1.NetworkPolicyList, bool, string) {
	networkPolicyList, err := kubeClient.NetworkingV1().NetworkPolicies(namespaceName).List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("victim=%s", podName),
	})
	if err != nil {
		if apierror.IsNotFound(err) {
			return nil, false, i18n.NetworkPolicyNotFound
		}
		log.Logger.Warningf("Failed to list Pod %s NetworkPolicy: %v", podName, err)
		return nil, false, i18n.GetNetworkPolicyError
	}
	return networkPolicyList, true, i18n.Success
}

func DeleteNetworkPolicy(ctx context.Context, name string) (bool, string) {
	err := kubeClient.NetworkingV1().NetworkPolicies(namespaceName).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete NetworkPolicy %s: %v", name, err)
		return false, i18n.DeleteNetworkPolicyError
	}
	return true, i18n.Success
}

// DeleteNetworkPolicyListByPodName TODO: 有可能删不干净
func DeleteNetworkPolicyListByPodName(ctx context.Context, podName string) (bool, string) {
	networkPolicyList, ok, msg := GetNetworkPolicyListByPodName(ctx, podName)
	if !ok {
		if msg != i18n.NetworkPolicyNotFound {
			return false, i18n.GetNetworkPolicyError
		}
		return true, i18n.Success
	}
	for _, np := range networkPolicyList.Items {
		if ok, msg = DeleteNetworkPolicy(ctx, np.Name); !ok {
			return false, msg
		}
	}
	return true, i18n.Success
}
