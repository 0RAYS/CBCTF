package k8s

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"fmt"
	"strings"

	netv1 "k8s.io/api/networking/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CreateNetworkPolicyOptions struct {
	Name        string
	Labels      map[string]string
	MatchLabels map[string]string
	From        []*netv1.IPBlock
	To          []*netv1.IPBlock
}

func CreateNetworkPolicy(ctx context.Context, options CreateNetworkPolicyOptions) (*netv1.NetworkPolicy, model.RetVal) {
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
			Namespace: globalNamespace,
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
	networkPolicy, err = kubeClient.NetworkingV1().NetworkPolicies(globalNamespace).Create(ctx, networkPolicy, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create NetworkPolicy: %s", err)
		return nil, model.RetVal{Msg: i18n.K8S.CreateError, Attr: map[string]any{"Model": "NetworkPolicy", "Error": err.Error()}}
	}
	return networkPolicy, model.SuccessRetVal()
}

func GetNetworkPolicy(ctx context.Context, name string) (*netv1.NetworkPolicy, model.RetVal) {
	networkPolicy, err := kubeClient.NetworkingV1().NetworkPolicies(globalNamespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if apierror.IsNotFound(err) {
			return nil, model.RetVal{Msg: i18n.K8S.NotFound, Attr: map[string]any{"Model": "NetworkPolicy"}}
		}
		log.Logger.Warningf("Failed to get NetworkPolicy: %s", err)
		return nil, model.RetVal{Msg: i18n.K8S.GetError, Attr: map[string]any{"Model": "NetworkPolicy", "Error": err.Error()}}
	}
	return networkPolicy, model.SuccessRetVal()
}

func GetNetworkPolicyList(ctx context.Context, labels ...map[string]string) (*netv1.NetworkPolicyList, model.RetVal) {
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
	networkPolicyList, err := kubeClient.NetworkingV1().NetworkPolicies(globalNamespace).List(ctx, options)
	if err != nil {
		if apierror.IsNotFound(err) {
			return nil, model.RetVal{Msg: i18n.K8S.NotFound, Attr: map[string]any{"Model": "NetworkPolicy"}}
		}
		log.Logger.Warningf("Failed to list NetworkPolicy: %s", err)
		return nil, model.RetVal{Msg: i18n.K8S.GetError, Attr: map[string]any{"Model": "NetworkPolicy", "Error": err.Error()}}
	}
	return networkPolicyList, model.SuccessRetVal()
}

func DeleteNetworkPolicy(ctx context.Context, name string) model.RetVal {
	err := kubeClient.NetworkingV1().NetworkPolicies(globalNamespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete NetworkPolicy %s: %s", name, err)
		return model.RetVal{Msg: i18n.K8S.DeleteError, Attr: map[string]any{"Model": "NetworkPolicy", "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

func DeleteNetworkPolicyList(ctx context.Context, labels ...map[string]string) model.RetVal {
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
	err := kubeClient.NetworkingV1().NetworkPolicies(globalNamespace).DeleteCollection(ctx, metav1.DeleteOptions{}, options)
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete NetworkPolicy: %s", err)
		return model.RetVal{Msg: i18n.K8S.DeleteError, Attr: map[string]any{"Model": "NetworkPolicy", "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}
