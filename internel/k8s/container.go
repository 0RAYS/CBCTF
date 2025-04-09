package k8s

import (
	"CBCTF/internel/config"
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"math/rand"
	"strconv"
	"time"
)

func StartContainer(container model.Container) (string, bool, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	if container.Image == "" {
		return "", false, "EmptyContainerImage"
	}
	log.Logger.Debugf("Creating Pod for usage %d:%s", container.ID, container.Image)
	service, ok, msg := CreateService(ctx, container)
	if !ok {
		return "", false, msg
	}
	for _, policy := range container.NetworkPolicies {
		_, ok, msg := CreateNetworkPolicy(ctx, container, policy)
		if !ok {
			return "", false, msg
		}
	}
	containers := []corev1.Container{
		{
			Name:  container.ContainerName,
			Image: container.Image,
			Env: func() []corev1.EnvVar {
				tmp := make([]corev1.EnvVar, 0)
				if len(container.Flags) == 1 {
					tmp = append(tmp, corev1.EnvVar{
						Name:  "FLAG",
						Value: container.Flags[0],
					})
				} else {
					for i, f := range container.Flags {
						tmp = append(tmp, corev1.EnvVar{
							Name:  fmt.Sprintf("FLAG%d", i+1),
							Value: f,
						})
					}
				}
				return tmp
			}(),
			Ports: func() []corev1.ContainerPort {
				tmp := make([]corev1.ContainerPort, 0)
				for _, p := range container.Exposes {
					tmp = append(tmp, corev1.ContainerPort{
						ContainerPort: p,
					})
				}
				return tmp
			}(),
		},
		{
			Name:    "tcpdump",
			Image:   config.Env.K8S.TCPDumpImage,
			Command: []string{"/bin/sh", "-c", "tcpdump -i any -w /root/traffic.pcap"},
		},
	}
	rand.New(rand.NewSource(time.Now().UnixNano()))
	frps := config.Env.K8S.Frpc.Frps[rand.Intn(len(config.Env.K8S.Frpc.Frps))]
	var ip string
	if config.Env.K8S.Frpc.On {
		for _, port := range service.Spec.Ports {
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
						Value: port.TargetPort.StrVal,
					},
					{
						Name:  "remotePort",
						Value: strconv.Itoa(int(port.NodePort)),
					},
				},
			}
			containers = append(containers, frpc)
			log.Logger.Infof("Frpc started: %s:%d -> %s:%s", frps.Host, port.NodePort, container.PodName, port.TargetPort.StrVal)
		}
		ip = frps.Host
	}
	pod, ok, msg := CreatePod(ctx, container.PodName, containers)
	if !ok {
		return "", false, msg
	}
	if !config.Env.K8S.Frpc.On {
		ip = pod.Status.HostIP
	}
	log.Logger.Infof("Pod %s:%s is running on %s", container.PodName, pod.Name, ip)
	return ip, true, "Success"
}

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
