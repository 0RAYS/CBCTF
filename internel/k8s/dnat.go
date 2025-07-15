package k8s

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"context"
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

func CreateDNat(ctx context.Context, options CreateDNatOptions) (*kubeovnv1.IptablesDnatRule, bool, string) {
	var (
		dnat *kubeovnv1.IptablesDnatRule
		err  error
	)
	dnat = &kubeovnv1.IptablesDnatRule{
		ObjectMeta: metav1.ObjectMeta{
			Name: options.Name,
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
		log.Logger.Warningf("Failed to create iptables DnatRule: %v", err)
		return dnat, false, i18n.CreateDNatError
	}
	return dnat, true, i18n.Success
}

func GetDNat(ctx context.Context, name string) (*kubeovnv1.IptablesDnatRule, bool, string) {
	dnat, err := kubeOVNClient.KubeovnV1().IptablesDnatRules().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if apierror.IsNotFound(err) {
			return nil, false, i18n.DNatNotFound
		}
		log.Logger.Warningf("Failed to get iptables DnatRule: %v", err)
		return nil, false, i18n.GetDNatError
	}
	return dnat, true, i18n.Success
}

func GetDNatList(ctx context.Context) (*kubeovnv1.IptablesDnatRuleList, bool, string) {
	dnats, err := kubeOVNClient.KubeovnV1().IptablesDnatRules().List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to list iptables DnatRules: %v", err)
		return nil, false, i18n.GetDNatError
	}
	return dnats, true, i18n.Success
}

func DeleteDNat(ctx context.Context, name string) (bool, string) {
	err := kubeOVNClient.KubeovnV1().IptablesDnatRules().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to delete iptables DnatRule: %v", err)
		return false, i18n.DeleteDNatError
	}
	return true, i18n.Success
}
