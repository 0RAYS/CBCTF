package k8s

import (
	"CBCTF/internal/config"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
	"CBCTF/internal/utils"
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"sync"

	kubeovnv1 "github.com/JBNRZ/kubeovn-api/pkg/apis/kubeovn/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// StartVictim model.Victim Preload model.Pod
func StartVictim(ctx context.Context, victim model.Victim) (map[string]model.Exposes, model.RetVal) {
	log.Logger.Infof("Starting Victim for Team %d Challenge %d", victim.TeamID.V, victim.ChallengeID)
	// 添加一个独立tag, 防止 NetworkPolicy 影响 frpc 通信
	labels := map[string]string{
		"victim_id":            strconv.Itoa(int(victim.ID)),
		"user_id":              strconv.Itoa(int(victim.UserID)),
		"team_id":              strconv.Itoa(int(victim.TeamID.V)),
		"contest_id":           strconv.Itoa(int(victim.ContestID.V)),
		"challenge_id":         strconv.Itoa(int(victim.ChallengeID)),
		"contest_challenge_id": strconv.Itoa(int(victim.ContestChallengeID.V)),
		VictimPodTag:           VictimPodTag,
	}
	subnetMap := make(map[string]*model.Subnet)
	netAttchDefMap := make(map[string]*model.NetAttachDef)
	ipExposesMap := make(map[string]model.Exposes)
	ipExposesMapMutex := &sync.Mutex{}
	wg := utils.NewGroup(ctx)
	wg.Go(func() error {
		name := fmt.Sprintf("np-%s", utils.RandStr(20))
		_, ret := CreateNetworkPolicy(ctx, CreateNetworkPolicyOptions{
			Name:        name,
			Labels:      labels,
			MatchLabels: labels,
			From: func() []*netv1.IPBlock {
				tmp := make([]*netv1.IPBlock, 0)
				for _, p := range victim.NetworkPolicies {
					tmp = append(tmp, p.From...)
				}
				return tmp
			}(),
			To: func() []*netv1.IPBlock {
				tmp := make([]*netv1.IPBlock, 0)
				for _, p := range victim.NetworkPolicies {
					tmp = append(tmp, p.To...)
				}
				return tmp
			}(),
		})
		log.Logger.Debugf("Create NetworkPolicy %s: %s", name, ret.Msg)
		if err, ok := ret.Attr["Error"]; ok && !ret.OK {
			return errors.New(err.(string))
		}
		return nil
	})
	if victim.VPC.Name != "" {
		// 首先创建 VPC 资源, 导致多跑一个循环
		var policyRoutes []*kubeovnv1.PolicyRoute
		for _, subnet := range victim.VPC.Subnets {
			if subnet.NatGateway != nil {
				policyRoutes = append(policyRoutes, &kubeovnv1.PolicyRoute{
					Action:    kubeovnv1.PolicyRouteActionReroute,
					Match:     fmt.Sprintf("ip4.src == %s", subnet.CIDRBlock),
					NextHopIP: subnet.NatGateway.LanIP,
					Priority:  1,
				})
			}
		}
		wg.Go(func() error {
			_, ret := CreateVPC(ctx, CreateVPCOptions{
				Name:         victim.VPC.Name,
				Labels:       labels,
				PolicyRoutes: policyRoutes,
			})
			log.Logger.Debugf("Create VPC %s: %s", victim.VPC.Name, ret.Msg)
			if err, ok := ret.Attr["Error"]; ok && !ret.OK {
				return errors.New(err.(string))
			}
			return nil
		})
		for _, subnet := range victim.VPC.Subnets {
			wg.Go(func() error {
				_, ret := CreateNetAttachDef(ctx, CreateNetAttachDefOptions{
					Name:      subnet.NetAttachDef.Name,
					Namespace: globalNamespace,
					Labels:    labels,
					Config: fmt.Sprintf(`{
						"cniVersion": "0.3.0",
						"type": "kube-ovn",
						"server_socket": "/run/openvswitch/kube-ovn-daemon.sock",
						"provider": "%s.%s.ovn"
					}`, subnet.NetAttachDef.Name, globalNamespace),
				})
				log.Logger.Debugf("Create NetAttachDef %s: %s", subnet.NetAttachDef.Name, ret.Msg)
				if err, ok := ret.Attr["Error"]; ok && !ret.OK {
					return errors.New(err.(string))
				}
				return nil
			})
			wg.Go(func() error {
				_, ret := CreateSubnet(ctx, CreateSubnetOptions{
					Name:       subnet.Name,
					Labels:     labels,
					VPC:        victim.VPC.Name,
					CIDR:       subnet.CIDRBlock,
					Gateway:    subnet.Gateway,
					ExcludeIPs: subnet.ExcludeIps,
					Provider:   fmt.Sprintf("%s.%s.ovn", subnet.NetAttachDef.Name, globalNamespace),
				})
				log.Logger.Debugf("Create Subnet %s: %s", subnet.Name, ret.Msg)
				if err, ok := ret.Attr["Error"]; ok && !ret.OK {
					return errors.New(err.(string))
				}
				return nil
			})
			subnetMap[subnet.DefName] = subnet
			netAttchDefMap[subnet.DefName] = subnet.NetAttachDef
			if subnet.NatGateway != nil {
				wg.Go(func() error {
					_, ret := CreateVPCNatGateway(ctx, CreateVPCNatGatewayOptions{
						Name:           subnet.NatGateway.Name,
						Labels:         labels,
						VPC:            victim.VPC.Name,
						Subnet:         subnet.Name,
						LanIP:          subnet.NatGateway.LanIP,
						ExternalSubnet: []string{externalSubnetName},
					})
					log.Logger.Debugf("Create VPCNatGateway %s: %s", subnet.NatGateway.Name, ret.Msg)
					if err, ok := ret.Attr["Error"]; ok && !ret.OK {
						return errors.New(err.(string))
					}
					return nil
				})
				for _, eip := range subnet.NatGateway.EIPs {
					wg.Go(func() error {
						e, ret := CreateEIP(ctx, CreateEIPOptions{
							Name:           eip.Name,
							Labels:         labels,
							NatGw:          subnet.NatGateway.Name,
							ExternalSubnet: externalSubnetName,
						})
						log.Logger.Debugf("Create EIP %s: %s", eip.Name, ret.Msg)
						if err, ok := ret.Attr["Error"]; ok && !ret.OK {
							return errors.New(err.(string))
						}
						// 后续会用到
						eip.IP = e.Spec.V4ip
						for _, dnat := range eip.DNats {
							_, ret = CreateDNat(ctx, CreateDNatOptions{
								Name:         dnat.Name,
								Labels:       labels,
								EIP:          eip.Name,
								ExternalPort: strconv.Itoa(int(dnat.ExternalPort)),
								InternalPort: strconv.Itoa(int(dnat.InternalPort)),
								InternalIP:   dnat.InternalIP,
								Protocol:     dnat.Protocol,
							})
							log.Logger.Debugf("Create DNat %s: %s", dnat.Name, ret.Msg)
							if err, ok := ret.Attr["Error"]; ok && !ret.OK {
								return errors.New(err.(string))
							}
							ipExposesMapMutex.Lock()
							if !slices.ContainsFunc(ipExposesMap[e.Spec.V4ip], func(e model.Expose) bool {
								return dnat.ExternalPort == e.Port && dnat.Protocol == e.Protocol
							}) {
								ipExposesMap[e.Spec.V4ip] = append(ipExposesMap[e.Spec.V4ip], model.Expose{
									Port:     dnat.ExternalPort,
									Protocol: dnat.Protocol,
								})
							}
							ipExposesMapMutex.Unlock()
						}
						for _, snat := range eip.SNats {
							_, ret = CreateSNat(ctx, CreateSNatOptions{
								Name:         snat.Name,
								Labels:       labels,
								EIP:          eip.Name,
								InternalCIDR: subnet.CIDRBlock,
							})
							log.Logger.Debugf("Create SNat %s: %s", snat.Name, ret.Msg)
							if err, ok := ret.Attr["Error"]; ok && !ret.OK {
								return errors.New(err.(string))
							}
						}
						return nil
					})
				}
			}
		}
	}
	for _, pod := range victim.Pods {
		wg.Go(func() error {
			nfsName := fmt.Sprintf("vol-%s", utils.RandStr(20))
			containers := []corev1.Container{
				{
					Name:    "tcpdump",
					Image:   config.Env.K8S.TCPDumpImage,
					Command: []string{"/bin/sh", "-c", fmt.Sprintf("tcpdump -i any -w /root/mnt/pod-%d.pcap", pod.ID)},
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      nfsName,
							MountPath: "/root/mnt",
							SubPath: strings.TrimPrefix(
								strings.TrimPrefix(victim.TrafficBasePath(), config.Env.Path), "/",
							),
						},
					},
				},
			}
			volumes := []corev1.Volume{
				{
					Name: nfsName,
					VolumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: nfsVolumeName,
						},
					},
				},
			}
			for _, container := range pod.Containers {
				volumeMounts := make([]corev1.VolumeMount, 0)
				for path, volumeFlag := range container.VolumeFlags {
					filename := strings.Split(path, "/")[len(strings.Split(path, "/"))-1]
					flagConfigMap, ret := CreateConfigMap(ctx, CreateConfigMapOptions{
						Name:   fmt.Sprintf("cm-%s", utils.RandStr(20)),
						Labels: labels,
						Data:   map[string]string{filename: volumeFlag},
					})
					if err, ok := ret.Attr["Error"]; ok && !ret.OK {
						return errors.New(err.(string))
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
				for name, envFlag := range container.EnvFlags {
					envs = append(envs, corev1.EnvVar{
						Name:  name,
						Value: envFlag,
					})
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
					limit["memory"] = resource.MustParse(strconv.Itoa(int(container.Memory)))
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
					return fmt.Errorf("subnet %s not found", network.Name)
				}
				netAttachDef, ok := netAttchDefMap[network.Name]
				if !ok {
					return fmt.Errorf("netAttachDef %s not found", network.Name)
				}
				if i == 0 {
					annotations["ovn.kubernetes.io/logical_switch"] = subnet.Name
					annotations["ovn.kubernetes.io/ip_address"] = network.IP
				} else {
					annotations["k8s.v1.cni.cncf.io/networks"] += fmt.Sprintf(",%s/%s", globalNamespace, netAttachDef.Name)
					annotations["k8s.v1.cni.cncf.io/networks"] = strings.Trim(annotations["k8s.v1.cni.cncf.io/networks"], ",")
					annotations[fmt.Sprintf("%s.%s.ovn.kubernetes.io/logical_switch", netAttachDef.Name, globalNamespace)] = subnet.Name
					annotations[fmt.Sprintf("%s.%s.ovn.kubernetes.io/ip_address", netAttachDef.Name, globalNamespace)] = network.IP
				}
			}
			pOptions := CreatePodOptions{
				Name:        pod.Name,
				Labels:      labels,
				Annotations: annotations,
				Containers:  containers,
				Volumes:     volumes,
			}
			// 容忍不支持VPC网络的节点
			if len(annotations) == 0 {
				pOptions.Tolerations = map[string]string{VPCNetworkTolerationKey: VPCNetworkTolerationVal}
			}
			p, ret := CreatePod(ctx, pOptions)
			if err, ok := ret.Attr["Error"]; ok && !ret.OK {
				return errors.New(err.(string))
			}
			if len(annotations) == 0 {
				service, ret := CreateService(ctx, CreateServiceOptions{
					Name:     fmt.Sprintf("svc-%s", utils.RandStr(20)),
					Ports:    pod.PodPorts,
					Labels:   labels,
					Selector: labels,
				})
				if err, ok := ret.Attr["Error"]; ok && !ret.OK {
					return errors.New(err.(string))
				}
				ipExposesMapMutex.Lock()
				for _, port := range service.Spec.Ports {
					if !slices.ContainsFunc(ipExposesMap[p.Status.HostIP], func(e model.Expose) bool {
						return port.NodePort == e.Port && strings.ToLower(e.Protocol) == strings.ToLower(string(port.Protocol))
					}) {
						ipExposesMap[p.Status.HostIP] = append(ipExposesMap[p.Status.HostIP], model.Expose{
							Port:     port.NodePort,
							Protocol: string(port.Protocol),
						})
					}
				}
				ipExposesMapMutex.Unlock()
			}
			log.Logger.Debugf("Create Pod %s: %s", pod.Name, ret.Msg)
			return nil
		})
	}
	if err := wg.Wait(); err != nil {
		return nil, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}}
	}
	return ipExposesMap, model.SuccessRetVal()
}

