package k8s

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"context"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type CreateServiceOptions struct {
	Name     string
	Labels   map[string]string
	Ports    model.Exposes
	Selector map[string]string
}

func CreateService(ctx context.Context, options CreateServiceOptions) (*corev1.Service, model.RetVal) {
	var (
		service *corev1.Service
		err     error
	)
	service = &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      options.Name,
			Namespace: globalNamespace,
			Labels:    options.Labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: options.Selector,
			Ports: func() []corev1.ServicePort {
				tmp := make([]corev1.ServicePort, 0)
				for _, p := range options.Ports {
					tmp = append(tmp, corev1.ServicePort{
						Name:       utils.UUID(),
						Protocol:   corev1.Protocol(strings.ToUpper(p.Protocol)),
						Port:       p.Port,
						TargetPort: intstr.FromInt32(p.Port),
					})
				}
				return tmp
			}(),
			Type:                  corev1.ServiceTypeNodePort,
			ExternalTrafficPolicy: corev1.ServiceExternalTrafficPolicyTypeLocal,
		},
	}
	service, err = kubeClient.CoreV1().Services(globalNamespace).Create(ctx, service, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create Service: %s", err)
		return nil, model.RetVal{Msg: i18n.K8S.CreateError, Attr: map[string]any{"Model": "Service", "Error": err.Error()}}
	}
	return service, model.SuccessRetVal()
}

func GetServiceList(ctx context.Context, labels ...map[string]string) (*corev1.ServiceList, model.RetVal) {
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

// DeleteServiceList Service 不支持 DeleteCollection
func DeleteServiceList(ctx context.Context, labels ...map[string]string) model.RetVal {
	services, ret := GetServiceList(ctx, labels...)
	if !ret.OK || services == nil {
		return ret
	}
	for _, service := range services.Items {
		if ret = DeleteService(ctx, service.Name); !ret.OK {
			return ret
		}
	}
	return model.SuccessRetVal()
}
