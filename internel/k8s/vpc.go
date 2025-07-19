package k8s

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"context"
	"fmt"
	kubeovnv1 "github.com/JBNRZ/kubeovn-api/pkg/apis/kubeovn/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type CreateVPCOptions struct {
	Name         string
	Labels       map[string]string
	StaticRoutes []*kubeovnv1.StaticRoute
	PolicyRoutes []*kubeovnv1.PolicyRoute
}

func CreateVPC(ctx context.Context, options CreateVPCOptions) (*kubeovnv1.Vpc, bool, string) {
	var (
		vpc *kubeovnv1.Vpc
		err error
	)
	vpc = &kubeovnv1.Vpc{
		ObjectMeta: metav1.ObjectMeta{
			Name:      options.Name,
			Namespace: GlobalNamespace,
			Labels:    options.Labels,
		},
		Spec: kubeovnv1.VpcSpec{
			StaticRoutes: options.StaticRoutes,
			PolicyRoutes: options.PolicyRoutes,
		},
	}
	vpc, err = kubeOVNClient.KubeovnV1().Vpcs().Create(ctx, vpc, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create VPC: %v", err)
		return nil, false, i18n.CreateVPCError
	}
	return vpc, true, i18n.Success
}

func GetVPC(ctx context.Context, name string) (*kubeovnv1.Vpc, bool, string) {
	vpc, err := kubeOVNClient.KubeovnV1().Vpcs().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if apierror.IsNotFound(err) {
			return nil, false, i18n.VPCNotFound
		}
		log.Logger.Warningf("Failed to get VPC: %v", err)
		return nil, false, i18n.GetVPCError
	}
	return vpc, true, i18n.Success
}

func GetVPCList(ctx context.Context, labels ...map[string]string) (*kubeovnv1.VpcList, bool, string) {
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
	vpcList, err := kubeOVNClient.KubeovnV1().Vpcs().List(ctx, options)
	if err != nil {
		log.Logger.Warningf("Failed to list VPC: %v", err)
		return nil, false, i18n.GetVPCError
	}
	return vpcList, true, i18n.Success
}

func DeleteVPC(ctx context.Context, name string) (bool, string) {
	err := kubeOVNClient.KubeovnV1().Vpcs().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete VPC: %v", err)
		return false, i18n.DeleteVPCError
	}
	return true, i18n.Success
}

func DeleteVPCByLabels(ctx context.Context, labels map[string]string) (bool, string) {
	vpcList, ok, msg := GetVPCList(ctx, labels)
	if !ok {
		return false, msg
	}
	for _, vpc := range vpcList.Items {
		if ok, msg = DeleteVPC(ctx, vpc.Name); !ok {
			return false, msg
		}
	}
	return true, i18n.Success
}
