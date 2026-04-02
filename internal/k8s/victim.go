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

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func VictimLabels(victim model.Victim, tags ...map[string]string) map[string]string {
	labels := map[string]string{
		"victim_id":            strconv.Itoa(int(victim.ID)),
		"user_id":              strconv.Itoa(int(victim.UserID)),
		"team_id":              strconv.Itoa(int(victim.TeamID.V)),
		"contest_id":           strconv.Itoa(int(victim.ContestID.V)),
		"challenge_id":         strconv.Itoa(int(victim.ChallengeID)),
		"contest_challenge_id": strconv.Itoa(int(victim.ContestChallengeID.V)),
	}
	if len(tags) > 0 {
		for tag, values := range tags[0] {
			labels[tag] = values
		}
	}
	return labels
}

// StartVictim expects victim.Spec and workload pod records to be preloaded from DB.
func StartVictim(ctx context.Context, victim model.Victim) (model.Victim, model.RetVal) {
	log.Logger.Infof("Starting Victim for Team %d Challenge %d", victim.TeamID.V, victim.ChallengeID)
	labels := VictimLabels(victim, map[string]string{VictimPodTag: VictimPodTag})

	subnetMap, netAttachDefMap, endpoints, ret := createVictimNetworkResources(ctx, &victim, labels)
	if !ret.OK {
		return victim, ret
	}
	victim.Endpoints = endpoints

	endpointsMutex := &sync.Mutex{}
	pods := append([]model.Pod(nil), victim.Pods...)
	wg := utils.NewGroup(ctx)
	for _, pod := range pods {
		pod := pod
		wg.Go(func() error {
			podSpec := pod.Spec
			containers := []corev1.Container{
				{
					Name:    "tcpdump",
					Image:   config.Env.K8S.TCPDumpImage,
					Command: []string{"/bin/sh", "-c", fmt.Sprintf("tcpdump -i any -w /root/mnt/pod-%s.pcap", pod.Name)},
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      nfsVolumeName,
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
					Name: nfsVolumeName,
					VolumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: nfsVolumeName,
						},
					},
				},
			}
			for _, container := range podSpec.Containers {
				volumeMounts := make([]corev1.VolumeMount, 0)
				for path, volumeFlag := range container.Files {
					filename := path[strings.LastIndex(path, "/")+1:]
					flagConfigMap, ret := CreateConfigMap(ctx, CreateConfigMapOptions{
						Name:   fmt.Sprintf("flag-%s", utils.RandStr(20)),
						Labels: labels,
						Data:   map[string]string{filename: volumeFlag},
					})
					if err, ok := ret.Attr["Error"]; ok && !ret.OK {
						return errors.New(err.(string))
					}
					volumeName := fmt.Sprintf("flag-%s", utils.RandStr(10))
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

				envs := make([]corev1.EnvVar, 0, len(container.Environment))
				for key, value := range container.Environment {
					envs = append(envs, corev1.EnvVar{Name: key, Value: value})
				}

				ports := make([]corev1.ContainerPort, 0, len(container.Exposes))
				for _, p := range container.Exposes {
					ports = append(ports, corev1.ContainerPort{ContainerPort: p.Port})
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
			for i, network := range podSpec.Networks {
				subnet, ok := subnetMap[network.Name]
				if !ok {
					return fmt.Errorf("subnet %s not found", network.Name)
				}
				netAttachDef, ok := netAttachDefMap[network.Name]
				if !ok {
					return fmt.Errorf("netAttachDef %s not found", network.Name)
				}
				if i == 0 {
					annotations["ovn.kubernetes.io/logical_switch"] = subnet.Name
					annotations["ovn.kubernetes.io/ip_address"] = network.IP
					annotations["v1.multus-cni.io/default-network"] = fmt.Sprintf("%s/%s", globalNamespace, netAttachDef.Name)
				} else {
					annotations["k8s.v1.cni.cncf.io/networks"] += fmt.Sprintf(",%s/%s", globalNamespace, netAttachDef.Name)
					annotations["k8s.v1.cni.cncf.io/networks"] = strings.Trim(annotations["k8s.v1.cni.cncf.io/networks"], ",")
				}
				annotations[fmt.Sprintf("%s.%s.ovn.kubernetes.io/logical_switch", netAttachDef.Name, globalNamespace)] = subnet.Name
				annotations[fmt.Sprintf("%s.%s.ovn.kubernetes.io/ip_address", netAttachDef.Name, globalNamespace)] = network.IP
				if network.External {
					annotations[fmt.Sprintf("%s.%s.ovn.kubernetes.io/routes", netAttachDef.Name, globalNamespace)] = fmt.Sprintf("[{\"gw\":\"%s\"}]", network.Gateway)
				}
			}

			pOptions := CreatePodOptions{
				Name:        pod.Name,
				Labels:      labels,
				Annotations: annotations,
				Containers:  containers,
				Volumes:     volumes,
			}
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
					Ports:    podSpec.ServicePorts,
					Labels:   labels,
					Selector: labels,
				})
				if err, ok := ret.Attr["Error"]; ok && !ret.OK {
					return errors.New(err.(string))
				}
				endpointsMutex.Lock()
				for _, port := range service.Spec.Ports {
					endpoint := model.Endpoint{
						IP:       p.Status.HostIP,
						Port:     port.NodePort,
						Protocol: string(port.Protocol),
					}
					if !slices.ContainsFunc(victim.Endpoints, func(e model.Endpoint) bool {
						return e.IP == endpoint.IP && e.Port == endpoint.Port && strings.EqualFold(e.Protocol, endpoint.Protocol)
					}) {
						victim.Endpoints = append(victim.Endpoints, endpoint)
					}
				}
				endpointsMutex.Unlock()
			}

			log.Logger.Debugf("Create Pod %s: %s", pod.Name, ret.Msg)
			return nil
		})
	}

	if err := wg.Wait(); err != nil {
		return victim, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}}
	}

	victim.Resources.NetworkPlan = victim.Spec.NetworkPlan
	victim.Resources.PodNames = make(model.StringList, 0, len(pods))
	for _, pod := range pods {
		victim.Resources.PodNames = append(victim.Resources.PodNames, pod.Name)
	}
	victim.ExposedEndpoints = append(model.Endpoints(nil), victim.Endpoints...)
	if config.Env.K8S.Frp.On {
		return AddFrpc(ctx, victim)
	}

	return victim, model.SuccessRetVal()
}

