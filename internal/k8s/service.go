package k8s

import (
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"time"
)

func CreateService(ctx context.Context, docker model.Docker, usage model.Usage) (*corev1.Service, bool, string) {
	var (
		service *corev1.Service
		err     error
	)
	service = &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      docker.ServiceName,
			Namespace: NamespaceName,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": docker.PodName,
			},
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       usage.Port,
					TargetPort: intstr.FromInt32(usage.Port),
				},
			},
			Type:                  corev1.ServiceTypeNodePort,
			ExternalIPs:           config.Env.K8S.Nodes,
			ExternalTrafficPolicy: corev1.ServiceExternalTrafficPolicyTypeLocal,
		},
	}
	service, err = Client.CoreV1().Services(NamespaceName).Create(ctx, service, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Error creating service: %s", err)
		return nil, false, "CreateServiceError"
	}
	return service, true, "CreateService"
}

// DeleteService 删除 Service, 目前主要是靶机的端口映射
func DeleteService(name string) (bool, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	err := Client.CoreV1().Services(NamespaceName).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to delete service: %v", err)
		return false, "DeleteServiceError"
	}
	return true, "Success"
}
