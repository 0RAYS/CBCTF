package k8s

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"context"

	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CreateEndpointOptions struct {
	Name     string
	Labels   map[string]string
	IP       string
	Port     int32
	Protocol string
}

func CreateEndpoint(ctx context.Context, options CreateEndpointOptions) (*discoveryv1.EndpointSlice, bool, string) {
	var (
		endpoint *discoveryv1.EndpointSlice
		err      error
	)
	endpoint = &discoveryv1.EndpointSlice{
		ObjectMeta: metav1.ObjectMeta{
			Name:   options.Name,
			Labels: options.Labels,
		},
		Endpoints: []discoveryv1.Endpoint{
			{
				Addresses: []string{options.IP},
			},
		},
		Ports: []discoveryv1.EndpointPort{
			{
				Protocol: (*corev1.Protocol)(&options.Protocol),
				Port:     &options.Port,
			},
		},
	}
	endpoint, err = kubeClient.DiscoveryV1().EndpointSlices(globalNamespace).Create(ctx, endpoint, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create EndpointSlice for %s", err)
		return nil, false, i18n.CreateEndpointError
	}
	return endpoint, true, i18n.Success
}

func GetEndpoint(ctx context.Context, name string) (*discoveryv1.EndpointSlice, bool, string) {
	endpoints, err := kubeClient.DiscoveryV1().EndpointSlices(globalNamespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if apierror.IsNotFound(err) {
			return nil, false, i18n.EndpointNotFound
		}
		log.Logger.Warningf("Failed to get EndpointSlice for %s", name)
		return nil, false, i18n.GetEndpointError
	}
	return endpoints, true, i18n.Success
}

func DeleteEndpoint(ctx context.Context, name string) (bool, string) {
	err := kubeClient.DiscoveryV1().EndpointSlices(globalNamespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to delete EndpointSlice for %s", name)
		return false, i18n.DeleteEndpointError
	}
	return true, i18n.Success
}
