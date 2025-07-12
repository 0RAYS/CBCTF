package k8s

import (
	"CBCTF/internel/config"
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"CBCTF/internel/utils"
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type CreateServiceOptions struct {
	Name     string
	Labels   map[string]string
	Ports    []int32
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
						Protocol:   corev1.ProtocolTCP,
						Port:       p,
						TargetPort: intstr.FromInt32(p),
					})
					tmp = append(tmp, corev1.ServicePort{
						Name:       utils.UUID(),
						Protocol:   corev1.ProtocolUDP,
						Port:       p,
						TargetPort: intstr.FromInt32(p),
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

func GetServiceList(ctx context.Context) (*corev1.ServiceList, bool, string) {
	serviceList, err := kubeClient.CoreV1().Services(GlobalNamespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		if apierror.IsNotFound(err) {
			return nil, false, i18n.ServiceNotFound
		}
		log.Logger.Warningf("Failed to list Service: %v", err)
		return nil, false, i18n.GetServiceError
	}
	return serviceList, true, i18n.Success
}

func GetServiceListByPodName(ctx context.Context, key, podName string) (*corev1.ServiceList, bool, string) {
	serviceList, err := kubeClient.CoreV1().Services(GlobalNamespace).List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", key, podName),
	})
	if err != nil {
		if apierror.IsNotFound(err) {
			return nil, false, i18n.ServiceNotFound
		}
		log.Logger.Warningf("Failed to list Pod %s Service: %v", podName, err)
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

// DeleteServiceListByPodName TODO: 有可能删不干净
func DeleteServiceListByPodName(ctx context.Context, key, podName string) (bool, string) {
	serviceList, ok, msg := GetServiceListByPodName(ctx, key, podName)
	if !ok {
		if msg != i18n.ServiceNotFound {
			return false, msg
		}
		return true, i18n.Success
	}
	for _, svc := range serviceList.Items {
		if ok, msg = DeleteService(ctx, svc.Name); !ok {
			return false, msg
		}
	}
	return true, i18n.Success
}
