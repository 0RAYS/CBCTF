package k8s

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"fmt"
	"strings"

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

func CreateEndpoint(ctx context.Context, options CreateEndpointOptions) (*discoveryv1.EndpointSlice, model.RetVal) {
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
		return nil, model.RetVal{Msg: i18n.K8S.CreateError, Attr: map[string]any{"Model": "EndpointSlice", "Error": err.Error()}}
	}
	return endpoint, model.SuccessRetVal()
}

func GetEndpoint(ctx context.Context, name string) (*discoveryv1.EndpointSlice, model.RetVal) {
	endpoint, err := kubeClient.DiscoveryV1().EndpointSlices(globalNamespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if apierror.IsNotFound(err) {
			return nil, model.RetVal{Msg: i18n.K8S.NotFound, Attr: map[string]any{"Model": "EndpointSlice"}}
		}
		log.Logger.Warningf("Failed to get EndpointSlice: %s", err)
		return nil, model.RetVal{Msg: i18n.K8S.GetError, Attr: map[string]any{"Model": "EndpointSlice", "Error": err.Error()}}
	}
	return endpoint, model.SuccessRetVal()
}

func GetEndpointList(ctx context.Context, labels ...map[string]string) (*discoveryv1.EndpointSliceList, model.RetVal) {
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
	endpoints, err := kubeClient.DiscoveryV1().EndpointSlices(globalNamespace).List(ctx, options)
	if err != nil {
		if apierror.IsNotFound(err) {
			return nil, model.RetVal{Msg: i18n.K8S.NotFound, Attr: map[string]any{"Model": "EndpointSlice"}}
		}
		log.Logger.Warningf("Failed to get EndpointSlice: %s", err)
		return nil, model.RetVal{Msg: i18n.K8S.GetError, Attr: map[string]any{"Model": "EndpointSlice", "Error": err.Error()}}
	}
	return endpoints, model.SuccessRetVal()
}

func DeleteEndpoint(ctx context.Context, name string) model.RetVal {
	err := kubeClient.DiscoveryV1().EndpointSlices(globalNamespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to delete EndpointSlice for %s", name)
		return model.RetVal{Msg: i18n.K8S.DeleteError, Attr: map[string]any{"Model": "EndpointSlice", "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

func DeleteEndpointList(ctx context.Context, labels ...map[string]string) model.RetVal {
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
	err := kubeClient.DiscoveryV1().EndpointSlices(globalNamespace).DeleteCollection(ctx, metav1.DeleteOptions{}, options)
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete EndpointSlice: %s", err)
		return model.RetVal{Msg: i18n.K8S.DeleteError, Attr: map[string]any{"Model": "EndpointSlice", "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}