func createVictimNetworkResources(
	ctx context.Context,
	victim *model.Victim,
	labels map[string]string,
) (map[string]*model.Subnet, map[string]*model.NetAttachDef, model.Endpoints, model.RetVal) {
	subnetMap := make(map[string]*model.Subnet)
	netAttachDefMap := make(map[string]*model.NetAttachDef)
	endpoints := make(model.Endpoints, 0)
	endpointsMutex := &sync.Mutex{}
	wg := utils.NewGroup(ctx)

	wg.Go(func() error {
		name := fmt.Sprintf("np-%s", utils.RandStr(20))
		_, ret := CreateNetworkPolicy(ctx, CreateNetworkPolicyOptions{
			Name:        name,
			Labels:      labels,
			MatchLabels: labels,
			From: func() []*netv1.IPBlock {
				tmp := make([]*netv1.IPBlock, 0)
				for _, p := range victim.Spec.NetworkPolicies {
					tmp = append(tmp, p.From...)
				}
				return tmp
			}(),
			To: func() []*netv1.IPBlock {
				tmp := make([]*netv1.IPBlock, 0)
				for _, p := range victim.Spec.NetworkPolicies {
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

	if victim.Spec.NetworkPlan.Name == "" {
		if err := wg.Wait(); err != nil {
			return nil, nil, nil, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}}
		}
		return subnetMap, netAttachDefMap, endpoints, model.SuccessRetVal()
	}

	policyRoutes := make([]*kubeovnv1.PolicyRoute, 0)
	for _, subnet := range victim.Spec.NetworkPlan.Subnets {
		if subnet != nil && subnet.NatGateway != nil {
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
			Name:         victim.Spec.NetworkPlan.Name,
			Labels:       labels,
			PolicyRoutes: policyRoutes,
		})
		log.Logger.Debugf("Create VPC %s: %s", victim.Spec.NetworkPlan.Name, ret.Msg)
		if err, ok := ret.Attr["Error"]; ok && !ret.OK {
			return errors.New(err.(string))
		}
		return nil
	})

	for _, subnet := range victim.Spec.NetworkPlan.Subnets {
		if subnet == nil {
			continue
		}
		subnet := subnet
		subnetMap[subnet.DefName] = subnet
		netAttachDefMap[subnet.DefName] = subnet.NetAttachDef

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
				VPC:        victim.Spec.NetworkPlan.Name,
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

		if subnet.NatGateway == nil {
			continue
		}

		wg.Go(func() error {
			_, ret := CreateVPCNatGateway(ctx, CreateVPCNatGatewayOptions{
				Name:           subnet.NatGateway.Name,
				Labels:         labels,
				VPC:            victim.Spec.NetworkPlan.Name,
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

		wg.Go(func() error {
			e, ret := CreateEIP(ctx, CreateEIPOptions{
				Name:           subnet.NatGateway.EIP.Name,
				Labels:         labels,
				NatGw:          subnet.NatGateway.Name,
				ExternalSubnet: externalSubnetName,
			})
			log.Logger.Debugf("Create EIP %s: %s", subnet.NatGateway.EIP.Name, ret.Msg)
			if !ret.OK {
				if err, ok := ret.Attr["Error"].(string); ok {
					return errors.New(err)
				}
				return fmt.Errorf("create EIP %s failed: %s", subnet.NatGateway.EIP.Name, ret.Msg)
			}

			for _, dnat := range subnet.NatGateway.EIP.DNats {
				_, ret = CreateDNat(ctx, CreateDNatOptions{
					Name:         dnat.Name,
					Labels:       labels,
					EIP:          subnet.NatGateway.EIP.Name,
					ExternalPort: strconv.Itoa(int(dnat.ExternalPort)),
					InternalPort: strconv.Itoa(int(dnat.InternalPort)),
					InternalIP:   dnat.InternalIP,
					Protocol:     dnat.Protocol,
				})
				log.Logger.Debugf("Create DNat %s: %s", dnat.Name, ret.Msg)
				if err, ok := ret.Attr["Error"]; ok && !ret.OK {
					return errors.New(err.(string))
				}
				endpointsMutex.Lock()
				endpoint := model.Endpoint{
					IP:       e.Spec.V4ip,
					Port:     dnat.ExternalPort,
					Protocol: dnat.Protocol,
				}
				if !slices.ContainsFunc(victim.Endpoints, func(e model.Endpoint) bool {
					return e.IP == endpoint.IP && e.Port == endpoint.Port && strings.EqualFold(e.Protocol, endpoint.Protocol)
				}) {
					victim.Endpoints = append(victim.Endpoints, endpoint)
				}
				endpointsMutex.Unlock()
			}

			for _, snat := range subnet.NatGateway.EIP.SNats {
				_, ret = CreateSNat(ctx, CreateSNatOptions{
					Name:         snat.Name,
					Labels:       labels,
					EIP:          subnet.NatGateway.EIP.Name,
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

	if err := wg.Wait(); err != nil {
		return nil, nil, nil, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}}
	}

	return subnetMap, netAttachDefMap, endpoints, model.SuccessRetVal()
}

func StopVictim(ctx context.Context, victim model.Victim) model.RetVal {
	log.Logger.Infof("Stopping Victim for Team %d Challenge %d", victim.TeamID.V, victim.ChallengeID)
	labels := VictimLabels(victim)
	for _, endpoint := range victim.ExposedEndpoints {
		redis.UnlockFrpsPort(endpoint.IP, endpoint.Port, endpoint.Protocol)
	}
	var firstErr model.RetVal
	tryDelete := func(ret model.RetVal) {
		if !ret.OK && firstErr.OK {
			firstErr = ret
		}
	}
	firstErr = model.SuccessRetVal()
	tryDelete(DeleteDNatList(ctx, labels))
	tryDelete(DeleteSNatList(ctx, labels))
	tryDelete(DeleteEIPList(ctx, labels))
	tryDelete(DeleteSubnetList(ctx, labels))
	tryDelete(DeleteNetAttachDefList(ctx, globalNamespace, labels))
	tryDelete(DeleteVPCNatGatewayList(ctx, labels))
	tryDelete(DeleteVPCList(ctx, labels))
	tryDelete(DeleteConfigMapList(ctx, labels))
	tryDelete(DeleteNetworkPolicyList(ctx, labels))
	tryDelete(DeleteEndpointList(ctx, labels))
	tryDelete(DeleteServiceList(ctx, labels))
	tryDelete(DeletePodList(ctx, labels))
	for _, subnet := range victim.Spec.NetworkPlan.Subnets {
		tryDelete(DeleteIPList(ctx, map[string]string{"ovn.kubernetes.io/subnet": subnet.Name}))
	}
	return firstErr
}
