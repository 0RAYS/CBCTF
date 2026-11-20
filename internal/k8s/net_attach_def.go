package k8s

import (
	"context"
	"fmt"

	netattv1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
)

const OVNNetworkConfigTmpl = `{
	"cniVersion": "0.3.0",
	"type": "kube-ovn",
	"server_socket": "/run/openvswitch/kube-ovn-daemon.sock",
	"provider": "%s.%s.ovn"
}`

type CreateNetAttachDefOptions struct {
	Name   string
	Labels map[string]string
}

func CreateNetAttachDef(ctx context.Context, options CreateNetAttachDefOptions) (*netattv1.NetworkAttachmentDefinition, model.RetVal) {
	var (
		netAttachDef *netattv1.NetworkAttachmentDefinition
		err          error
	)
	netAttachDef = &netattv1.NetworkAttachmentDefinition{
		OwnerReferences: resourceOwners(ctx),
		Name:            options.Name,
		Namespace:       globalNamespace,
		Labels:          options.Labels,
		Spec: netattv1.NetworkAttachmentDefinitionSpec{
			Config: fmt.Sprintf(OVNNetworkConfigTmpl, options.Name, globalNamespace),
		},
	}
	netAttachDef, err = netattClient.K8sCniCncfIoV1().NetworkAttachmentDefinitions(globalNamespace).Create(ctx, netAttachDef, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create NetworkAttachmentDefinition %s: %s", options.Name, err)
		return nil, model.FailRetVal(i18n.K8S.CreateError, "NetworkAttachmentDefinition", options.Name, err)
	}
	return netAttachDef, model.SuccessRetVal()
}

func DeleteNetAttachDefCollection(ctx context.Context, namespace string, labels map[string]string) model.RetVal {
	if namespace == "" {
		return model.FailRetVal(i18n.K8S.DeleteError, "NetworkAttachmentDefinition", "", errMissingNamespace)
	}
	return deleteCollection("NetworkAttachmentDefinition", "NetworkAttachmentDefinitions", labels, func(options metav1.ListOptions) error {
		return netattClient.K8sCniCncfIoV1().NetworkAttachmentDefinitions(namespace).DeleteCollection(ctx, metav1.DeleteOptions{}, options)
	})
}
