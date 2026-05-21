package k8s

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"fmt"
	"strings"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	corev1 "k8s.io/api/core/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	VPCUnsupportedNodeLabelKey   = "node.cbctf.io/vpc-unsupported"
	VPCUnsupportedNodeLabelValue = "true"
	ExternalNetworkNodeLabelKey  = "node.cbctf.io/external-network"
)

type CreateVPCNatGatewayOptions struct {
	Name           string
	Labels         map[string]string
	VPC            string
	Subnet         string
	NetAttachDef   string
	LanIP          string
	ExternalSubnet string
	Interface      string
}

func CreateVPCNatGateway(ctx context.Context, options CreateVPCNatGatewayOptions) (*kubeovnv1.VpcNatGateway, model.RetVal) {
	var (
		gateway *kubeovnv1.VpcNatGateway
		err     error
	)
	gateway = &kubeovnv1.VpcNatGateway{
		ObjectMeta: metav1.ObjectMeta{
			Name:      options.Name,
			Namespace: "kube-system",
			Labels:    options.Labels,
		},
		Spec: kubeovnv1.VpcNatGatewaySpec{
			Vpc:             options.VPC,
			Subnet:          options.Subnet,
			LanIP:           options.LanIP,
			ExternalSubnets: []string{options.ExternalSubnet},
			Affinity: corev1.Affinity{
				NodeAffinity: &corev1.NodeAffinity{
					RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
						NodeSelectorTerms: []corev1.NodeSelectorTerm{
							{
								MatchExpressions: func() []corev1.NodeSelectorRequirement {
									matchExpressions := []corev1.NodeSelectorRequirement{
										{
											Key:      VPCUnsupportedNodeLabelKey,
											Operator: corev1.NodeSelectorOpNotIn,
											Values:   []string{VPCUnsupportedNodeLabelValue},
										},
									}
									if options.Interface != "" {
										matchExpressions = append(matchExpressions, corev1.NodeSelectorRequirement{
											Key:      ExternalNetworkNodeLabelKey,
											Operator: corev1.NodeSelectorOpIn,
											Values:   []string{options.Interface},
										})
									}
									return matchExpressions
								}(),
							},
						},
					},
				},
			},
			// kube-ovn 作为副 CNI 时, 通过注解注入 lanIp 到 eth0
			// https://github.com/kubeovn/kube-ovn/issues/6744
			Annotations: map[string]string{"v1.multus-cni.io/default-network": fmt.Sprintf("%s/%s", globalNamespace, options.NetAttachDef)},
		},
	}
	gateway, err = ovnClient.KubeovnV1().VpcNatGateways().Create(ctx, gateway, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create VPCNatGateway: %s", err)
		return nil, model.RetVal{Msg: i18n.K8S.CreateError, Attr: map[string]any{"Model": "VPCNatGateway", "Error": err.Error()}}
	}
	return gateway, model.SuccessRetVal()
}

func DeleteVPCNatGatewayCollection(ctx context.Context, labels ...map[string]string) model.RetVal {
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
	err := ovnClient.KubeovnV1().VpcNatGateways().DeleteCollection(ctx, metav1.DeleteOptions{}, options)
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete VPCNatGateway: %s", err)
		return model.RetVal{Msg: i18n.K8S.DeleteError, Attr: map[string]any{"Model": "VPCNatGateway", "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}
