package k8s

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"context"
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
			Namespace: GlobalNamespace,
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

func GetVPCNatGatewayList(ctx context.Context) (*kubeovnv1.VpcNatGatewayList, bool, string) {
	gatewayList, err := kubeOVNClient.KubeovnV1().VpcNatGateways().List(ctx, metav1.ListOptions{})
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
