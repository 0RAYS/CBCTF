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
	Name     string
	Labels   map[string]string
	Policies model.NetworkPolicies
}

func CreateNetworkPolicy(ctx context.Context, options CreateNetworkPolicyOptions) (*netv1.NetworkPolicy, model.RetVal) {
	var (
		networkPolicy *netv1.NetworkPolicy
		err           error
	)
	ingress, egress := func(policies model.NetworkPolicies) ([]netv1.NetworkPolicyIngressRule, []netv1.NetworkPolicyEgressRule) {
		var ingress []netv1.NetworkPolicyIngressRule
		var egress []netv1.NetworkPolicyEgressRule
		var from []*netv1.IPBlock
		var to []*netv1.IPBlock
		for _, policy := range policies {
			from = append(from, policy.From...)
			to = append(to, policy.To...)
		}
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
	}(options.Policies)
	networkPolicy = &netv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      options.Name,
			Namespace: globalNamespace,
			Labels:    options.Labels,
		},
		Spec: netv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: options.Labels,
			},
			PolicyTypes: func() []netv1.PolicyType {
				policyTypes := []netv1.PolicyType{netv1.PolicyTypeEgress}
				if len(ingress) > 0 {
					policyTypes = append(policyTypes, netv1.PolicyTypeIngress)
				}
				return policyTypes
			}(),
			Ingress: ingress,
			Egress:  egress,
		},
	}
	networkPolicy, err = kubeClient.NetworkingV1().NetworkPolicies(globalNamespace).Create(ctx, networkPolicy, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create NetworkPolicy: %s", err)
		return nil, model.RetVal{Msg: i18n.K8S.CreateError, Attr: map[string]any{"Model": "NetworkPolicy", "Error": err.Error()}}
	}
	return networkPolicy, model.SuccessRetVal()
}

func DeleteNetworkPolicyCollection(ctx context.Context, labels ...map[string]string) model.RetVal {
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
