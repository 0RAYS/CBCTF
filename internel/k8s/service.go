package k8s

import (
	"CBCTF/internel/config"
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"CBCTF/internel/utils"
	"context"
	corev1 "k8s.io/api/core/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func GetService(ctx context.Context, name string) (*corev1.Service, bool, string) {
	service, err := client.CoreV1().Services(NamespaceName).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if apierror.IsNotFound(err) {
			return nil, false, "ServiceNotFound"
		}
		log.Logger.Warningf("Failed to get Service %s: %v", name, err)
		return nil, false, i18n.GetServiceError
	}
	return service, true, i18n.Success
}

func CreateService(ctx context.Context, pod model.Pod) (*corev1.Service, bool, string) {
	var (
		service *corev1.Service
		err     error
	)
	if _, ok, _ := GetService(ctx, pod.ServiceName); ok {
		DeleteService(ctx, pod.ServiceName)
	}
	service = &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pod.ServiceName,
			Namespace: NamespaceName,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"victim": pod.Name,
			},
			Ports: func() []corev1.ServicePort {
				tmp := make([]corev1.ServicePort, 0)
				for _, p := range pod.ExposePorts {
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
	service, err = client.CoreV1().Services(NamespaceName).Create(ctx, service, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create Service %s: %s", pod.ServiceName, err)
		return nil, false, i18n.CreateServiceError
	}
	return service, true, i18n.Success
}

// DeleteService 删除 Service, 目前主要是靶机的端口映射
func DeleteService(ctx context.Context, name string) (bool, string) {
	err := client.CoreV1().Services(NamespaceName).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete Service %s: %v", name, err)
		return false, i18n.DeleteServiceError
	}
	return true, i18n.Success
}
