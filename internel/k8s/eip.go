package k8s

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"context"
	"fmt"
	kubeovnv1 "github.com/JBNRZ/kubeovn-api/pkg/apis/kubeovn/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net"
	"strings"
	"time"
)

type CreateEIPOptions struct {
	Name           string
	Labels         map[string]string
	NatGw          string
	ExternalSubnet string
}

func CreateEIP(ctx context.Context, options CreateEIPOptions) (*kubeovnv1.IptablesEIP, bool, string) {
	var (
		eip *kubeovnv1.IptablesEIP
		ok  bool
		err error
	)
	eip = &kubeovnv1.IptablesEIP{
		ObjectMeta: metav1.ObjectMeta{
			Name:      options.Name,
			Namespace: GlobalNamespace,
			Labels:    options.Labels,
		},
		Spec: kubeovnv1.IptablesEIPSpec{
			NatGwDp:        options.NatGw,
			ExternalSubnet: options.ExternalSubnet,
		},
	}
	eip, err = kubeOVNClient.KubeovnV1().IptablesEIPs().Create(ctx, eip, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create EIP: %v", err)
		return nil, false, i18n.CreateEIPError
	}
	for {
		eip, ok, _ = GetEIP(ctx, options.Name)
		if !ok {
			return nil, false, i18n.GetEIPError
		}
		if net.ParseIP(eip.Spec.V4ip) != nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	return eip, true, i18n.Success
}

func GetEIP(ctx context.Context, name string) (*kubeovnv1.IptablesEIP, bool, string) {
	eip, err := kubeOVNClient.KubeovnV1().IptablesEIPs().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if apierror.IsNotFound(err) {
			return nil, false, i18n.EIPNotFound
		}
		log.Logger.Warningf("Failed to get EIP: %v", err)
		return nil, false, i18n.GetEIPError
	}
	return eip, true, i18n.Success
}

func GetEIPList(ctx context.Context, labels ...map[string]string) (*kubeovnv1.IptablesEIPList, bool, string) {
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
	eips, err := kubeOVNClient.KubeovnV1().IptablesEIPs().List(ctx, options)
	if err != nil {
		log.Logger.Warningf("Failed to get EIP list: %v", err)
		return nil, false, i18n.GetEIPError
	}
	return eips, true, i18n.Success
}

func DeleteEIP(ctx context.Context, name string) (bool, string) {
	err := kubeOVNClient.KubeovnV1().IptablesEIPs().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete EIP: %v", err)
		return false, i18n.DeleteEIPError
	}
	return true, i18n.Success
}

func DeleteEIPByLabels(ctx context.Context, labels ...map[string]string) (bool, string) {
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
	err := kubeOVNClient.KubeovnV1().IptablesEIPs().DeleteCollection(ctx, metav1.DeleteOptions{}, options)
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete EIP: %v", err)
		return false, i18n.DeleteEIPError
	}
	return true, i18n.Success
}
