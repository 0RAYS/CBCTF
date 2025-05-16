package k8s

import (
	"CBCTF/internel/config"
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"CBCTF/internel/utils"
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"math/rand"
	"sync"
	"time"
)

// StartVictim model.Victim 需要预加载 model.Pod, 嵌套预加载 model.Container
func StartVictim(victim model.Victim, dns map[string]string) (map[string]map[string]any, bool, string) {
	log.Logger.Debugf("Creating Victim for team %d usage %d", victim.TeamID, victim.UsageID)
	type result struct {
		PodName string
		IP      string
		Ports   []int32
		OK      bool
		Msg     string
	}
	var wg sync.WaitGroup
	resultCh := make(chan result, len(victim.Pods))
	for _, pod := range victim.Pods {
		wg.Add(1)
		go func(pod model.Pod) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
			defer cancel()
			service, ok, msg := CreateService(ctx, pod)
			if !ok {
				resultCh <- result{PodName: pod.Name, OK: false, Msg: msg}
				return
			}
			containers := []corev1.Container{
				{
					Name:    "tcpdump",
					Image:   config.Env.K8S.TCPDumpImage,
					Command: []string{"/bin/sh", "-c", "tcpdump -i any -w /root/traffic.pcap"},
				},
			}
			volumes := make([]corev1.Volume, 0)
			var ip string
			if len(pod.PodPorts) > 0 && config.Env.K8S.Frpc.On {
				rand.New(rand.NewSource(time.Now().UnixNano()))
				frps := config.Env.K8S.Frpc.Frps[rand.Intn(len(config.Env.K8S.Frpc.Frps))]
				ip = frps.Host
				policy := model.NetworkPolicy{
					To: []model.Target{
						{
							CIDR: fmt.Sprintf("%s/32", frps.Host),
						},
					},
				}
				for _, p := range pod.NetworkPolicies {
					if len(p.From) > 0 {
						policy.From = []model.Target{
							{
								CIDR: fmt.Sprintf("%s/32", frps.Host),
							},
						}
						break
					}
				}
				pod.NetworkPolicies = append(pod.NetworkPolicies, policy)
				frpcConfigMap, ok, msg := CreateFrpcConfig(ctx, frps.Host, frps.Port, frps.Token, pod, service)
				if !ok {
					resultCh <- result{PodName: pod.Name, OK: false, Msg: msg}
					return
				}
				volumeName := fmt.Sprintf("%s-vol", pod.Name)
				frpc := corev1.Container{
					Name:  "frpc",
					Image: config.Env.K8S.Frpc.Image,
					Args:  []string{"-c", "/etc/frp/frpc.toml"},
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      volumeName,
							MountPath: "/etc/frp/frpc.toml",
							SubPath:   "frpc.toml",
						},
					},
				}
				volumes = append(volumes, corev1.Volume{
					Name: volumeName,
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: frpcConfigMap.Name,
							},
						},
					},
				})
				containers = append(containers, frpc)
			}
			for _, policy := range pod.NetworkPolicies {
				_, ok, msg := CreateNetworkPolicy(ctx, pod, policy)
				if !ok {
					resultCh <- result{PodName: pod.Name, OK: false, Msg: msg}
					return
				}
			}
			for _, container := range pod.Containers {
				if container.Image == "" {
					resultCh <- result{PodName: pod.Name, OK: false, Msg: i18n.EmptyContainerImage}
					return
				}
				containers = append(containers, corev1.Container{
					Name:  container.Name,
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
						for _, p := range container.PodPorts {
							tmp = append(tmp, corev1.ContainerPort{
								ContainerPort: p,
							})
						}
						return tmp
					}(),
				})
			}
			ports := make([]int32, 0)
			for _, port := range service.Spec.Ports {
				if !utils.In(port.NodePort, ports) {
					ports = append(ports, port.NodePort)
				}
			}
			p, ok, msg := CreatePod(ctx, pod.Name, containers, volumes, pod.PodIP, dns)
			if !ok {
				resultCh <- result{PodName: pod.Name, OK: false, Msg: msg}
				return
			}
			if !config.Env.K8S.Frpc.On {
				ip = p.Status.HostIP
			}
			log.Logger.Infof("Pod %s is running on %s", pod.Name, ip)
			resultCh <- result{PodName: pod.Name, Ports: ports, IP: ip, OK: true, Msg: msg}
		}(pod)
	}
	wg.Wait()
	close(resultCh)
	targets := make(map[string]map[string]any)
	for res := range resultCh {
		if !res.OK {
			return nil, false, res.Msg
		}
		targets[res.PodName] = map[string]any{
			"ip":    res.IP,
			"ports": res.Ports,
		}
	}
	return targets, true, i18n.Success
}

// StopVictim 需要预加载 model.Pod
func StopVictim(victim model.Victim) (bool, string) {
	log.Logger.Debugf("Stopping Victim for team %d usage %d", victim.TeamID, victim.UsageID)
	type result struct {
		OK  bool
		Msg string
	}
	var wg sync.WaitGroup
	resultCh := make(chan result, len(victim.Pods))
	for _, pod := range victim.Pods {
		wg.Add(1)
		go func(pod model.Pod) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
			defer cancel()
			var err error
			err = CopyFromPod(
				pod.Name, "tcpdump", "/root/traffic.pcap",
				pod.TrafficPath(),
			)
			if err != nil {
				log.Logger.Warningf("Failed to copy %d traffic: %v", victim.TeamID, err)
			}
			if ok, msg := DeleteNetworkPolicy(ctx, pod.NetworkPolicyName); !ok {
				resultCh <- result{OK: false, Msg: msg}
				return
			}
			if ok, msg := DeleteService(ctx, pod.ServiceName); !ok {
				resultCh <- result{OK: false, Msg: msg}
				return
			}
			if ok, msg := DeleteConfigMap(ctx, fmt.Sprintf("%s-cm", pod.Name)); !ok {
				resultCh <- result{OK: false, Msg: msg}
				return
			}
			ok, msg := DeletePod(ctx, pod.Name)
			resultCh <- result{OK: ok, Msg: msg}
		}(pod)
	}
	wg.Wait()
	close(resultCh)
	for res := range resultCh {
		if !res.OK {
			return false, res.Msg
		}
	}
	return true, i18n.Success
}
