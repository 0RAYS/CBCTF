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

type CreateSNatOptions struct {
	Name         string
	Labels       map[string]string
	EIP          string
	InternalCIDR string
}

func CreateSNat(ctx context.Context, options CreateSNatOptions) (*kubeovnv1.IptablesSnatRule, model.RetVal) {
	var (
		snat *kubeovnv1.IptablesSnatRule
		err  error
	)
	snat = &kubeovnv1.IptablesSnatRule{
		ObjectMeta: metav1.ObjectMeta{
			Name:   options.Name,
			Labels: options.Labels,
		},
		Spec: kubeovnv1.IptablesSnatRuleSpec{
			EIP:          options.EIP,
			InternalCIDR: options.InternalCIDR,
		},
	}
	snat, err = kubeOVNClient.KubeovnV1().IptablesSnatRules().Create(ctx, snat, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warnf("Failed to create iptables SnatRule: %s", err)
		return nil, model.RetVal{Msg: i18n.K8S.CreateError, Attr: map[string]any{"Model": "SnatRule", "Error": err.Error()}}
	}
	return snat, model.SuccessRetVal()
}

func GetSNat(ctx context.Context, name string) (*kubeovnv1.IptablesSnatRule, model.RetVal) {
	snat, err := kubeOVNClient.KubeovnV1().IptablesSnatRules().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		log.Logger.Warnf("Failed to get iptables SnatRule: %s", err)
		return nil, model.RetVal{Msg: i18n.K8S.GetError, Attr: map[string]any{"Model": "SnatRule", "Error": err.Error()}}
	}
	return snat, model.SuccessRetVal()
}

func GetSNatList(ctx context.Context, labels ...map[string]string) (*kubeovnv1.IptablesSnatRuleList, model.RetVal) {
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
	snats, err := kubeOVNClient.KubeovnV1().IptablesSnatRules().List(ctx, options)
	if err != nil {
		log.Logger.Warnf("Failed to list iptables SnatRules: %s", err)
		return nil, model.RetVal{Msg: i18n.K8S.GetError, Attr: map[string]any{"Model": "SnatRule", "Error": err.Error()}}
	}
	return snats, model.SuccessRetVal()
}

func DeleteSNat(ctx context.Context, name string) model.RetVal {
	err := kubeOVNClient.KubeovnV1().IptablesSnatRules().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warnf("Failed to delete iptables SnatRule: %s", err)
		return model.RetVal{Msg: i18n.K8S.DeleteError, Attr: map[string]any{"Model": "SnatRule", "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

func DeleteSNatList(ctx context.Context, labels ...map[string]string) model.RetVal {
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
	err := kubeOVNClient.KubeovnV1().IptablesSnatRules().DeleteCollection(ctx, metav1.DeleteOptions{}, options)
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warnf("Failed to delete iptables SnatRule: %s", err)
		return model.RetVal{Msg: i18n.K8S.DeleteError, Attr: map[string]any{"Model": "SnatRule", "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}
