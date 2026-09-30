package k8s

import (
	"context"

	netv1 "k8s.io/api/networking/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
)

type CreateNetworkPolicyOptions struct {
	Name     string
	Labels   map[string]string
	Policies model.NetworkPolicies
}

func CreateNetworkPolicy(ctx context.Context, options CreateNetworkPolicyOptions) (*netv1.NetworkPolicy, model.RetVal) {
	var from, to []netv1.NetworkPolicyPeer
	for _, policy := range options.Policies {
		for _, block := range policy.From {
			from = append(from, netv1.NetworkPolicyPeer{IPBlock: block})
		}
		for _, block := range policy.To {
			to = append(to, netv1.NetworkPolicyPeer{IPBlock: block})
		}
	}
	spec := netv1.NetworkPolicySpec{
		PodSelector: metav1.LabelSelector{MatchLabels: options.Labels},
		PolicyTypes: []netv1.PolicyType{netv1.PolicyTypeEgress},
	}
	if len(from) > 0 {
		spec.Ingress = []netv1.NetworkPolicyIngressRule{{From: from}}
		spec.PolicyTypes = append(spec.PolicyTypes, netv1.PolicyTypeIngress)
	}
	if len(to) > 0 {
		spec.Egress = []netv1.NetworkPolicyEgressRule{{To: to}}
	}
	networkPolicy := &netv1.NetworkPolicy{
		OwnerReferences: resourceOwners(ctx),
		Name:            options.Name,
		Namespace:       globalNamespace,
		Labels:          options.Labels,
		Spec:            spec,
	}
	networkPolicy, err := kubeClient.NetworkingV1().NetworkPolicies(globalNamespace).Create(ctx, networkPolicy, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create NetworkPolicy: %s", err)
		return nil, model.RetVal{Msg: i18n.K8S.CreateError, Attr: map[string]any{"Model": "NetworkPolicy", "Error": err.Error()}}
	}
	return networkPolicy, model.SuccessRetVal()
}

func DeleteNetworkPolicyCollection(ctx context.Context, labels map[string]string) model.RetVal {
	options, ret := deleteCollectionOptions("NetworkPolicy", labels)
	if !ret.OK {
		return ret
	}
	err := kubeClient.NetworkingV1().NetworkPolicies(globalNamespace).DeleteCollection(ctx, metav1.DeleteOptions{}, options)
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete NetworkPolicy: %s", err)
		return model.RetVal{Msg: i18n.K8S.DeleteError, Attr: map[string]any{"Model": "NetworkPolicy", "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}
