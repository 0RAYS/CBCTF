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
	"math/rand"
	"strconv"
	"time"
)

// StartContainer 启动容器, 并且注入 flag
func StartContainer(challenge model.Challenge, flag model.Flag, docker model.Docker) (string, int32, bool, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	var (
		err error
		ok  bool
		ip  string
	)
	if challenge.Type != model.Container {
		return "", -1, false, "InvalidChallengeType"
	}
	if challenge.DockerImage == "" {
		return "", -1, false, "EmptyDockerImage"
	}
	log.Logger.Debugf("Creating pod for challenge %s:%s", challenge.Name, challenge.ID)
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
			Type:                  corev1.ServiceTypeNodePort,
			ExternalIPs:           config.Env.K8S.Nodes,
			ExternalTrafficPolicy: corev1.ServiceExternalTrafficPolicyTypeLocal,
		},
	}
	service, err = Client.CoreV1().Services(NamespaceName).Create(ctx, service, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Error creating service: %s", err)
		return "", -1, false, "CreateServiceError"
	}
	port := service.Spec.Ports[0].NodePort
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      docker.PodName,
			Namespace: NamespaceName,
			Labels: map[string]string{
				"app": docker.PodName,
			},
		},
		Spec: corev1.PodSpec{
			Containers: func() []corev1.Container {
				containers := []corev1.Container{
					{
						Name:  docker.ContainerName,
						Image: challenge.DockerImage,
						Env: []corev1.EnvVar{
							{
								Name:  "FLAG",
								Value: flag.Value,
							},
						},
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: challenge.Port,
							},
						},
					},
					{
						Name:    "tcpdump",
						Image:   config.Env.K8S.TCPDumpImage,
						Command: []string{"/bin/sh", "-c", "tcpdump -i any -w /root/traffic.pcap"},
					},
				}

				if config.Env.K8S.Frpc.On {
					rand.New(rand.NewSource(time.Now().UnixNano()))
					frps := config.Env.K8S.Frpc.Frps[rand.Intn(len(config.Env.K8S.Frpc.Frps))]
					frpc := corev1.Container{
						Name:  "frpc",
						Image: config.Env.K8S.Frpc.Image,
						Env: []corev1.EnvVar{
							{
								Name:  "serverAddr",
								Value: frps.Host,
							},
							{
								Name:  "serverPort",
								Value: strconv.Itoa(frps.Port),
							},
							{
								Name:  "token",
								Value: frps.Token,
							},
							{
								Name:  "name",
								Value: docker.PodName,
							},
							{
								Name:  "type",
								Value: "tcp",
							},
							{
								Name:  "localIP",
								Value: "127.0.0.1",
							},
							{
								Name:  "localPort",
								Value: strconv.Itoa(int(challenge.Port)),
							},
							{
								Name:  "remotePort",
								Value: strconv.Itoa(int(port)),
							},
						},
					}
					containers = append(containers, frpc)
					ip = frps.Host
					log.Logger.Infof("Frpc started: %s:%d -> %s:%d", frps.Host, port, docker.PodName, port)
				}
				return containers
			}(),
			TerminationGracePeriodSeconds: ptr.To[int64](3),
			RestartPolicy:                 corev1.RestartPolicyNever,
		},
	}
	pod, err = Client.CoreV1().Pods(NamespaceName).Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create pod: %v", err)
		return "", -1, false, "CreatePodError"
	}
	for {
		pod, ok, _ = GetPod(pod.Name)
		if !ok {
			log.Logger.Warningf("Failed to get pod: %v", err)
			return "", -1, false, "GetPodError"
		}
		if pod.Status.Phase == corev1.PodRunning {
			break
		}
		if pod.Status.Phase != corev1.PodPending {
			log.Logger.Warningf("Pod %s:%s failed to run", challenge.Name, pod.Name)
			return "", -1, false, "CreatePodError"
		}
	}
	if !config.Env.K8S.Frpc.On {
		node, ok, msg := GetNode(pod.Spec.NodeName)
		if !ok {
			return "", -1, false, msg
		}
		for _, address := range node.Status.Addresses {
			if address.Type == corev1.NodeInternalIP && address.Address != "" {
				ip = address.Address
				continue
			}
			if address.Type == corev1.NodeExternalIP && address.Address != "" {
				ip = address.Address
				break
			}
		}
	}
	log.Logger.Infof("Pod %s:%s is running on %s:%d", challenge.Name, pod.Name, ip, port)
	return ip, port, true, "Success"
}

// StopContainer 停止容器
func StopContainer(docker model.Docker) (bool, string) {
	var err error
	err = CopyFromPod(
		docker.PodName, "tcpdump", "/root/traffic.pcap",
		docker.TrafficPath(),
	)
	if err != nil {
		log.Logger.Warningf("Failed to copy %d traffic: %v", docker.TeamID, err)
	}
	if ok, msg := DeleteService(docker.ServiceName); !ok {
		return false, msg
	}
	return DeletePod(docker.PodName)
}
