package k8s

import (
	"CBCTF/internel/config"
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"CBCTF/internel/redis"
	"CBCTF/internel/utils"
	"context"
	"fmt"
	kubeovnv1 "github.com/JBNRZ/kubeovn-api/pkg/apis/kubeovn/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

func StartVictim(victim model.Victim) (map[string]model.Exposes, bool, string) {
	log.Logger.Infof("Creating Victim for team %d challenge %d", victim.TeamID, victim.ContestChallengeID)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	// 添加一个独立tag, 防止 NetworkPolicy 影响 frpc 通信
	labels := map[string]string{
		"user_id":              fmt.Sprintf("%d", victim.UserID),
		"team_id":              fmt.Sprintf("%d", victim.TeamID),
		"contest_challenge_id": fmt.Sprintf("%d", victim.ContestChallengeID),
		VictimPodTag:           fmt.Sprintf("victim-%s", utils.RandStr(20)),
	}
	subnetMap := make(map[string]*model.Subnet)
	netAttchDefMap := make(map[string]*model.NetAttachDef)
	ipExposesMap := make(map[string]model.Exposes)
	ipExposesMapMutex := &sync.Mutex{}
	if victim.VPC.Name != "" {
		var policyRoutes []*kubeovnv1.PolicyRoute
		for _, subnet := range victim.VPC.Subnets {
			if _, ok, msg := CreateSubnet(ctx, CreateSubnetOptions{
				Name:       subnet.Name,
				Labels:     labels,
				VPC:        victim.VPC.Name,
				CIDR:       subnet.CIDRBlock,
				Gateway:    subnet.Gateway,
				ExcludeIPs: subnet.ExcludeIps,
				Provider:   fmt.Sprintf("%s.%s.ovn", subnet.NetAttachDef.Name, GlobalNamespace),
			}); !ok {
				return ipExposesMap, false, msg
			}
			subnetMap[subnet.DefName] = subnet
			if _, ok, msg := CreateNetAttachDef(ctx, CreateNetAttachDefOptions{
				Name:      subnet.NetAttachDef.Name,
				Namespace: GlobalNamespace,
				Labels:    labels,
				Config: fmt.Sprintf(`{
					"cniVersion": "0.3.0",
					"type": "kube-ovn",
					"server_socket": "/run/openvswitch/kube-ovn-daemon.sock",
					"provider": "%s.%s.ovn"
				}`, subnet.NetAttachDef.Name, GlobalNamespace),
			}); !ok {
				return ipExposesMap, false, msg
			}
			netAttchDefMap[subnet.DefName] = subnet.NetAttachDef
			if subnet.NatGateway != nil {
				if _, ok, msg := CreateVPCNatGateway(ctx, CreateVPCNatGatewayOptions{
					Name:           subnet.NatGateway.Name,
					Labels:         labels,
					VPC:            victim.VPC.Name,
					Subnet:         subnet.Name,
					LanIP:          subnet.NatGateway.LanIP,
					ExternalSubnet: []string{ExternalSubnetName},
				}); !ok {
					return ipExposesMap, false, msg
				}
				policyRoutes = append(policyRoutes, &kubeovnv1.PolicyRoute{
					Action:    kubeovnv1.PolicyRouteActionReroute,
					Match:     fmt.Sprintf("ip4.src == %s", subnet.CIDRBlock),
					NextHopIP: subnet.NatGateway.LanIP,
					Priority:  1,
				})
				for _, eip := range subnet.NatGateway.EIPs {
					ip, ok, msg := CreateEIP(ctx, CreateEIPOptions{
						Name:           eip.Name,
						Labels:         labels,
						NatGw:          subnet.NatGateway.Name,
						ExternalSubnet: ExternalSubnetName,
					})
					if !ok {
						return ipExposesMap, false, msg
					}
					for _, dnat := range eip.DNats {
						if _, ok, msg = CreateDNat(ctx, CreateDNatOptions{
							Name:         dnat.Name,
							Labels:       labels,
							EIP:          eip.Name,
							ExternalPort: dnat.ExternalPort,
							InternalPort: dnat.InternalPort,
							InternalIP:   dnat.InternalIP,
							Protocol:     dnat.Protocol,
						}); !ok {
							return ipExposesMap, false, msg
						}
						port, err := strconv.ParseInt(dnat.ExternalPort, 10, 64)
						if err != nil {
							return ipExposesMap, false, msg
						}
						if !slices.ContainsFunc(ipExposesMap[ip.Spec.V4ip], func(e model.Expose) bool {
							return int32(port) == e.Port && dnat.Protocol == e.Protocol
						}) {
							ipExposesMap[ip.Spec.V4ip] = append(ipExposesMap[ip.Spec.V4ip], model.Expose{
								Port:     int32(port),
								Protocol: dnat.Protocol,
							})
						}
					}
					for _, snat := range eip.SNats {
						if _, ok, msg = CreateSNat(ctx, CreateSNatOptions{
							Name:         snat.Name,
							Labels:       labels,
							EIP:          eip.Name,
							InternalCIDR: subnet.CIDRBlock,
						}); !ok {
							return ipExposesMap, false, msg
						}
					}
				}
			}
		}
		_, ok, msg := CreateVPC(ctx, CreateVPCOptions{
			Name:         victim.VPC.Name,
			Labels:       labels,
			PolicyRoutes: policyRoutes,
		})
		if !ok {
			return ipExposesMap, false, msg
		}
	}
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
			containers := []corev1.Container{
				{
					Name:    "tcpdump",
					Image:   config.Env.K8S.TCPDumpImage,
					Command: []string{"/bin/sh", "-c", "tcpdump -i any -w /root/traffic.pcap"},
				},
			}
			volumes := make([]corev1.Volume, 0)
			for _, container := range pod.Containers {
				volumeMounts := make([]corev1.VolumeMount, 0)
				for path, volumeFlag := range container.VolumeFlags {
					filename := strings.Split(path, "/")[len(strings.Split(path, "/"))-1]
					flagConfigMap, ok, msg := CreateConfigMap(ctx, CreateConfigMapOptions{
						Name:   fmt.Sprintf("cm-%s", utils.RandStr(20)),
						Labels: labels,
						Data:   map[string]string{filename: volumeFlag},
					})
					if !ok {
						resultCh <- result{OK: false, Msg: msg}
						return
					}
					volumeName := fmt.Sprintf("vol-%s", utils.RandStr(10))
					volumeMounts = append(volumeMounts, corev1.VolumeMount{
						Name:      volumeName,
						MountPath: path,
						SubPath:   filename,
					})
					volumes = append(volumes, corev1.Volume{
						Name: volumeName,
						VolumeSource: corev1.VolumeSource{
							ConfigMap: &corev1.ConfigMapVolumeSource{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: flagConfigMap.Name,
								},
							},
						},
					})
				}
				envs := make([]corev1.EnvVar, 0)
				for key, value := range container.Environment {
					envs = append(envs, corev1.EnvVar{
						Name:  key,
						Value: value,
					})
				}
				if len(container.EnvFlags) == 1 {
					envs = append(envs, corev1.EnvVar{
						Name:  "FLAG",
						Value: container.EnvFlags[0],
					})
				} else {
					for i, envFlag := range container.EnvFlags {
						envs = append(envs, corev1.EnvVar{
							Name:  fmt.Sprintf("FLAG%d", i+1),
							Value: envFlag,
						})
					}
				}
				ports := make([]corev1.ContainerPort, 0)
				for _, p := range container.Exposes {
					ports = append(ports, corev1.ContainerPort{
						ContainerPort: p.Port,
					})
				}
				limit := make(corev1.ResourceList)
				if container.CPU > 0 {
					limit["cpu"] = resource.MustParse(fmt.Sprintf("%dm", int(container.CPU*1000)))
				}
				if container.Memory > 0 {
					limit["memory"] = resource.MustParse(fmt.Sprintf("%d", container.Memory))
				}
				tmp := corev1.Container{
					Name:         container.Name,
					Image:        container.Image,
					Env:          envs,
					Ports:        ports,
					VolumeMounts: volumeMounts,
					Resources: corev1.ResourceRequirements{
						Limits: limit,
					},
				}
				if len(container.Command) > 0 {
					tmp.Command = container.Command
				}
				if container.WorkingDir != "" {
					tmp.WorkingDir = container.WorkingDir
				}
				containers = append(containers, tmp)
			}
			annotations := make(map[string]string)
			for i, network := range pod.Networks {
				subnet, ok := subnetMap[network.Name]
				if !ok {
					resultCh <- result{OK: false, Msg: i18n.SubnetNotFound}
					return
				}
				netAttachDef, ok := netAttchDefMap[network.Name]
				if !ok {
					resultCh <- result{OK: false, Msg: i18n.NetAttNotFound}
					return
				}
				if i == 0 {
					annotations["ovn.kubernetes.io/logical_switch"] = subnet.Name
					annotations["ovn.kubernetes.io/ip_address"] = network.IP
				} else {
					annotations["k8s.v1.cni.cncf.io/networks"] += fmt.Sprintf(",%s/%s", GlobalNamespace, netAttachDef.Name)
					annotations["k8s.v1.cni.cncf.io/networks"] = strings.Trim(annotations["k8s.v1.cni.cncf.io/networks"], ",")
					annotations[fmt.Sprintf("%s.%s.ovn.kubernetes.io/logical_switch", netAttachDef.Name, GlobalNamespace)] = subnet.Name
					annotations[fmt.Sprintf("%s.%s.ovn.kubernetes.io/ip_address", netAttachDef.Name, GlobalNamespace)] = network.IP
				}
			}
			p, ok, msg := CreatePod(ctx, CreatePodOptions{
				Name:        pod.Name,
				Labels:      labels,
				Annotations: annotations,
				Containers:  containers,
				Volumes:     volumes,
			})
			if !ok {
				resultCh <- result{OK: false, Msg: msg}
				return
			}
			if len(annotations) == 0 {
				ipExposesMapMutex.Lock()
				for _, port := range pod.PodPorts {
					if !slices.ContainsFunc(ipExposesMap[p.Status.PodIP], func(e model.Expose) bool {
						return port.Port == e.Port && e.Protocol == port.Protocol
					}) {
						ipExposesMap[p.Status.PodIP] = append(ipExposesMap[p.Status.PodIP], port)
					}
				}
				ipExposesMapMutex.Unlock()
			}
			resultCh <- result{OK: true, Msg: msg}
		}(pod)
	}
	wg.Wait()
	close(resultCh)
	for res := range resultCh {
		if !res.OK {
			return ipExposesMap, false, res.Msg
		}
	}
	return ipExposesMap, true, i18n.Success
}

