package k8s

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"fmt"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CreateSubnetOptions struct {
	Name         string
	Labels       map[string]string
	VPC          string
	CIDR         string
	Gateway      string
	ExcludeIPs   []string
	NetAttachDef string
}

func CreateSubnet(ctx context.Context, options CreateSubnetOptions) (*kubeovnv1.Subnet, model.RetVal) {
	var (
		subnet *kubeovnv1.Subnet
		err    error
	)
	subnet = &kubeovnv1.Subnet{
		OwnerReferences: resourceOwners(ctx),
		Name:            options.Name,
		Labels:          options.Labels,
		Spec: kubeovnv1.SubnetSpec{
			Vpc:        options.VPC,
			Protocol:   "IPv4",
			CIDRBlock:  options.CIDR,
			Gateway:    options.Gateway,
			ExcludeIps: options.ExcludeIPs,
			Provider:   fmt.Sprintf("%s.%s.ovn", options.NetAttachDef, globalNamespace),
		},
	}
	subnet, err = ovnClient.KubeovnV1().Subnets().Create(ctx, subnet, metav1.CreateOptions{})
	if err != nil {
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

func DeleteSubnetCollection(ctx context.Context, labels map[string]string) model.RetVal {
	options, ret := deleteCollectionOptions("Subnet", labels)
	if !ret.OK {
		return ret
	}
	err := ovnClient.KubeovnV1().Subnets().DeleteCollection(ctx, metav1.DeleteOptions{}, options)
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete Subnet: %s", err)
		return model.RetVal{Msg: i18n.K8S.DeleteError, Attr: map[string]any{"Model": "Subnet", "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}
