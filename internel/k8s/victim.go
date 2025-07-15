package k8s

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"context"
	"fmt"
	kubeovnv1 "github.com/JBNRZ/kubeovn-api/pkg/apis/kubeovn/v1"
	netattv1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
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
	if victim.VPC != "" {
		subnetMap := make(map[string]model.Subnet)
		for _, subnet := range victim.Subnets {
			subnetMap[subnet.Name] = subnet
		}
		subnetGatewayMap := make(map[string]model.Gateway)
		var policyRoutes []*kubeovnv1.PolicyRoute
		for _, gateway := range victim.Gateways {
			subnetGatewayMap[gateway.Subnet] = gateway
			policyRoutes = append(policyRoutes, &kubeovnv1.PolicyRoute{
				Action:    kubeovnv1.PolicyRouteActionReroute,
				Match:     fmt.Sprintf("ip.src == %s", subnetMap[gateway.Subnet].CIDR),
				NextHopIP: gateway.LanIP,
				Priority:  1,
			})
		}
		vpc, ok, msg := CreateVPC(ctx, CreateVPCOptions{
			Name:         victim.VPC,
			Labels:       labels,
			PolicyRoutes: policyRoutes,
		})
		if !ok {
			return false, msg
		}
		netAttachDefMap := make(map[string]*netattv1.NetworkAttachmentDefinition)
		for subnetName, netAttachDefName := range victim.NetAttachDefs {
			nad, ok, msg := CreateNetAttachDef(ctx, CreateNetAttachDefOptions{
				Name:      netAttachDefName,
				Namespace: GlobalNamespace,
				Labels:    labels,
				Config: fmt.Sprintf(`{
				  "cniVersion": "0.3.0",
				  "type": "kube-ovn",
				  "server_socket": "/run/openvswitch/kube-ovn-daemon.sock",
				  "provider": "%s.%s.ovn"
				}`, subnetName, GlobalNamespace),
			})
			if !ok {
				return false, msg
			}
			netAttachDefMap[netAttachDefName] = nad
		}
		subnets := make(map[string]*kubeovnv1.Subnet)
		for _, net := range subnetMap {
			subnet, ok, msg := CreateSubnet(ctx, CreateSubnetOptions{
				Name:       net.Name,
				Labels:     labels,
				VPC:        vpc.Name,
				CIDR:       net.CIDR,
				Gateway:    net.Gateway,
				ExcludeIPs: []string{net.Gateway, subnetGatewayMap[net.Name].LanIP},
			})
			if !ok {
				return false, msg
			}
			subnets[subnet.Name] = subnet
		}
		gateways := make(map[string]*kubeovnv1.VpcNatGateway)
		for _, gateway := range victim.Gateways {
			natGateway, ok, msg := CreateVPCNatGateway(ctx, CreateVPCNatGatewayOptions{
				Name:           gateway.Name,
				Labels:         labels,
				VPC:            gateway.VPC,
				Subnet:         gateway.Subnet,
				LanIP:          gateway.LanIP,
				ExternalSubnet: []string{ExternalSubnetName},
			})
			if !ok {
				return false, msg
			}
			gateways[natGateway.Name] = natGateway
		}
		eips := make(map[string]*kubeovnv1.IptablesEIP)
		for _, e := range victim.EIPs {
			eip, ok, msg := CreateEIP(ctx, CreateEIPOptions{
				Name:           e.Name,
				Labels:         labels,
				NatGw:          e.Gateway,
				ExternalSubnet: ExternalSubnetName,
			})
			if !ok {
				return false, msg
			}
			eips[eip.Name] = eip
		}
		dnats := make(map[string]*kubeovnv1.IptablesDnatRule)
		for _, d := range victim.DNats {
			dnat, ok, msg := CreateDNat(ctx, CreateDNatOptions{
				Name:         d.Name,
				Labels:       labels,
				ExternalPort: fmt.Sprintf("%d", d.ExternalPort),
				InternalPort: d.InternalIP,
				InternalIP:   d.InternalIP,
				Protocol:     d.Protocol,
			})
			if !ok {
				return false, msg
			}
			dnats[dnat.Name] = dnat
		}
		snats := make(map[string]*kubeovnv1.IptablesSnatRule)
		for _, s := range victim.SNats {
			snat, ok, msg := CreateSNat(ctx, CreateSNatOptions{
				Name:         s.Name,
				Labels:       labels,
				EIP:          s.EIP,
				InternalCIDR: s.InternalCIDR,
			})
			if !ok {
				return false, msg
			}
			snats[snat.Name] = snat
		}
	}
	//type result struct {
	//	PodName string
	//	IP      string
	//	Ports   []int32
	//	OK      bool
	//	Msg     string
	//}
	//var wg sync.WaitGroup
	//resultCh := make(chan result, len(victim.Pods))
	//for _, pod := range victim.Pods {
	//	wg.Add(1)
	//	go func(pod model.Pod) {
	//		defer wg.Done()
	//		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	//		defer cancel()
	//		service, ok, msg := CreateService(ctx, CreateServiceOptions{
	//			Name:     fmt.Sprintf("svc-%s", strings.ToLower(utils.RandStr(10))),
	//			Ports:    pod.PodPorts,
	//			Labels:   map[string]string{VictimPodTag: pod.Name},
	//			Selector: map[string]string{VictimPodTag: pod.Name},
	//		})
	//		if !ok {
	//			resultCh <- result{PodName: pod.Name, OK: false, Msg: msg}
	//			return
	//		}
	//		containers := []corev1.Container{
	//			{
	//				Name:    "tcpdump",
	//				Image:   config.Env.K8S.TCPDumpImage,
	//				Command: []string{"/bin/sh", "-c", "tcpdump -i any -w /root/traffic.pcap"},
	//			},
	//		}
	//		volumes := make([]corev1.Volume, 0)
	//		rand.New(rand.NewSource(time.Now().UnixNano()))
	//		frps := config.Env.K8S.Frpc.Frps[rand.Intn(len(config.Env.K8S.Frpc.Frps))]
	//		if len(pod.PodPorts) > 0 && config.Env.K8S.Frpc.On {
	//			policy := model.NetworkPolicy{
	//				To: []model.Target{
	//					{
	//						CIDR: fmt.Sprintf("%s/32", frps.Host),
	//					},
	//				},
	//			}
	//			for _, p := range pod.NetworkPolicies {
	//				// 当已经存在来源策略时, 需要设frps的ip为白名单; 不存在时, 默认允许
	//				if len(p.From) > 0 {
	//					policy.From = []model.Target{
	//						{
	//							CIDR: fmt.Sprintf("%s/32", frps.Host),
	//						},
	//					}
	//					break
	//				}
	//			}
	//			pod.NetworkPolicies = append(pod.NetworkPolicies, policy)
	//			frpcConfigMap, ok, msg := CreateFrpcConfig(ctx, frps.Host, frps.Port, frps.Token, pod, service)
	//			if !ok {
	//				resultCh <- result{PodName: pod.Name, OK: false, Msg: msg}
	//				return
	//			}
	//			volumeName := fmt.Sprintf("vol-%s", strings.ToLower(utils.RandStr(5)))
	//			frpc := corev1.Container{
	//				Name:  "frpc",
	//				Image: config.Env.K8S.Frpc.Image,
	//				Args:  []string{"-c", "/etc/frp/frpc.toml"},
	//				VolumeMounts: []corev1.VolumeMount{
	//					{
	//						Name:      volumeName,
	//						MountPath: "/etc/frp/frpc.toml",
	//						SubPath:   "frpc.toml",
	//					},
	//				},
	//			}
	//			volumes = append(volumes, corev1.Volume{
	//				Name: volumeName,
	//				VolumeSource: corev1.VolumeSource{
	//					ConfigMap: &corev1.ConfigMapVolumeSource{
	//						LocalObjectReference: corev1.LocalObjectReference{
	//							Name: frpcConfigMap.Name,
	//						},
	//					},
	//				},
	//			})
	//			containers = append(containers, frpc)
	//		}
	//		for _, policy := range pod.NetworkPolicies {
	//			_, ok, msg = CreateNetworkPolicy(ctx, CreateNetworkPolicyOptions{
	//				Name:        fmt.Sprintf("np-%s", strings.ToLower(utils.RandStr(10))),
	//				Labels:      map[string]string{VictimPodTag: pod.Name},
	//				MatchLabels: map[string]string{VictimPodTag: pod.Name},
	//				From: func() []*netv1.IPBlock {
	//					from := make([]*netv1.IPBlock, 0)
	//					for _, v := range policy.From {
	//						from = append(from, &netv1.IPBlock{
	//							CIDR:   v.CIDR,
	//							Except: v.Except,
	//						})
	//					}
	//					return from
	//				}(),
	//				To: func() []*netv1.IPBlock {
	//					to := make([]*netv1.IPBlock, 0)
	//					for _, v := range policy.To {
	//						to = append(to, &netv1.IPBlock{
	//							CIDR:   v.CIDR,
	//							Except: v.Except,
	//						})
	//					}
	//					return to
	//				}(),
	//			})
	//			if !ok {
	//				resultCh <- result{PodName: pod.Name, OK: false, Msg: msg}
	//				return
	//			}
	//		}
	//		for _, container := range pod.Containers {
	//			if container.Image == "" {
	//				resultCh <- result{PodName: pod.Name, OK: false, Msg: i18n.InvalidDockerImage}
	//				return
	//			}
	//			volumeMounts := make([]corev1.VolumeMount, 0)
	//			for path, volumeFlag := range container.VolumeFlags {
	//				filename := strings.Split(path, "/")[len(strings.Split(path, "/"))-1]
	//				flagConfigMap, ok, msg := CreateConfigMap(ctx, CreateConfigMapOptions{
	//					Name:   fmt.Sprintf("cm-%s", strings.ToLower(utils.RandStr(10))),
	//					Labels: map[string]string{VictimPodTag: pod.Name},
	//					Data:   map[string]string{filename: volumeFlag},
	//				})
	//				if !ok {
	//					resultCh <- result{PodName: pod.Name, OK: false, Msg: msg}
	//					return
	//				}
	//				volumeName := fmt.Sprintf("vol-%s", strings.ToLower(utils.RandStr(5)))
	//				volumeMounts = append(volumeMounts, corev1.VolumeMount{
	//					Name:      volumeName,
	//					MountPath: path,
	//					SubPath:   filename,
	//				})
	//				volumes = append(volumes, corev1.Volume{
	//					Name: volumeName,
	//					VolumeSource: corev1.VolumeSource{
	//						ConfigMap: &corev1.ConfigMapVolumeSource{
	//							LocalObjectReference: corev1.LocalObjectReference{
	//								Name: flagConfigMap.Name,
	//							},
	//						},
	//					},
	//				})
	//			}
	//			envs := make([]corev1.EnvVar, 0)
	//			for key, value := range container.Environment {
	//				envs = append(envs, corev1.EnvVar{
	//					Name:  key,
	//					Value: value,
	//				})
	//			}
	//			if len(container.EnvFlags) == 1 {
	//				envs = append(envs, corev1.EnvVar{
	//					Name:  "FLAG",
	//					Value: container.EnvFlags[0],
	//				})
	//			} else {
	//				for i, envFlag := range container.EnvFlags {
	//					envs = append(envs, corev1.EnvVar{
	//						Name:  fmt.Sprintf("FLAG%d", i+1),
	//						Value: envFlag,
	//					})
	//				}
	//			}
	//			ports := make([]corev1.ContainerPort, 0)
	//			for _, p := range container.Exposes {
	//				port, _ := strconv.ParseInt(p, 10, 32)
	//				ports = append(ports, corev1.ContainerPort{
	//					ContainerPort: int32(port),
	//				})
	//			}
	//			tmp := corev1.Container{
	//				Name:         container.Name,
	//				Image:        container.Image,
	//				Env:          envs,
	//				Ports:        ports,
	//				VolumeMounts: volumeMounts,
	//			}
	//			if len(container.Command) > 0 {
	//				tmp.Command = container.Command
	//			}
	//			if container.WorkingDir != "" {
	//				tmp.WorkingDir = container.WorkingDir
	//			}
	//			containers = append(containers, tmp)
	//		}
	//		p, ok, msg := CreatePod(ctx, CreatePodOptions{
	//			Name: pod.Name,
	//			Labels: map[string]string{
	//				VictimPodTag:           pod.Name,
	//				"team_id":              strconv.Itoa(int(victim.TeamID)),
	//				"contest_challenge_id": strconv.Itoa(int(victim.ContestChallengeID)),
	//				"user_id":              strconv.Itoa(int(victim.UserID)),
	//			},
	//			Containers: containers,
	//			Volumes:    volumes,
	//		})
	//		if !ok {
	//			resultCh <- result{PodName: pod.Name, OK: false, Msg: msg}
	//			return
	//		}
	//		var ip string
	//		if config.Env.K8S.Frpc.On {
	//			ip = frps.Host
	//		} else {
	//			ip = p.Status.HostIP
	//		}
	//		log.Logger.Infof("Pod %s is running on %s", pod.Name, ip)
	//		ports := make([]int32, 0)
	//		for _, port := range service.Spec.Ports {
	//			if !slices.Contains(ports, port.NodePort) {
	//				ports = append(ports, port.NodePort)
	//			}
	//		}
	//		resultCh <- result{PodName: pod.Name, Ports: ports, IP: ip, OK: true, Msg: msg}
	//	}(pod)
	//}
	//wg.Wait()
	//close(resultCh)
	//targets := make(map[string]map[string]any)
	//for res := range resultCh {
	//	if !res.OK {
	//		return nil, false, res.Msg
	//	}
	//	targets[res.PodName] = map[string]any{
	//		"ip":    res.IP,
	//		"ports": res.Ports,
	//	}
	//}
	//return targets, true, i18n.Success
	return true, i18n.Success
}

func StopVictim(victim model.Victim) (bool, string) {
	log.Logger.Infof("Stopping Victim for team %d challenge %d", victim.TeamID, victim.ContestChallengeID)
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
			if ok, msg := DeleteNetworkPolicyListByPodName(ctx, VictimPodTag, pod.Name); !ok {
				resultCh <- result{OK: false, Msg: msg}
				return
			}
			if ok, msg := DeleteServiceListByPodName(ctx, VictimPodTag, pod.Name); !ok {
				resultCh <- result{OK: false, Msg: msg}
				return
			}
			if ok, msg := DeleteConfigMapListByPodName(ctx, VictimPodTag, pod.Name); !ok {
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
