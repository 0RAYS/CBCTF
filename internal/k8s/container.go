package k8s

import (
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	corev1 "k8s.io/api/core/v1"
	"math/rand"
	"strconv"
	"time"
)

// StartContainer 启动容器, 并且注入 flag
func StartContainer(usage model.Usage, flag model.Flag, container model.Container) (string, int32, bool, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	var (
		service *corev1.Service
		pod     *corev1.Pod
		ok      bool
		ip      string
		msg     string
	)
	if usage.Type != model.Docker {
		return "", -1, false, "InvalidChallengeType"
	}
	if usage.DockerImage == "" {
		return "", -1, false, "EmptyContainerImage"
	}
	log.Logger.Debugf("Creating Pod for challenge %s:%s", usage.Name, usage.ChallengeID)
	service, ok, msg = CreateService(ctx, container, usage)
	if !ok {
		return "", -1, false, msg
	}
	_, ok, msg = CreateNetworkPolicy(ctx, container, usage)
	if !ok {
		return "", -1, false, msg
	}
	port := service.Spec.Ports[0].NodePort
	containers := []corev1.Container{
		{
			Name:  container.ContainerName,
			Image: usage.DockerImage,
			Env: []corev1.EnvVar{
				{
					Name:  "FLAG",
					Value: flag.Value,
				},
			},
			Ports: []corev1.ContainerPort{
				{
					ContainerPort: usage.Port,
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
					Value: container.PodName,
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
					Value: strconv.Itoa(int(usage.Port)),
				},
				{
					Name:  "remotePort",
					Value: strconv.Itoa(int(port)),
				},
			},
		}
		containers = append(containers, frpc)
		ip = frps.Host
		log.Logger.Infof("Frpc started: %s:%d -> %s:%d", frps.Host, port, container.PodName, port)
	}
	pod, ok, msg = CreatePod(ctx, usage, container.PodName, containers)
	if !ok {
		return "", -1, false, msg
	}
	if !config.Env.K8S.Frpc.On {
		ip = pod.Status.HostIP
	}
	log.Logger.Infof("Pod %s:%s is running on %s:%d", usage.Name, pod.Name, ip, port)
	return ip, port, true, "Success"
}

// StopContainer 停止容器
func StopContainer(container model.Container) (bool, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	var err error
	err = CopyFromPod(
		container.PodName, "tcpdump", "/root/traffic.pcap",
		container.TrafficPath(),
	)
	if err != nil {
		log.Logger.Warningf("Failed to copy %d traffic: %v", container.TeamID, err)
	}
	if ok, msg := DeleteNetworkPolicy(ctx, container.NetworkPolicyName); !ok {
		return false, msg
	}
	if ok, msg := DeleteService(ctx, container.ServiceName); !ok {
		return false, msg
	}
	return DeletePod(ctx, container.PodName)
}
