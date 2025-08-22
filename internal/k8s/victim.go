package k8s

import (
	"CBCTF/internal/config"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
	"CBCTF/internal/utils"
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	kubeovnv1 "github.com/JBNRZ/kubeovn-api/pkg/apis/kubeovn/v1"
	netattv1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

type Result struct {
	OK  bool
	MSG string
}

type CreateNADResult struct {
	NAD *netattv1.NetworkAttachmentDefinition
	OK  bool
	MSG string
}

type CreateSubnetResult struct {
	Subnet *kubeovnv1.Subnet
	OK     bool
	MSG    string
}

type CreateVPCNatGWResult struct {
	VPCNatGW *kubeovnv1.VpcNatGateway
	OK       bool
	MSG      string
}

type CreateEIPResult struct {
	EIP *kubeovnv1.IptablesEIP
	OK  bool
	MSG string
}

type CreateDNatResult struct {
	DNat *kubeovnv1.IptablesDnatRule
	OK   bool
	MSG  string
}

type CreateSNatResult struct {
	SNat *kubeovnv1.IptablesSnatRule
	OK   bool
	MSG  string
}

// StartVictim model.Victim Preload model.Pod
func StartVictim(victim model.Victim) (map[string]model.Exposes, bool, string) {
	log.Logger.Infof("Starting Victim for Team %d ContestChallenge %d", victim.TeamID, victim.ContestChallengeID)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	// 添加一个独立tag, 防止 NetworkPolicy 影响 frpc 通信
	labels := map[string]string{
		"victim_id":            strconv.Itoa(int(victim.ID)),
		"user_id":              strconv.Itoa(int(victim.UserID)),
		"team_id":              strconv.Itoa(int(victim.TeamID)),
		"challenge_id":         strconv.Itoa(int(victim.ChallengeID)),
		"contest_challenge_id": strconv.Itoa(int(victim.ContestChallengeID)),
		VictimPodTag:           fmt.Sprintf("victim-%s", utils.RandStr(20)),
	}
	subnetMap := make(map[string]*model.Subnet)
	netAttchDefMap := make(map[string]*model.NetAttachDef)
	ipExposesMap := make(map[string]model.Exposes)
	ipExposesMapMutex := &sync.Mutex{}
	if _, ok, msg := CreateNetworkPolicy(ctx, CreateNetworkPolicyOptions{
		Name:        fmt.Sprintf("np-%s", utils.RandStr(20)),
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
	}); !ok {
		return nil, false, msg
	}
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
		if _, ok, msg := CreateVPC(ctx, CreateVPCOptions{
			Name:         victim.VPC.Name,
			Labels:       labels,
			PolicyRoutes: policyRoutes,
		}); !ok {
			return nil, false, msg
		}
		createNADFuncL := make([]func() CreateNADResult, 0)
		createSubnetFuncL := make([]func() CreateSubnetResult, 0)
		createVPCNatGWFuncL := make([]func() CreateVPCNatGWResult, 0)
		createEIPFuncL := make([]func() CreateEIPResult, 0)
		for _, subnet := range victim.VPC.Subnets {
			createNADFuncL = append(createNADFuncL, func() CreateNADResult {
				object, ok, msg := CreateNetAttachDef(ctx, CreateNetAttachDefOptions{
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
				return CreateNADResult{object, ok, msg}
			})
			createSubnetFuncL = append(createSubnetFuncL, func() CreateSubnetResult {
				object, ok, msg := CreateSubnet(ctx, CreateSubnetOptions{
					Name:       subnet.Name,
					Labels:     labels,
					VPC:        victim.VPC.Name,
					CIDR:       subnet.CIDRBlock,
					Gateway:    subnet.Gateway,
					ExcludeIPs: subnet.ExcludeIps,
					Provider:   fmt.Sprintf("%s.%s.ovn", subnet.NetAttachDef.Name, globalNamespace),
				})
				log.Logger.Debugf("CreateSubnet %s: %s", subnet.Name, msg)
				return CreateSubnetResult{object, ok, msg}
			})
			subnetMap[subnet.DefName] = subnet
			netAttchDefMap[subnet.DefName] = subnet.NetAttachDef
			if subnet.NatGateway != nil {
				createVPCNatGWFuncL = append(createVPCNatGWFuncL, func() CreateVPCNatGWResult {
					object, ok, msg := CreateVPCNatGateway(ctx, CreateVPCNatGatewayOptions{
						Name:           subnet.NatGateway.Name,
						Labels:         labels,
						VPC:            victim.VPC.Name,
						Subnet:         subnet.Name,
						LanIP:          subnet.NatGateway.LanIP,
						ExternalSubnet: []string{externalSubnetName},
					})
					log.Logger.Debugf("CreateVPCNatGateway %s: %s", subnet.NatGateway.Name, msg)
					return CreateVPCNatGWResult{object, ok, msg}
				})
				for _, eip := range subnet.NatGateway.EIPs {
					createEIPFuncL = append(createEIPFuncL, func() CreateEIPResult {
						return func() CreateEIPResult {
							e, ok, msg := CreateEIP(ctx, CreateEIPOptions{
								Name:           eip.Name,
								Labels:         labels,
								NatGw:          subnet.NatGateway.Name,
								ExternalSubnet: externalSubnetName,
							})
							log.Logger.Debugf("CreateEIP %s: %s", eip.Name, msg)
							if !ok {
								return CreateEIPResult{e, false, msg}
							}
							eip.IP = e.Spec.V4ip
							createDNatFuncL := make([]func() CreateDNatResult, 0)
							for _, dnat := range eip.DNats {
								createDNatFuncL = append(createDNatFuncL, func() CreateDNatResult {
									object, ok, msg := CreateDNat(ctx, CreateDNatOptions{
										Name:         dnat.Name,
										Labels:       labels,
										EIP:          eip.Name,
										ExternalPort: dnat.ExternalPort,
										InternalPort: dnat.InternalPort,
										InternalIP:   dnat.InternalIP,
										Protocol:     dnat.Protocol,
									})
									log.Logger.Debugf("CreateDNat %s: %s", dnat.Name, msg)
									return CreateDNatResult{object, ok, msg}
								})
							}
							for _, res := range utils.RunFuncLConcurrently(createDNatFuncL) {
								if !res.OK {
									return CreateEIPResult{nil, false, res.MSG}
								}
								port, err := strconv.ParseInt(res.DNat.Spec.ExternalPort, 10, 64)
								if err != nil {
									log.Logger.Warningf("Failed to parse external port: %s", err)
									return CreateEIPResult{nil, false, i18n.UnknownError}
								}
								ipExposesMapMutex.Lock()
								if !slices.ContainsFunc(ipExposesMap[e.Spec.V4ip], func(e model.Expose) bool {
									return int32(port) == e.Port && res.DNat.Spec.Protocol == e.Protocol
								}) {
									ipExposesMap[e.Spec.V4ip] = append(ipExposesMap[e.Spec.V4ip], model.Expose{
										Port:     int32(port),
										Protocol: res.DNat.Spec.Protocol,
									})
								}
								ipExposesMapMutex.Unlock()
							}
							createSNatFuncL := make([]func() CreateSNatResult, 0)
							for _, snat := range eip.SNats {
								createSNatFuncL = append(createSNatFuncL, func() CreateSNatResult {
									object, ok, msg := CreateSNat(ctx, CreateSNatOptions{
										Name:         snat.Name,
										Labels:       labels,
										EIP:          eip.Name,
										InternalCIDR: subnet.CIDRBlock,
									})
									log.Logger.Debugf("CreateSNat %s: %s", snat.Name, msg)
									return CreateSNatResult{object, ok, msg}
								})
							}
							for _, res := range utils.RunFuncLConcurrently(createSNatFuncL) {
								if !res.OK {
									return CreateEIPResult{nil, false, res.MSG}
								}
							}
							return CreateEIPResult{e, true, msg}
						}()
					})
				}
			}
		}
		for _, res := range utils.RunFuncLConcurrently(createSubnetFuncL) {
			if !res.OK {
				return nil, false, res.MSG
			}
		}
		for _, res := range utils.RunFuncLConcurrently(createNADFuncL) {
			if !res.OK {
				return nil, false, res.MSG
			}
		}
		for _, res := range utils.RunFuncLConcurrently(createVPCNatGWFuncL) {
			if !res.OK {
				return nil, false, res.MSG
			}
		}
		for _, res := range utils.RunFuncLConcurrently(createEIPFuncL) {
			if !res.OK {
				return nil, false, res.MSG
			}
		}
	}
	createPodFuncL := make([]func() Result, 0)
	for _, pod := range victim.Pods {
		createPodFuncL = append(createPodFuncL, func() Result {
			return func(pod model.Pod) Result {
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
						flagConfigMap, ok, msg := CreateConfigMap(ctx, CreateConfigMapOptions{
							Name:   fmt.Sprintf("cm-%s", utils.RandStr(20)),
							Labels: labels,
							Data:   map[string]string{filename: volumeFlag},
						})
						if !ok {
							return Result{OK: false, MSG: msg}
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
						return Result{OK: false, MSG: i18n.SubnetNotFound}
					}
					netAttachDef, ok := netAttchDefMap[network.Name]
					if !ok {
						return Result{OK: false, MSG: i18n.NetAttNotFound}
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
				p, ok, msg := CreatePod(ctx, CreatePodOptions{
					Name:        pod.Name,
					Labels:      labels,
					Annotations: annotations,
					Containers:  containers,
					Volumes:     volumes,
				})
				if !ok {
					return Result{OK: false, MSG: msg}
				}
				if len(annotations) == 0 {
					service, ok, msg := CreateService(ctx, CreateServiceOptions{
						Name:     fmt.Sprintf("svc-%s", utils.RandStr(20)),
						Ports:    pod.PodPorts,
						Labels:   labels,
						Selector: labels,
					})
					if !ok {
						log.Logger.Warningf("Failed to create service for generator: %s", msg)
						return Result{OK: false, MSG: msg}
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
				log.Logger.Debugf("Create Pod %s: %s", pod.Name, msg)
				return Result{OK: true, MSG: msg}
			}(pod)
		})
	}
	for _, res := range utils.RunFuncLConcurrently(createPodFuncL) {
		if !res.OK {
			return nil, false, res.MSG
		}
	}
	return ipExposesMap, true, i18n.Success
}

func StopVictim(victim model.Victim) (bool, string) {
	log.Logger.Infof("Stopping Victim for Team %d ContestChallenge %d", victim.TeamID, victim.ContestChallengeID)
	// 不添加独立 tag, 删除时直接删除所有相关资源
	labels := map[string]string{
		"victim_id":            strconv.Itoa(int(victim.ID)),
		"user_id":              strconv.Itoa(int(victim.UserID)),
		"team_id":              strconv.Itoa(int(victim.TeamID)),
		"challenge_id":         strconv.Itoa(int(victim.ChallengeID)),
		"contest_challenge_id": strconv.Itoa(int(victim.ContestChallengeID)),
	}
	for _, endpoint := range victim.Endpoints {
		if err := redis.UnlockFrpsPort(endpoint.IP, endpoint.Port, endpoint.Protocol); err != nil {
			log.Logger.Warningf("Failed to unlock frps port: %s", err)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
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
	if ok, msg := DeleteSubnetList(ctx, labels); !ok {
		return false, msg
	}
	if ok, msg := DeleteNetAttachDefList(ctx, globalNamespace, labels); !ok {
		return false, msg
	}
	if ok, msg := DeleteVPCNatGatewayList(ctx, labels); !ok {
		return false, msg
	}
	if ok, msg := DeleteVPCList(ctx, labels); !ok {
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
	if ok, msg := DeletePodList(ctx, labels); !ok {
		return false, msg
	}
	for _, subnet := range victim.VPC.Subnets {
		if ok, msg := DeleteIPList(ctx, map[string]string{"ovn.kubernetes.io/subnet": subnet.Name}); !ok {
			return false, msg
		}
	}
	return true, i18n.Success
}
