package k8s

import (
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"
)

func StartContainer(challenge model.Challenge, flag model.Flag, docker model.Docker) (int32, bool, string) {
	var err error
	if challenge.Type != model.Container {
		return -1, false, "InvalidChallengeType"
	}
	if challenge.DockerImage == "" {
		return -1, false, "EmptyDockerImage"
	}
	log.Logger.Debugf("Creating pod for challenge %s:%s", challenge.Name, challenge.ID)
	env := []corev1.EnvVar{
		{
			Name:  "FLAG",
			Value: flag.Value,
		},
	}
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      docker.PodName,
			Namespace: NamespaceName,
			Labels: map[string]string{
				"app": docker.PodName,
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  docker.ContainerName,
					Image: challenge.DockerImage,
					Env:   env,
					Ports: []corev1.ContainerPort{
						{
							ContainerPort: challenge.Port,
						},
					},
				},
			},
			TerminationGracePeriodSeconds: ptr.To[int64](3),
			RestartPolicy:                 corev1.RestartPolicyNever,
		},
	}
	pod, err = Client.CoreV1().Pods(NamespaceName).Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create pod: %v", err)
		return -1, false, "CreatePodError"
	}
	service := &corev1.Service{
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
					Port:       challenge.Port,
					TargetPort: intstr.FromInt32(challenge.Port),
				},
			},
			Type: corev1.ServiceTypeNodePort,
			ExternalIPs: []string{
				config.Env.K8S.Master,
			},
		},
	}
	service, err = Client.CoreV1().Services(NamespaceName).Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Error creating service: %s", err)
		return -1, false, "CreateServiceError"
	}
	for {
		pod, err = Client.CoreV1().Pods(NamespaceName).Get(context.TODO(), docker.PodName, metav1.GetOptions{})
		if err != nil {
			log.Logger.Warningf("Failed to get pod: %v", err)
			return -1, false, "GetPodError"
		}
		if pod.Status.Phase == corev1.PodRunning {
			break
		}
		if pod.Status.Phase != corev1.PodPending {
			log.Logger.Warningf("Pod %s:%s failed to run", challenge.Name, pod.Name)
			return -1, false, "CreatePodError"
		}
	}
	port := service.Spec.Ports[0].NodePort
	return port, true, "Success"
}

func StopContainer(docker model.Docker) (bool, string) {
	var err error
	err = Client.CoreV1().Services(NamespaceName).Delete(context.TODO(), docker.ServiceName, metav1.DeleteOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to delete service: %v", err)
		return false, "DeleteServiceError"
	}
	err = Client.CoreV1().Pods(NamespaceName).Delete(context.TODO(), docker.PodName, metav1.DeleteOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to delete pod: %v", err)
		return false, "DeletePodError"
	}
	return true, "Success"
}
