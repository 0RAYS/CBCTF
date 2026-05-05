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

type CreateSubnetOptions struct {
	Name       string
	Labels     map[string]string
	VPC        string
	CIDR       string
	Gateway    string
	ExcludeIPs []string
	Provider   string
}

func CreateSubnet(ctx context.Context, options CreateSubnetOptions) (*kubeovnv1.Subnet, model.RetVal) {
	var (
		subnet *kubeovnv1.Subnet
		err    error
	)
	subnet = &kubeovnv1.Subnet{
		ObjectMeta: metav1.ObjectMeta{
			Name:   options.Name,
			Labels: options.Labels,
		},
		Spec: kubeovnv1.SubnetSpec{
			Vpc:        options.VPC,
			Protocol:   "IPv4",
			CIDRBlock:  options.CIDR,
			Gateway:    options.Gateway,
			ExcludeIps: options.ExcludeIPs,
			Provider:   options.Provider,
		},
	}
	subnet, err = ovnClient.KubeovnV1().Subnets().Create(ctx, subnet, metav1.CreateOptions{})
	if err != nil {
		if apierror.IsAlreadyExists(err) {
			return GetSubnet(ctx, options.Name)
		}
		log.Logger.Warningf("Failed to create Subnet: %s", err)
		return nil, model.RetVal{Msg: i18n.K8S.CreateError, Attr: map[string]any{"Model": "Subnet", "Error": err.Error()}}
	}
	return subnet, model.SuccessRetVal()
}

func GetSubnet(ctx context.Context, name string) (*kubeovnv1.Subnet, model.RetVal) {
	subnet, err := ovnClient.KubeovnV1().Subnets().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if apierror.IsNotFound(err) {
			return nil, model.RetVal{Msg: i18n.K8S.NotFound, Attr: map[string]any{"Model": "Subnet"}}
		}
		log.Logger.Warningf("Failed to get Subnet: %s", err)
		return nil, model.RetVal{Msg: i18n.K8S.GetError, Attr: map[string]any{"Model": "Subnet", "Error": err.Error()}}
	}
	return subnet, model.SuccessRetVal()
}

func DeleteSubnetList(ctx context.Context, labels ...map[string]string) model.RetVal {
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
	err := ovnClient.KubeovnV1().Subnets().DeleteCollection(ctx, metav1.DeleteOptions{}, options)
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete Subnet: %s", err)
		return model.RetVal{Msg: i18n.K8S.DeleteError, Attr: map[string]any{"Model": "Subnet", "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}
