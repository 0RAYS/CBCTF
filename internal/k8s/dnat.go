package k8s

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"fmt"
	"strings"

	kubeovnv1 "github.com/JBNRZ/kubeovn-api/pkg/apis/kubeovn/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CreateDNatOptions struct {
	Name         string
	Labels       map[string]string
	EIP          string
	ExternalPort string
	InternalIP   string
	InternalPort string
	Protocol     string
}

func CreateDNat(ctx context.Context, options CreateDNatOptions) (*kubeovnv1.IptablesDnatRule, model.RetVal) {
	var (
		dnat *kubeovnv1.IptablesDnatRule
		err  error
	)
	dnat = &kubeovnv1.IptablesDnatRule{
		ObjectMeta: metav1.ObjectMeta{
			Name:   options.Name,
			Labels: options.Labels,
		},
		Spec: kubeovnv1.IptablesDnatRuleSpec{
			EIP:          options.EIP,
			ExternalPort: options.ExternalPort,
			InternalIP:   options.InternalIP,
			InternalPort: options.InternalPort,
			Protocol:     options.Protocol,
		},
	}
	dnat, err = kubeOVNClient.KubeovnV1().IptablesDnatRules().Create(ctx, dnat, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create iptables DnatRule: %s", err)
		return nil, model.RetVal{Msg: i18n.K8S.CreateError, Attr: map[string]any{"Model": "DnatRule", "Error": err.Error()}}
	}
	return dnat, model.SuccessRetVal()
}

func GetDNat(ctx context.Context, name string) (*kubeovnv1.IptablesDnatRule, model.RetVal) {
	dnat, err := kubeOVNClient.KubeovnV1().IptablesDnatRules().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if apierror.IsNotFound(err) {
			return nil, model.RetVal{Msg: i18n.K8S.NotFound, Attr: map[string]any{"Model": "DnatRule"}}
		}
		log.Logger.Warningf("Failed to get iptables DnatRule: %s", err)
		return nil, model.RetVal{Msg: i18n.K8S.GetError, Attr: map[string]any{"Model": "DnatRule", "Error": err.Error()}}
	}
	return dnat, model.SuccessRetVal()
}

func GetDNatList(ctx context.Context, labels ...map[string]string) (*kubeovnv1.IptablesDnatRuleList, model.RetVal) {
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
	dnats, err := kubeOVNClient.KubeovnV1().IptablesDnatRules().List(ctx, options)
	if err != nil {
		log.Logger.Warningf("Failed to list iptables DnatRules: %s", err)
		return nil, model.RetVal{Msg: i18n.K8S.GetError, Attr: map[string]any{"Model": "DnatRule", "Error": err.Error()}}
	}
	return dnats, model.SuccessRetVal()
}

func DeleteDNat(ctx context.Context, name string) model.RetVal {
	err := kubeOVNClient.KubeovnV1().IptablesDnatRules().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete iptables DnatRule: %s", err)
		return model.RetVal{Msg: i18n.K8S.DeleteError, Attr: map[string]any{"Model": "DnatRule", "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

func DeleteDNatList(ctx context.Context, labels ...map[string]string) model.RetVal {
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
	err := kubeOVNClient.KubeovnV1().IptablesDnatRules().DeleteCollection(ctx, metav1.DeleteOptions{}, options)
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete iptables DnatRule: %s", err)
		return model.RetVal{Msg: i18n.K8S.DeleteError, Attr: map[string]any{"Model": "DnatRule", "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}
