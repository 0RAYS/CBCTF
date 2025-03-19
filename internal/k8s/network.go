package k8s

import (
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"context"
	netv1 "k8s.io/api/networking/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateNetworkPolicy(ctx context.Context, docker model.Docker, usage model.Usage) (*netv1.NetworkPolicy, bool, string) {
	if len(usage.NetworkPolicy.From) < 1 && len(usage.NetworkPolicy.To) < 1 {
		return nil, true, "Success"
	}
	networkPolicy := &netv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      docker.NetworkPolicyName,
			Namespace: NamespaceName,
		},
		Spec: netv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": docker.PodName,
				},
			},
		},
	}
	policyTypes, ingress, egress := func() ([]netv1.PolicyType, []netv1.NetworkPolicyIngressRule, []netv1.NetworkPolicyEgressRule) {
		var t []netv1.PolicyType
		var ingress []netv1.NetworkPolicyIngressRule
		var egress []netv1.NetworkPolicyEgressRule
		if len(usage.NetworkPolicy.From) > 0 {
			var peers []netv1.NetworkPolicyPeer
			for _, from := range usage.NetworkPolicy.From {
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
		if len(usage.NetworkPolicy.To) > 0 {
			var peers []netv1.NetworkPolicyPeer
			for _, to := range usage.NetworkPolicy.To {
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
	networkPolicy, err = Client.NetworkingV1().NetworkPolicies(NamespaceName).Create(ctx, networkPolicy, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create NetworkPolicy %s: %v", docker.NetworkPolicyName, err)
		return nil, false, "CreateNetworkPolicyError"
	}
	return networkPolicy, true, "Success"
}

func DeleteNetworkPolicy(ctx context.Context, name string) (bool, string) {
	err := Client.NetworkingV1().NetworkPolicies(NamespaceName).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete NetworkPolicy %s: %v", name, err)
		return false, "DeleteNetworkPolicyError"
	}
	return true, "Success"
}

func checkBlock(block utils.IPBlock) bool {

}
