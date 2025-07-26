package k8s

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"context"
	"fmt"
	kubeovnv1 "github.com/JBNRZ/kubeovn-api/pkg/apis/kubeovn/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type CreateSNatOptions struct {
	Name         string
	Labels       map[string]string
	EIP          string
	InternalCIDR string
}

func CreateSNat(ctx context.Context, options CreateSNatOptions) (*kubeovnv1.IptablesSnatRule, bool, string) {
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
		log.Logger.Warnf("Failed to create iptables SnatRule: %v", err)
		return snat, false, i18n.CreateSNatError
	}
	return snat, true, i18n.Success
}

func GetSNat(ctx context.Context, name string) (*kubeovnv1.IptablesSnatRule, bool, string) {
	snat, err := kubeOVNClient.KubeovnV1().IptablesSnatRules().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		log.Logger.Warnf("Failed to get iptables SnatRule: %v", err)
		return snat, false, i18n.GetSNatError
	}
	return snat, true, i18n.Success
}

func GetSNatList(ctx context.Context, labels ...map[string]string) (*kubeovnv1.IptablesSnatRuleList, bool, string) {
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
		log.Logger.Warnf("Failed to list iptables SnatRules: %v", err)
		return snats, false, i18n.GetSNatError
	}
	return snats, true, i18n.Success
}

func DeleteSNat(ctx context.Context, name string) (bool, string) {
	err := kubeOVNClient.KubeovnV1().IptablesSnatRules().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		log.Logger.Warnf("Failed to delete iptables SnatRule: %v", err)
		return false, i18n.DeleteSNatError
	}
	return true, i18n.Success
}

func DeleteSNatList(ctx context.Context, labels ...map[string]string) (bool, string) {
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
	if err != nil {
		log.Logger.Warnf("Failed to delete iptables SnatRule: %v", err)
		return false, i18n.DeleteSNatError
	}
	return true, i18n.Success
}
