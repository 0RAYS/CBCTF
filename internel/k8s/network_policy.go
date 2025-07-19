package k8s

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"context"
	"fmt"
	netv1 "k8s.io/api/networking/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type CreateNetworkPolicyOptions struct {
	Name        string
	Labels      map[string]string
	MatchLabels map[string]string
	From        []*netv1.IPBlock
	To          []*netv1.IPBlock
}

func CreateNetworkPolicy(ctx context.Context, options CreateNetworkPolicyOptions) (*netv1.NetworkPolicy, bool, string) {
	var (
		networkPolicy *netv1.NetworkPolicy
		err           error
	)
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
			Name:      options.Name,
			Namespace: GlobalNamespace,
			Labels:    options.Labels,
		},
		Spec: netv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: options.MatchLabels,
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
	networkPolicy, err = kubeClient.NetworkingV1().NetworkPolicies(GlobalNamespace).Create(ctx, networkPolicy, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create NetworkPolicy: %v", err)
		return nil, false, i18n.CreateNetworkPolicyError
	}
	return networkPolicy, true, i18n.Success
}

func GetNetworkPolicyList(ctx context.Context, labels ...map[string]string) (*netv1.NetworkPolicyList, bool, string) {
	var options metav1.ListOptions
	if len(labels) > 0 {
		var selector string
		for k, v := range labels[0] {
			selector += fmt.Sprintf("%s=%s,", k, v)
		}
		options = metav1.ListOptions{
			LabelSelector: strings.TrimSuffix(selector, ","),
		}
	}
	networkPolicyList, err := kubeClient.NetworkingV1().NetworkPolicies(GlobalNamespace).List(ctx, options)
	if err != nil {
		if apierror.IsNotFound(err) {
			return nil, false, i18n.NetworkPolicyNotFound
		}
		log.Logger.Warningf("Failed to list NetworkPolicy: %v", err)
		return nil, false, i18n.GetNetworkPolicyError
	}
	return networkPolicyList, true, i18n.Success
}

func GetNetworkPolicyListByPodName(ctx context.Context, key, podName string) (*netv1.NetworkPolicyList, bool, string) {
	networkPolicyList, err := kubeClient.NetworkingV1().NetworkPolicies(GlobalNamespace).List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", key, podName),
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
	err := kubeClient.NetworkingV1().NetworkPolicies(GlobalNamespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete NetworkPolicy %s: %v", name, err)
		return false, i18n.DeleteNetworkPolicyError
	}
	return true, i18n.Success
}

// DeleteNetworkPolicyListByPodName TODO: 有可能删不干净
func DeleteNetworkPolicyListByPodName(ctx context.Context, key, podName string) (bool, string) {
	networkPolicyList, ok, msg := GetNetworkPolicyListByPodName(ctx, key, podName)
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

func DeleteNetworkPolicyByLabels(ctx context.Context, labels map[string]string) (bool, string) {
	networkPolicyList, ok, msg := GetNetworkPolicyList(ctx, labels)
	if !ok {
		return false, msg
	}
	for _, np := range networkPolicyList.Items {
		if ok, msg = DeleteNetworkPolicy(ctx, np.Name); !ok {
			return false, msg
		}
	}
	return true, i18n.Success
}
