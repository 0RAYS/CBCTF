package k8s

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"context"
	"fmt"
	netattv1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type CreateNetAttachDefOptions struct {
	Name      string
	Namespace string
	Labels    map[string]string
	Config    string
}

func CreateNetAttachDef(ctx context.Context, options CreateNetAttachDefOptions) (*netattv1.NetworkAttachmentDefinition, bool, string) {
	var (
		netAttachDef *netattv1.NetworkAttachmentDefinition
		err          error
	)
	if options.Namespace == "" {
		options.Namespace = globalNamespace
	}
	netAttachDef = &netattv1.NetworkAttachmentDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name:      options.Name,
			Namespace: options.Namespace,
			Labels:    options.Labels,
		},
		Spec: netattv1.NetworkAttachmentDefinitionSpec{
			Config: options.Config,
		},
	}
	netAttachDef, err = natattClient.K8sCniCncfIoV1().NetworkAttachmentDefinitions(options.Namespace).Create(ctx, netAttachDef, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create NetworkAttachmentDefinition: %v", err)
		return nil, false, i18n.CreateNetAttError
	}
	return netAttachDef, true, i18n.Success
}

func GetNetAttachDef(ctx context.Context, name string, namespace ...string) (*netattv1.NetworkAttachmentDefinition, bool, string) {
	if len(namespace) == 0 {
		namespace = append(namespace, globalNamespace)
	}
	netAttachDef, err := natattClient.K8sCniCncfIoV1().NetworkAttachmentDefinitions(namespace[0]).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if apierror.IsNotFound(err) {
			return nil, false, i18n.NetAttNotFound
		}
		log.Logger.Warningf("Failed to get NetworkAttachmentDefinition: %v", err)
		return nil, false, i18n.GetNetAttError
	}
	return netAttachDef, true, i18n.Success
}

func GetNetAttachDefList(ctx context.Context, labels ...map[string]string) (*netattv1.NetworkAttachmentDefinitionList, bool, string) {
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
	netAttachDefList, err := natattClient.K8sCniCncfIoV1().NetworkAttachmentDefinitions(globalNamespace).List(ctx, options)
	if err != nil {
		log.Logger.Warningf("Failed to list NetworkAttachmentDefinitions: %v", err)
		return nil, false, i18n.GetNetAttError
	}
	return netAttachDefList, true, i18n.Success
}

func DeleteNetAttachDef(ctx context.Context, name string, namespace string) (bool, string) {
	err := natattClient.K8sCniCncfIoV1().NetworkAttachmentDefinitions(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete NetworkAttachmentDefinition: %v", err)
		return false, i18n.DeleteNetAttError
	}
	return true, i18n.Success
}

func DeleteNetAttachDefList(ctx context.Context, namespace string, labels ...map[string]string) (bool, string) {
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
	err := natattClient.K8sCniCncfIoV1().NetworkAttachmentDefinitions(namespace).DeleteCollection(ctx, metav1.DeleteOptions{}, options)
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete NetworkAttachmentDefinition: %v", err)
		return false, i18n.DeleteNetAttError
	}
	return true, i18n.Success
}
