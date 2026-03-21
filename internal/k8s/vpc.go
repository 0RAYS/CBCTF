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

type CreateVPCOptions struct {
	Name         string
	Labels       map[string]string
	StaticRoutes []*kubeovnv1.StaticRoute
	PolicyRoutes []*kubeovnv1.PolicyRoute
}

func CreateVPC(ctx context.Context, options CreateVPCOptions) (*kubeovnv1.Vpc, model.RetVal) {
	var (
		vpc *kubeovnv1.Vpc
		err error
	)
	vpc = &kubeovnv1.Vpc{
		ObjectMeta: metav1.ObjectMeta{
			Name:      options.Name,
			Namespace: globalNamespace,
			Labels:    options.Labels,
		},
		Spec: kubeovnv1.VpcSpec{
			StaticRoutes: options.StaticRoutes,
			PolicyRoutes: options.PolicyRoutes,
		},
	}
	vpc, err = ovnClient.KubeovnV1().Vpcs().Create(ctx, vpc, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create VPC: %s", err)
		return nil, model.RetVal{Msg: i18n.K8S.CreateError, Attr: map[string]any{"Model": "VPC", "Error": err.Error()}}
	}
	return vpc, model.SuccessRetVal()
}

func GetVPC(ctx context.Context, name string) (*kubeovnv1.Vpc, model.RetVal) {
	vpc, err := ovnClient.KubeovnV1().Vpcs().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if apierror.IsNotFound(err) {
			return nil, model.RetVal{Msg: i18n.K8S.NotFound, Attr: map[string]any{"Model": "VPC"}}
		}
		log.Logger.Warningf("Failed to get VPC: %s", err)
		return nil, model.RetVal{Msg: i18n.K8S.GetError, Attr: map[string]any{"Model": "VPC", "Error": err.Error()}}
	}
	return vpc, model.SuccessRetVal()
}

func GetVPCList(ctx context.Context, labels ...map[string]string) (*kubeovnv1.VpcList, model.RetVal) {
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
	vpcList, err := ovnClient.KubeovnV1().Vpcs().List(ctx, options)
	if err != nil {
		log.Logger.Warningf("Failed to list VPC: %s", err)
		return nil, model.RetVal{Msg: i18n.K8S.GetError, Attr: map[string]any{"Model": "VPC", "Error": err.Error()}}
	}
	return vpcList, model.SuccessRetVal()
}

func DeleteVPC(ctx context.Context, name string) model.RetVal {
	err := ovnClient.KubeovnV1().Vpcs().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete VPC: %s", err)
		return model.RetVal{Msg: i18n.K8S.DeleteError, Attr: map[string]any{"Model": "VPC", "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

func DeleteVPCList(ctx context.Context, labels ...map[string]string) model.RetVal {
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
	err := ovnClient.KubeovnV1().Vpcs().DeleteCollection(ctx, metav1.DeleteOptions{}, options)
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete VPC: %s", err)
		return model.RetVal{Msg: i18n.K8S.DeleteError, Attr: map[string]any{"Model": "VPC", "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}
