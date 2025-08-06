package k8s

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"context"
	"fmt"
	"strings"

	kubeovnv1 "github.com/JBNRZ/kubeovn-api/pkg/apis/kubeovn/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CreateVPCNatGatewayOptions struct {
	Name           string
	Labels         map[string]string
	VPC            string
	Subnet         string
	LanIP          string
	ExternalSubnet []string
}

func CreateVPCNatGateway(ctx context.Context, options CreateVPCNatGatewayOptions) (*kubeovnv1.VpcNatGateway, bool, string) {
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
			ExternalSubnets: options.ExternalSubnet,
		},
	}
	gateway, err = kubeOVNClient.KubeovnV1().VpcNatGateways().Create(ctx, gateway, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create VPCNatGateway: %v", err)
		return nil, false, i18n.CreateVPCNatGatewayError
	}
	return gateway, true, i18n.Success
}

func GetVPCNatGateway(ctx context.Context, name string) (*kubeovnv1.VpcNatGateway, bool, string) {
	gateway, err := kubeOVNClient.KubeovnV1().VpcNatGateways().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if apierror.IsNotFound(err) {
			return nil, false, i18n.VPCNatGatewayNotFound
		}
		log.Logger.Warningf("Failed to get VPCNatGateway: %v", err)
		return nil, false, i18n.GetVPCNatGatewayError
	}
	return gateway, true, i18n.Success
}

func GetVPCNatGatewayList(ctx context.Context, labels ...map[string]string) (*kubeovnv1.VpcNatGatewayList, bool, string) {
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
	gatewayList, err := kubeOVNClient.KubeovnV1().VpcNatGateways().List(ctx, options)
	if err != nil {
		log.Logger.Warningf("Failed to list VPCNatGateway: %v", err)
		return nil, false, i18n.GetVPCNatGatewayError
	}
	return gatewayList, true, i18n.Success
}

func DeleteVPCNatGateway(ctx context.Context, name string) (bool, string) {
	err := kubeOVNClient.KubeovnV1().VpcNatGateways().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete VPCNatGateway: %v", err)
		return false, i18n.DeleteVPCNatGatewayError
	}
	return true, i18n.Success
}

func DeleteVPCNatGatewayList(ctx context.Context, labels ...map[string]string) (bool, string) {
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
	err := kubeOVNClient.KubeovnV1().VpcNatGateways().DeleteCollection(ctx, metav1.DeleteOptions{}, options)
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete VPCNatGateway: %v", err)
		return false, i18n.DeleteVPCNatGatewayError
	}
	return true, i18n.Success
}
