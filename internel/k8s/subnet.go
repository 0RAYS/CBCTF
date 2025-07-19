package k8s

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"context"
	kubeovnv1 "github.com/JBNRZ/kubeovn-api/pkg/apis/kubeovn/v1"
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

func CreateSubnet(ctx context.Context, options CreateSubnetOptions) (*kubeovnv1.Subnet, bool, string) {
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
	subnet, err = kubeOVNClient.KubeovnV1().Subnets().Create(ctx, subnet, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create Subnet: %v", err)
		return nil, false, i18n.CreateSubnetError
	}
	return subnet, true, i18n.Success
}

func GetSubnet(ctx context.Context, name string) (*kubeovnv1.Subnet, bool, string) {
	subnet, err := kubeOVNClient.KubeovnV1().Subnets().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if apierror.IsNotFound(err) {
			return nil, false, i18n.SubnetNotFound
		}
		log.Logger.Warningf("Failed to get Subnet: %v", err)
		return nil, false, i18n.GetSubnetError
	}
	return subnet, true, i18n.Success
}

func GetSubnetList(ctx context.Context) (*kubeovnv1.SubnetList, bool, string) {
	subnetList, err := kubeOVNClient.KubeovnV1().Subnets().List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to list Subnet: %v", err)
		return nil, false, i18n.GetSubnetError
	}
	return subnetList, true, i18n.Success
}

func DeleteSubnet(ctx context.Context, name string) (bool, string) {
	err := kubeOVNClient.KubeovnV1().Subnets().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete Subnet: %v", err)
		return false, i18n.DeleteSubnetError
	}
	return true, i18n.Success
}
