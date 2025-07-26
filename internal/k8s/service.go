package k8s

import (
	"CBCTF/internal/config"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strings"
)

type CreateServiceOptions struct {
	Name     string
	Labels   map[string]string
	Ports    model.Exposes
	Selector map[string]string
}

func CreateService(ctx context.Context, options CreateServiceOptions) (*corev1.Service, bool, string) {
	var (
		service *corev1.Service
		err     error
	)
	service = &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      options.Name,
			Namespace: GlobalNamespace,
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
			ExternalIPs:           config.Env.K8S.Nodes,
			ExternalTrafficPolicy: corev1.ServiceExternalTrafficPolicyTypeLocal,
		},
	}
	service, err = kubeClient.CoreV1().Services(GlobalNamespace).Create(ctx, service, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create Service: %v", err)
		return nil, false, i18n.CreateServiceError
	}
	return service, true, i18n.Success
}

func GetServiceList(ctx context.Context, labels ...map[string]string) (*corev1.ServiceList, bool, string) {
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
	serviceList, err := kubeClient.CoreV1().Services(GlobalNamespace).List(ctx, options)
	if err != nil {
		if apierror.IsNotFound(err) {
			return nil, false, i18n.ServiceNotFound
		}
		log.Logger.Warningf("Failed to list Service: %v", err)
		return nil, false, i18n.GetServiceError
	}
	return serviceList, true, i18n.Success
}

// DeleteService 删除 Service, 目前主要是靶机的端口映射
func DeleteService(ctx context.Context, name string) (bool, string) {
	err := kubeClient.CoreV1().Services(GlobalNamespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete Service %s: %v", name, err)
		return false, i18n.DeleteServiceError
	}
	return true, i18n.Success
}

// DeleteServiceList Service 不支持 DeleteCollection
func DeleteServiceList(ctx context.Context, labels ...map[string]string) (bool, string) {
	services, ok, msg := GetServiceList(ctx, labels...)
	if !ok {
		return false, msg
	}
	for _, service := range services.Items {
		if ok, msg = DeleteService(ctx, service.Name); !ok {
			return false, msg
		}
	}
	return true, i18n.Success
}
