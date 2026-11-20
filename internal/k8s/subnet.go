package k8s

import (
	"context"
	"fmt"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
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
		log.Logger.Warningf("Failed to create Subnet %s: %s", options.Name, err)
		return nil, model.FailRetVal(i18n.K8S.CreateError, "Subnet", options.Name, err)
	}
	return subnet, model.SuccessRetVal()
}

func DeleteSubnetCollection(ctx context.Context, labels map[string]string) model.RetVal {
	return deleteCollection("Subnet", "Subnets", labels, func(options metav1.ListOptions) error {
		return ovnClient.KubeovnV1().Subnets().DeleteCollection(ctx, metav1.DeleteOptions{}, options)
	})
}
