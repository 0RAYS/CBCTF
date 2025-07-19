package k8s

import (
	"CBCTF/internel/config"
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"CBCTF/internel/utils"
	"context"
	"fmt"
	kubeovnv1 "github.com/JBNRZ/kubeovn-api/pkg/apis/kubeovn/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"strings"
	"sync"
	"time"
)

func StartVictim(victim model.Victim) (bool, string) {
	log.Logger.Infof("Creating Victim for team %d challenge %d", victim.TeamID, victim.ContestChallengeID)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	labels := map[string]string{
		"contest_challenge_id": fmt.Sprintf("%d", victim.ContestChallengeID),
		"team_id":              fmt.Sprintf("%d", victim.TeamID),
		"user_id":              fmt.Sprintf("%d", victim.UserID),
	}
	//rand.New(rand.NewSource(time.Now().UnixNano()))
	//frps := config.Env.K8S.Frpc.Frps[rand.Intn(len(config.Env.K8S.Frpc.Frps))]
	//policy := model.NetworkPolicy{
	//	To: []*netv1.IPBlock{
	//		{
	//			CIDR: fmt.Sprintf("%s/32", frps.Host),
	//		},
	//	},
	//}
	//for _, p := range victim.NetworkPolicies {
	//	// 当已经存在来源策略时, 需要设frps的ip为白名单; 不存在时, 默认允许
	//	if len(p.From) > 0 {
	//		policy.From = []*netv1.IPBlock{
	//			{
	//				CIDR: fmt.Sprintf("%s/32", frps.Host),
	//			},
	//		}
	//		break
	//	}
	//}
	//victim.NetworkPolicies = append(victim.NetworkPolicies, policy)
	//for _, policy := range victim.NetworkPolicies {
	//	_, ok, msg := CreateNetworkPolicy(ctx, CreateNetworkPolicyOptions{
	//		Name:        fmt.Sprintf("np-%s", utils.RandStr(10)),
	//		Labels:      labels,
	//		MatchLabels: labels,
	//		From:        policy.From,
	//		To:          policy.To,
	//	})
	//	if !ok {
	//		return false, msg
	//	}
	//}
	subnetMap := make(map[string]*model.Subnet)
	netAttchDefMap := make(map[string]*model.NetAttachDef)
	if victim.VPC.Name != "" {
		var policyRoutes []*kubeovnv1.PolicyRoute
		for _, subnet := range victim.VPC.Subnets {
			_, ok, msg := CreateSubnet(ctx, CreateSubnetOptions{
				Name:       subnet.Name,
				Labels:     labels,
				VPC:        victim.VPC.Name,
				CIDR:       subnet.CIDRBlock,
				Gateway:    subnet.Gateway,
				ExcludeIPs: subnet.ExcludeIps,
				Provider:   fmt.Sprintf("%s.%s.ovn", subnet.NetAttachDef.Name, GlobalNamespace),
			})
			if !ok {
				return false, msg
			}
			subnetMap[subnet.DefName] = subnet
			_, ok, msg = CreateNetAttachDef(ctx, CreateNetAttachDefOptions{
				Name:      subnet.NetAttachDef.Name,
				Namespace: GlobalNamespace,
				Labels:    labels,
				Config: fmt.Sprintf(`{
					"cniVersion": "0.3.0",
					"type": "kube-ovn",
					"server_socket": "/run/openvswitch/kube-ovn-daemon.sock",
					"provider": "%s.%s.ovn"
				}`, subnet.NetAttachDef.Name, GlobalNamespace),
			})
			if !ok {
				return false, msg
			}
			netAttchDefMap[subnet.DefName] = subnet.NetAttachDef
			if subnet.NatGateway != nil {
				_, ok, msg = CreateVPCNatGateway(ctx, CreateVPCNatGatewayOptions{
					Name:           subnet.NatGateway.Name,
					Labels:         labels,
					VPC:            victim.VPC.Name,
					Subnet:         subnet.Name,
					LanIP:          subnet.NatGateway.LanIP,
					ExternalSubnet: []string{ExternalSubnetName},
				})
				if !ok {
					return false, msg
				}
				policyRoutes = append(policyRoutes, &kubeovnv1.PolicyRoute{
					Action:    kubeovnv1.PolicyRouteActionReroute,
					Match:     fmt.Sprintf("ip4.src == %s", subnet.CIDRBlock),
					NextHopIP: subnet.NatGateway.LanIP,
					Priority:  1,
				})
				for _, eip := range subnet.NatGateway.EIPs {
					_, ok, msg = CreateEIP(ctx, CreateEIPOptions{
						Name:           eip.Name,
						Labels:         labels,
						NatGw:          subnet.NatGateway.Name,
						ExternalSubnet: ExternalSubnetName,
					})
					if !ok {
						return false, msg
					}
					for _, dnat := range eip.DNats {
						_, ok, msg = CreateDNat(ctx, CreateDNatOptions{
							Name:         dnat.Name,
							Labels:       labels,
							EIP:          eip.Name,
							ExternalPort: dnat.ExternalPort,
							InternalPort: dnat.InternalPort,
							InternalIP:   dnat.InternalIP,
							Protocol:     dnat.Protocol,
						})
						if !ok {
							return false, msg
						}
					}
					for _, snat := range eip.SNats {
						_, ok, msg = CreateSNat(ctx, CreateSNatOptions{
							Name:         snat.Name,
							Labels:       labels,
							EIP:          eip.Name,
							InternalCIDR: subnet.CIDRBlock,
						})
						if !ok {
							return false, msg
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
			return false, msg
		}
	}
	for _, pod := range victim.Pods {
		containers := []corev1.Container{
			{
				Name:    "tcpdump",
				Image:   config.Env.K8S.TCPDumpImage,
				Command: []string{"/bin/sh", "-c", "tcpdump -i any -w /root/traffic.pcap"},
			},
		}
		volumes := make([]corev1.Volume, 0)
		//frpcConfigMap, ok, msg := CreateFrpcConfig(ctx, frps.Host, frps.Port, frps.Token, pod)
		//if !ok {
		//	return false, msg
		//}
		//volumeName := fmt.Sprintf("vol-%s", utils.RandStr(5))
		//frpc := corev1.Container{
		//	Name:  "frpc",
		//	Image: config.Env.K8S.Frpc.Image,
		//	Args:  []string{"-c", "/etc/frp/frpc.toml"},
		//	VolumeMounts: []corev1.VolumeMount{
		//		{
		//			Name:      volumeName,
		//			MountPath: "/etc/frp/frpc.toml",
		//			SubPath:   "frpc.toml",
		//		},
		//	},
		//}
		//volumes = append(volumes, corev1.Volume{
		//	Name: volumeName,
		//	VolumeSource: corev1.VolumeSource{
		//		ConfigMap: &corev1.ConfigMapVolumeSource{
		//			LocalObjectReference: corev1.LocalObjectReference{
		//				Name: frpcConfigMap.Name,
		//			},
		//		},
		//	},
		//})
		//containers = append(containers, frpc)
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
					return false, msg
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
				return false, i18n.SubnetNotFound
			}
			netAttachDef, ok := netAttchDefMap[network.Name]
			if !ok {
				return false, i18n.NetAttNotFound
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
		_, ok, msg := CreatePod(ctx, CreatePodOptions{
			Name:        pod.Name,
			Labels:      labels,
			Annotations: annotations,
			Containers:  containers,
			Volumes:     volumes,
		})
		if !ok {
			return false, msg
		}
	}
	return true, i18n.Success
}

func StopVictim(victim model.Victim) (bool, string) {
	log.Logger.Infof("Stopping Victim for team %d challenge %d", victim.TeamID, victim.ContestChallengeID)
	labels := map[string]string{
		"contest_challenge_id": fmt.Sprintf("%d", victim.ContestChallengeID),
		"team_id":              fmt.Sprintf("%d", victim.TeamID),
		"user_id":              fmt.Sprintf("%d", victim.UserID),
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
		}(pod)
	}
	wg.Wait()
	close(resultCh)
	for res := range resultCh {
		if !res.OK {
			return false, res.Msg
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
	return true, i18n.Success
}