func StopVictim(ctx context.Context, victim model.Victim) model.RetVal {
	log.Logger.Infof("Stopping Victim for Team %d Challenge %d", victim.TeamID.V, victim.ChallengeID)
	// 不添加独立 tag, 删除时直接删除所有相关资源
	labels := map[string]string{
		"victim_id":            strconv.Itoa(int(victim.ID)),
		"user_id":              strconv.Itoa(int(victim.UserID)),
		"team_id":              strconv.Itoa(int(victim.TeamID.V)),
		"contest_id":           strconv.Itoa(int(victim.ContestID.V)),
		"challenge_id":         strconv.Itoa(int(victim.ChallengeID)),
		"contest_challenge_id": strconv.Itoa(int(victim.ContestChallengeID.V)),
	}
	for _, endpoint := range victim.ExposedEndpoints {
		redis.UnlockFrpsPort(endpoint.IP, endpoint.Port, endpoint.Protocol)
	}
	if ret := DeleteDNatList(ctx, labels); !ret.OK {
		return ret
	}
	if ret := DeleteSNatList(ctx, labels); !ret.OK {
		return ret
	}
	if ret := DeleteEIPList(ctx, labels); !ret.OK {
		return ret
	}
	if ret := DeleteSubnetList(ctx, labels); !ret.OK {
		return ret
	}
	if ret := DeleteNetAttachDefList(ctx, globalNamespace, labels); !ret.OK {
		return ret
	}
	if ret := DeleteVPCNatGatewayList(ctx, labels); !ret.OK {
		return ret
	}
	if ret := DeleteVPCList(ctx, labels); !ret.OK {
		return ret
	}
	if ret := DeleteConfigMapList(ctx, labels); !ret.OK {
		return ret
	}
	if ret := DeleteNetworkPolicyList(ctx, labels); !ret.OK {
		return ret
	}
	if ret := DeleteEndpointList(ctx, labels); !ret.OK {
		return ret
	}
	if ret := DeleteServiceList(ctx, labels); !ret.OK {
		return ret
	}
	if ret := DeletePodList(ctx, labels); !ret.OK {
		return ret
	}
	for _, subnet := range victim.VPC.Subnets {
		if ret := DeleteIPList(ctx, map[string]string{"ovn.kubernetes.io/subnet": subnet.Name}); !ret.OK {
			return ret
		}
	}
	return model.SuccessRetVal()
}
