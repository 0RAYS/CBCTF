package k8s

import (
	"context"
	"strings"

	corev1 "k8s.io/api/core/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"

	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
)

type CreateServiceOptions struct {
	ClusterIP bool
	Name      string
	Labels    map[string]string
	Ports     model.Exposes
	Selector  map[string]string
}

func CreateService(ctx context.Context, options CreateServiceOptions) (*corev1.Service, model.RetVal) {
	ports := make([]corev1.ServicePort, 0, len(options.Ports))
	for _, port := range options.Ports {
		ports = append(ports, corev1.ServicePort{
			Name:       utils.UUID(),
			Protocol:   corev1.Protocol(strings.ToUpper(port.Protocol)),
			Port:       port.Port,
			TargetPort: intstr.FromInt32(port.Port),
		})
	}
	service := &corev1.Service{
		OwnerReferences: resourceOwners(ctx),
		Name:            options.Name,
		Namespace:       globalNamespace,
		Labels:          options.Labels,
		Spec: corev1.ServiceSpec{
			Selector:              options.Selector,
			Ports:                 ports,
			Type:                  corev1.ServiceTypeNodePort,
			ExternalTrafficPolicy: corev1.ServiceExternalTrafficPolicyTypeLocal,
		},
	}
	if options.ClusterIP {
		service.Spec.Type = corev1.ServiceTypeClusterIP
		service.Spec.ExternalTrafficPolicy = ""
	}
	service, err := kubeClient.CoreV1().Services(globalNamespace).Create(ctx, service, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create Service: %s", err)
		return nil, model.RetVal{Msg: i18n.K8S.CreateError, Attr: map[string]any{"Model": "Service", "Error": err.Error()}}
	}
	return service, model.SuccessRetVal()
}

func ListServices(ctx context.Context, selectors ...map[string]string) (*corev1.ServiceList, model.RetVal) {
	var options metav1.ListOptions
	if len(selectors) > 0 {
		options.LabelSelector = labels.SelectorFromSet(selectors[0]).String()
	}
	serviceList, err := kubeClient.CoreV1().Services(globalNamespace).List(ctx, options)
	if err != nil {
		if apierror.IsNotFound(err) {
			return nil, model.RetVal{Msg: i18n.K8S.NotFound, Attr: map[string]any{"Model": "Service"}}
		}
		log.Logger.Warningf("Failed to list Service: %s", err)
		return nil, model.RetVal{Msg: i18n.K8S.GetError, Attr: map[string]any{"Model": "Service", "Error": err.Error()}}
	}
	return serviceList, model.SuccessRetVal()
}

// DeleteService 删除 Service, 目前主要是靶机的端口映射
func DeleteService(ctx context.Context, name string) model.RetVal {
	err := kubeClient.CoreV1().Services(globalNamespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete Service %s: %s", name, err)
		return model.RetVal{Msg: i18n.K8S.DeleteError, Attr: map[string]any{"Model": "Service", "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

// DeleteServiceCollection Service 不支持 DeleteCollection
func DeleteServiceCollection(ctx context.Context, labels map[string]string) model.RetVal {
	if _, ret := deleteCollectionOptions("Service", labels); !ret.OK {
		return ret
	}
	services, ret := ListServices(ctx, labels)
	if !ret.OK || services == nil {
		return ret
	}
	batch := model.NewBatch(len(services.Items))
	for _, service := range services.Items {
		if ctx.Err() != nil {
			return batch.Result(ctx)
		}
		batch.Record(service.Name, "delete_service", DeleteService(ctx, service.Name))
	}
	return batch.Result(ctx)
}
