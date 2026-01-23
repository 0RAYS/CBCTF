package k8s

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	kubeovnv1 "github.com/JBNRZ/kubeovn-api/pkg/apis/kubeovn/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CreateEIPOptions struct {
	Name           string
	Labels         map[string]string
	NatGw          string
	ExternalSubnet string
}

func CreateEIP(ctx context.Context, options CreateEIPOptions) (*kubeovnv1.IptablesEIP, model.RetVal) {
	var (
		eip *kubeovnv1.IptablesEIP
		ret model.RetVal
		err error
	)
	eip = &kubeovnv1.IptablesEIP{
		ObjectMeta: metav1.ObjectMeta{
			Name:      options.Name,
			Namespace: globalNamespace,
			Labels:    options.Labels,
		},
		Spec: kubeovnv1.IptablesEIPSpec{
			NatGwDp:        options.NatGw,
			ExternalSubnet: options.ExternalSubnet,
		},
	}
	eip, err = kubeOVNClient.KubeovnV1().IptablesEIPs().Create(ctx, eip, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create EIP: %s", err)
		return nil, model.RetVal{Msg: i18n.K8S.CreateError, Attr: map[string]any{"Model": "EIP", "Error": err.Error()}}
	}
	for {
		eip, ret = GetEIP(ctx, options.Name)
		if !ret.OK {
			return nil, ret
		}
		if eip != nil && net.ParseIP(eip.Spec.V4ip) != nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	return eip, model.SuccessRetVal()
}

func GetEIP(ctx context.Context, name string) (*kubeovnv1.IptablesEIP, model.RetVal) {
	eip, err := kubeOVNClient.KubeovnV1().IptablesEIPs().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if apierror.IsNotFound(err) {
			return nil, model.RetVal{Msg: i18n.K8S.NotFound, Attr: map[string]any{"Model": "EIP"}}
		}
		log.Logger.Warningf("Failed to get EIP: %s", err)
		return nil, model.RetVal{Msg: i18n.K8S.GetError, Attr: map[string]any{"Model": "EIP", "Error": err.Error()}}
	}
	return eip, model.SuccessRetVal()
}

func GetEIPList(ctx context.Context, labels ...map[string]string) (*kubeovnv1.IptablesEIPList, model.RetVal) {
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
		log.Logger.Warningf("Failed to get EIP list: %s", err)
		return nil, model.RetVal{Msg: i18n.K8S.GetError, Attr: map[string]any{"Model": "EIP", "Error": err.Error()}}
	}
	return eips, model.SuccessRetVal()
}

func DeleteEIP(ctx context.Context, name string) model.RetVal {
	err := kubeOVNClient.KubeovnV1().IptablesEIPs().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete EIP: %s", err)
		return model.RetVal{Msg: i18n.K8S.DeleteError, Attr: map[string]any{"Model": "EIP", "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

func DeleteEIPList(ctx context.Context, labels ...map[string]string) model.RetVal {
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
		log.Logger.Warningf("Failed to delete EIP: %s", err)
		return model.RetVal{Msg: i18n.K8S.DeleteError, Attr: map[string]any{"Model": "EIP", "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}
