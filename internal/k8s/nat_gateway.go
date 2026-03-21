package k8s

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"fmt"
	"strings"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
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
			ExternalSubnets: options.ExternalSubnet,
		},
	}
	gateway, err = kubeOVNClient.KubeovnV1().VpcNatGateways().Create(ctx, gateway, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create VPCNatGateway: %s", err)
		return nil, model.RetVal{Msg: i18n.K8S.CreateError, Attr: map[string]any{"Model": "VPCNatGateway", "Error": err.Error()}}
	}
	return gateway, model.SuccessRetVal()
}

func GetVPCNatGateway(ctx context.Context, name string) (*kubeovnv1.VpcNatGateway, model.RetVal) {
	gateway, err := kubeOVNClient.KubeovnV1().VpcNatGateways().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if apierror.IsNotFound(err) {
			return nil, model.RetVal{Msg: i18n.K8S.NotFound, Attr: map[string]any{"Model": "VPCNatGateway"}}
		}
		log.Logger.Warningf("Failed to get VPCNatGateway: %s", err)
		return nil, model.RetVal{Msg: i18n.K8S.GetError, Attr: map[string]any{"Model": "VPCNatGateway", "Error": err.Error()}}
	}
	return gateway, model.SuccessRetVal()
}

func GetVPCNatGatewayList(ctx context.Context, labels ...map[string]string) (*kubeovnv1.VpcNatGatewayList, model.RetVal) {
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
		log.Logger.Warningf("Failed to list VPCNatGateway: %s", err)
		return nil, model.RetVal{Msg: i18n.K8S.GetError, Attr: map[string]any{"Model": "VPCNatGateway", "Error": err.Error()}}
	}
	return gatewayList, model.SuccessRetVal()
}

func DeleteVPCNatGateway(ctx context.Context, name string) model.RetVal {
	err := kubeOVNClient.KubeovnV1().VpcNatGateways().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete VPCNatGateway: %s", err)
		return model.RetVal{Msg: i18n.K8S.DeleteError, Attr: map[string]any{"Model": "VPCNatGateway", "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

func DeleteVPCNatGatewayList(ctx context.Context, labels ...map[string]string) model.RetVal {
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
		log.Logger.Warningf("Failed to delete VPCNatGateway: %s", err)
		return model.RetVal{Msg: i18n.K8S.DeleteError, Attr: map[string]any{"Model": "VPCNatGateway", "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}