func StopVictim(victim model.Victim) (bool, string) {
	log.Logger.Infof("Stopping Victim for team %d challenge %d", victim.TeamID, victim.ContestChallengeID)
	// 不添加独立 tag, 删除时直接删除所有相关资源
	labels := map[string]string{
		"user_id":              fmt.Sprintf("%d", victim.UserID),
		"team_id":              fmt.Sprintf("%d", victim.TeamID),
		"contest_challenge_id": fmt.Sprintf("%d", victim.ContestChallengeID),
	}
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
			if err := CopyFromPod(pod.Name, "tcpdump", "/root/traffic.pcap", pod.TrafficPath()); err != nil {
				log.Logger.Warningf("Failed to copy %d traffic: %v", victim.TeamID, err)
			}
			resultCh <- result{OK: true, Msg: i18n.Success}
		}(pod)
	}
	wg.Wait()
	close(resultCh)
	for res := range resultCh {
		if !res.OK {
			return false, res.Msg
		}
	}
	for _, endpoint := range victim.Endpoints {
		if err := redis.UnlockFrpsPort(endpoint.IP, endpoint.Port, endpoint.Protocol); err != nil {
			log.Logger.Warningf("Failed to unlock frps port: %v", err)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	if ok, msg := DeleteDNatList(ctx, labels); !ok {
		return false, msg
	}
	if ok, msg := DeleteSNatList(ctx, labels); !ok {
		return false, msg
	}
	if ok, msg := DeleteEIPByLabels(ctx, labels); !ok {
		return false, msg
	}
	if ok, msg := DeleteVPCNatGatewayList(ctx, labels); !ok {
		return false, msg
	}
	if ok, msg := DeletePodList(ctx, labels); !ok {
		return false, msg
	}
	if ok, msg := DeleteSubnetList(ctx, labels); !ok {
		return false, msg
	}
	if ok, msg := DeleteVPCList(ctx, labels); !ok {
		return false, msg
	}
	if ok, msg := DeleteNetAttachDefList(ctx, GlobalNamespace, labels); !ok {
		return false, msg
	}
	if ok, msg := DeleteConfigMapList(ctx, labels); !ok {
		return false, msg
	}
	if ok, msg := DeleteNetworkPolicyList(ctx, labels); !ok {
		return false, msg
	}
	if ok, msg := DeleteServiceList(ctx, labels); !ok {
		return false, msg
	}
	return true, i18n.Success
}
