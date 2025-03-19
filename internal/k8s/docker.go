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
func StartContainer(usage model.Usage, flag model.Flag, docker model.Docker) (string, int32, bool, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	var (
		service *corev1.Service
		pod     *corev1.Pod
		ok      bool
		ip      string
		msg     string
	)
	if usage.Type != model.Container {
		return "", -1, false, "InvalidChallengeType"
	}
	if usage.DockerImage == "" {
		return "", -1, false, "EmptyDockerImage"
	}
	log.Logger.Debugf("Creating pod for challenge %s:%s", usage.Name, usage.ChallengeID)
	service, ok, msg = CreateService(ctx, docker, usage)
	if !ok {
		return "", -1, false, msg
	}
	port := service.Spec.Ports[0].NodePort
	containers := []corev1.Container{
		{
			Name:  docker.ContainerName,
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
		log.Logger.Infof("Frpc started: %s:%d -> %s:%d", frps.Host, port, docker.PodName, port)
	}
	pod, ok, msg = CreatePod(ctx, docker, usage, containers)
	if !ok {
		return "", -1, false, msg
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
	log.Logger.Infof("Pod %s:%s is running on %s:%d", usage.Name, pod.Name, ip, port)
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
