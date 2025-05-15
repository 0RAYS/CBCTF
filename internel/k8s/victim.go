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
	"strconv"
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
			pod.NetworkPolicies = append(pod.NetworkPolicies, model.DefaultNetworkPolicy)
			for _, policy := range pod.NetworkPolicies {
				_, ok, msg := CreateNetworkPolicy(ctx, pod, policy)
				if !ok {
					resultCh <- result{PodName: pod.Name, OK: false, Msg: msg}
					return
				}
			}
			containers := []corev1.Container{
				{
					Name:    "tcpdump",
					Image:   config.Env.K8S.TCPDumpImage,
					Command: []string{"/bin/sh", "-c", "tcpdump -i any -w /root/traffic.pcap"},
				},
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
			rand.New(rand.NewSource(time.Now().UnixNano()))
			frps := config.Env.K8S.Frpc.Frps[rand.Intn(len(config.Env.K8S.Frpc.Frps))]
			var ip string
			ports := make([]int32, 0)
			// TODO 优化端口映射方案, 多个端口使用单个 frpc 容器
			for _, port := range service.Spec.Ports {
				if !utils.In(port.NodePort, ports) {
					ports = append(ports, port.NodePort)
				}
				if config.Env.K8S.Frpc.On {
					ip = frps.Host
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
								Value: fmt.Sprintf("%s-%d", pod.Name, port.Port),
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
					log.Logger.Infof("Frpc started: %s:%d -> %s:%s", frps.Host, port.NodePort, pod.Name, port.TargetPort.StrVal)
				}
			}
			p, ok, msg := CreatePod(ctx, pod.Name, containers, pod.PodIP, dns)
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
