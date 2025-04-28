package k8s

import (
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"context"
	netv1 "k8s.io/api/networking/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetNetworkPolicy(ctx context.Context, name string) (*netv1.NetworkPolicy, bool, string) {
	networkPolicy, err := client.NetworkingV1().NetworkPolicies(NamespaceName).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if apierror.IsNotFound(err) {
			return nil, false, "NetworkPolicyNotFound"
		}
		log.Logger.Warningf("Failed to get NetworkPolicy %s: %v", name, err)
		return nil, false, "GetNetworkPolicyError"
	}
	return networkPolicy, true, "Success"
}

func CreateNetworkPolicy(ctx context.Context, pod model.Pod, policy model.NetworkPolicy) (*netv1.NetworkPolicy, bool, string) {
	if _, ok, _ := GetNetworkPolicy(ctx, pod.NetworkPolicyName); ok {
		DeleteNetworkPolicy(ctx, pod.NetworkPolicyName)
	}
	if len(policy.From) < 1 && len(policy.To) < 1 {
		return nil, true, "Success"
	}
	networkPolicy := &netv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pod.NetworkPolicyName,
			Namespace: NamespaceName,
		},
		Spec: netv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"victim": pod.Name,
				},
			},
		},
	}
	policyTypes, ingress, egress := func() ([]netv1.PolicyType, []netv1.NetworkPolicyIngressRule, []netv1.NetworkPolicyEgressRule) {
		var t []netv1.PolicyType
		var ingress []netv1.NetworkPolicyIngressRule
		var egress []netv1.NetworkPolicyEgressRule
		if len(policy.From) > 0 {
			var peers []netv1.NetworkPolicyPeer
			for _, from := range policy.From {
				peers = append(peers, netv1.NetworkPolicyPeer{
					IPBlock: &netv1.IPBlock{
						CIDR:   from.CIDR,
						Except: from.Except,
					},
				})
			}
			ingress = append(ingress, netv1.NetworkPolicyIngressRule{From: peers})
			t = append(t, netv1.PolicyTypeIngress)
		}
		if len(policy.To) > 0 {
			var peers []netv1.NetworkPolicyPeer
			for _, to := range policy.To {
				peers = append(peers, netv1.NetworkPolicyPeer{
					IPBlock: &netv1.IPBlock{
						CIDR:   to.CIDR,
						Except: to.Except,
					},
				})
			}
			egress = append(egress, netv1.NetworkPolicyEgressRule{To: peers})
			t = append(t, netv1.PolicyTypeEgress)
		}
		return t, ingress, egress
	}()
	if len(policyTypes) > 0 {
		networkPolicy.Spec.PolicyTypes = policyTypes
	}
	if len(ingress) > 0 {
		networkPolicy.Spec.Ingress = ingress
	}
	if len(egress) > 0 {
		networkPolicy.Spec.Egress = egress
	}
	var err error
	networkPolicy, err = client.NetworkingV1().NetworkPolicies(NamespaceName).Create(ctx, networkPolicy, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create NetworkPolicy %s: %v", pod.NetworkPolicyName, err)
		return nil, false, "CreateNetworkPolicyError"
	}
	return networkPolicy, true, "Success"
}

func DeleteNetworkPolicy(ctx context.Context, name string) (bool, string) {
	err := client.NetworkingV1().NetworkPolicies(NamespaceName).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete NetworkPolicy %s: %v", name, err)
		return false, "DeleteNetworkPolicyError"
	}
	return true, "Success"
}
