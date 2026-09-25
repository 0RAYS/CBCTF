package k8s

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CreateVPCOptions struct {
	Name   string
	Labels map[string]string
}

func CreateVPC(ctx context.Context, options CreateVPCOptions) (*kubeovnv1.Vpc, model.RetVal) {
	var (
		vpc *kubeovnv1.Vpc
		err error
	)
	vpc = &kubeovnv1.Vpc{
		Name:   options.Name,
		Labels: options.Labels,
		Spec:   kubeovnv1.VpcSpec{},
	}
	vpc, err = ovnClient.KubeovnV1().Vpcs().Create(ctx, vpc, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create VPC: %s", err)
		return nil, model.RetVal{Msg: i18n.K8S.CreateError, Attr: map[string]any{"Model": "VPC", "Error": err.Error()}}
	}
	return vpc, model.SuccessRetVal()
}

func DeleteVPCCollection(ctx context.Context, labels map[string]string) model.RetVal {
	options, ret := deleteCollectionOptions("VPC", labels)
	if !ret.OK {
		return ret
	}
	err := ovnClient.KubeovnV1().Vpcs().DeleteCollection(ctx, metav1.DeleteOptions{}, options)
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete VPC: %s", err)
		return model.RetVal{Msg: i18n.K8S.DeleteError, Attr: map[string]any{"Model": "VPC", "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}
