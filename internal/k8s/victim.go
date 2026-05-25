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
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	corev1 "k8s.io/api/core/v1"
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

func findExposeDisplayName(exposes model.Exposes, port int32, protocol string) string {
	for _, expose := range exposes {
		if expose.Port == port && strings.EqualFold(expose.Protocol, protocol) {
			return expose.Published
		}
	}
	return ""
}

func podServiceName(podSpec model.PodSpec) string {
	if len(podSpec.Containers) > 0 {
		return podSpec.Containers[0].Name
	}
	return podSpec.Key
}

// StartVictim expects victim.Spec and workload pod records to be preloaded from DB.
func StartVictim(ctx context.Context, victim model.Victim) (model.Victim, model.RetVal) {
	log.Logger.Debugf(
		"Creating victim k8s resources: victim_id=%d team_id=%d challenge_id=%d pods=%d vpc=%t frp=%t",
		victim.ID, victim.TeamID.V, victim.ChallengeID, len(victim.Pods), victim.Spec.NetworkPlan.Name != "", config.Env.K8S.Frp.On,
	)
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
		wg.Go(func() error {
			// VPC 模式下, 支持多 NetworkPolicy 根据 Labels 绑定到指定 Pod 上
			podLabels := make(map[string]string, len(labels)+1)
			for key, value := range labels {
				podLabels[key] = value
			}
			if serviceName := podServiceName(pod.Spec); serviceName != "" {
				podLabels[ServiceLabel] = serviceName
			}

			networks := make([]Network, 0)
			for _, network := range pod.Spec.Networks {
				networks = append(networks, Network{
					Interface:    network.Attachment.Name,
					IPv4:         network.Attachment.IP,
					MAC:          network.Attachment.MAC,
					Gateway:      network.Definition.Gateway,
					Subnet:       subnetMap[network.Definition.Name].Name,
					NetAttachDef: netAttachDefMap[network.Definition.Name].Name,
				})
			}

			// 当 VPC 模式下, 一个 Pod 只有一个 Container
			if victim.Spec.NetworkPlan.Name != "" && pod.Spec.Containers[0].KubeVirt {
				container := pod.Spec.Containers[0]
				_, ret = CreateVM(ctx, CreateVMOptions{
					Name:        pod.Name,
					Labels:      podLabels,
					Image:       container.Image,
					Bootloader:  container.Bootloader,
					SecureBoot:  container.SecureBoot,
					CPUMillis:   container.Resources.CPUMillis,
					MemoryBytes: container.Resources.MemoryBytes,
					UserData:    container.UserData,
					Networks:    networks,
				})
				if err, ok := ret.Attr["Error"]; ok && !ret.OK {
					return errors.New(err.(string))
				}
				log.Logger.Debugf("Created victim vm: victim_id=%d vm=%s", victim.ID, pod.Name)
				return nil
			}
			capture := corev1.Container{
				Name:            "capture",
				Image:           config.Env.K8S.CaptureImage,
				ImagePullPolicy: corev1.PullIfNotPresent,
				Command:         []string{"/bin/sh", "-c"},
				VolumeMounts: []corev1.VolumeMount{
					{
						Name:      nfsVolumeName,
						MountPath: "/root/mnt",
						SubPath: strings.TrimPrefix(
							strings.TrimPrefix(victim.TrafficBasePath(), config.Env.Path), "/",
						),
					},
				},
				SecurityContext: &corev1.SecurityContext{
					Capabilities: &corev1.Capabilities{
						Add: []corev1.Capability{"NET_RAW", "SYS_ADMIN"},
					},
				},
				Stdin: true,
				TTY:   true,
			}
			command := fmt.Sprintf("rustnet -i any --pcap-export /root/mnt/pod-%s.pcap", pod.Name)
			if _, err := os.Stat(filepath.Join(config.Env.Path, "GeoLite2-City.mmdb")); err == nil {
				command = fmt.Sprintf("rustnet -i any --geoip-city /root/GeoLite2-City.mmdb --pcap-export /root/mnt/pod-%s.pcap", pod.Name)
				capture.VolumeMounts = append(capture.VolumeMounts, corev1.VolumeMount{
					Name:      nfsVolumeName,
					MountPath: "/root/GeoLite2-City.mmdb",
					SubPath:   "GeoLite2-City.mmdb",
					ReadOnly:  true,
				})
			}
			capture.Command = append(capture.Command, command)
			containers := []corev1.Container{capture}
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
			// 兼容非 VPC 模式下, 一个 Pod 多个 Container, 使用循环来处理
			for _, container := range pod.Spec.Containers {
				volumeMounts := make([]corev1.VolumeMount, 0)
				for _, fileMount := range container.FileMounts {
					path := fileMount.Path
					filename := path[strings.LastIndex(path, "/")+1:]
					flagConfigMap, ret := CreateConfigMap(ctx, CreateConfigMapOptions{
						Name:   fmt.Sprintf("flag-%s", utils.RandHexStr(20)),
						Labels: labels,
						Data:   map[string]string{filename: fileMount.Content},
					})
					if err, ok := ret.Attr["Error"]; ok && !ret.OK {
						return errors.New(err.(string))
					}
					volumeName := fmt.Sprintf("flag-%s", utils.RandHexStr(10))
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
				if container.Resources.CPUMillis > 0 {
					limit[corev1.ResourceCPU] = resource.MustParse(strconv.FormatInt(container.Resources.CPUMillis, 10) + "m")
				}
				if container.Resources.MemoryBytes > 0 {
					limit[corev1.ResourceMemory] = resource.MustParse(strconv.FormatInt(container.Resources.MemoryBytes, 10))
				}

				tmp := corev1.Container{
					Name:            container.Name,
					Image:           container.Image,
					ImagePullPolicy: corev1.PullIfNotPresent,
					Env:             envs,
					Ports:           ports,
					VolumeMounts:    volumeMounts,
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

			pOptions := CreatePodOptions{
				Name:       pod.Name,
				Labels:     podLabels,
				Networks:   networks,
				Containers: containers,
				Volumes:    volumes,
			}

			p, ret := CreatePod(ctx, pOptions)
			if err, ok := ret.Attr["Error"]; ok && !ret.OK {
				return errors.New(err.(string))
			}

			if len(pod.Spec.ServicePorts) > 0 {
				service, ret := CreateService(ctx, CreateServiceOptions{
					Name:     fmt.Sprintf("svc-%s", utils.RandHexStr(20)),
					Ports:    pod.Spec.ServicePorts,
					Labels:   labels,
					Selector: podLabels,
				})
				if err, ok := ret.Attr["Error"]; ok && !ret.OK {
					return errors.New(err.(string))
				}
				endpointsMutex.Lock()
				for _, port := range service.Spec.Ports {
					endpoint := model.Endpoint{
						Name:     findExposeDisplayName(pod.Spec.ServicePorts, port.Port, string(port.Protocol)),
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

			log.Logger.Debugf("Created victim pod: victim_id=%d pod=%s", victim.ID, pod.Name)
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

	log.Logger.Debugf("Victim k8s resources created: victim_id=%d pods=%d endpoints=%d", victim.ID, len(victim.Resources.PodNames), len(victim.ExposedEndpoints))
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
	wg := utils.NewGroup(ctx)

	createNetworkPolicy := func(labels map[string]string, policies model.NetworkPolicies) error {
		name := fmt.Sprintf("np-%s", utils.RandHexStr(20))
		_, ret := CreateNetworkPolicy(ctx, CreateNetworkPolicyOptions{
			Name:     name,
			Labels:   labels,
			Policies: policies,
		})
		if err, ok := ret.Attr["Error"]; ok && !ret.OK {
			return errors.New(err.(string))
		}
		log.Logger.Debugf("Created victim network policy: victim_id=%d network_policy=%s", victim.ID, name)
		return nil
	}

	if victim.Spec.NetworkPlan.Name == "" {
		wg.Go(func() error {
			return createNetworkPolicy(labels, victim.Spec.NetworkPolicies)
		})
	} else {
		for _, podSpec := range victim.Spec.Pods {
			serviceName := podServiceName(podSpec)
			matchLabels := make(map[string]string, len(labels)+1)
			for key, value := range labels {
				matchLabels[key] = value
			}
			if serviceName != "" {
				matchLabels[ServiceLabel] = serviceName
			}
			policies := make(model.NetworkPolicies, 0)
			for _, policy := range victim.Spec.NetworkPolicies {
				if policy.Service == serviceName {
					policies = append(policies, policy)
				}
			}
			wg.Go(func() error {
				return createNetworkPolicy(matchLabels, policies)
			})
		}
	}

	if victim.Spec.NetworkPlan.Name == "" {
		if err := wg.Wait(); err != nil {
			return nil, nil, nil, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}}
		}
		return subnetMap, netAttachDefMap, endpoints, model.SuccessRetVal()
	}

	wg.Go(func() error {
		_, ret := CreateVPC(ctx, CreateVPCOptions{
			Name:   victim.Spec.NetworkPlan.Name,
			Labels: labels,
		})
		if err, ok := ret.Attr["Error"]; ok && !ret.OK {
			return errors.New(err.(string))
		}
		log.Logger.Debugf("Created victim vpc: victim_id=%d vpc=%s", victim.ID, victim.Spec.NetworkPlan.Name)
		return nil
	})

	for _, subnet := range victim.Spec.NetworkPlan.Subnets {
		if subnet == nil {
			continue
		}
		subnetMap[subnet.DefName] = subnet
		netAttachDefMap[subnet.DefName] = subnet.NetAttachDef

		wg.Go(func() error {
			_, ret := CreateNetAttachDef(ctx, CreateNetAttachDefOptions{
				Name:   subnet.NetAttachDef.Name,
				Labels: labels,
			})
			if err, ok := ret.Attr["Error"]; ok && !ret.OK {
				return errors.New(err.(string))
			}
			log.Logger.Debugf("Created victim net attach def: victim_id=%d net_attach_def=%s", victim.ID, subnet.NetAttachDef.Name)
			return nil
		})

		wg.Go(func() error {
			_, ret := CreateSubnet(ctx, CreateSubnetOptions{
				Name:         subnet.Name,
				Labels:       labels,
				VPC:          victim.Spec.NetworkPlan.Name,
				CIDR:         subnet.CIDRBlock,
				Gateway:      subnet.Gateway,
				ExcludeIPs:   subnet.ExcludeIps,
				NetAttachDef: subnet.NetAttachDef.Name,
			})
			if err, ok := ret.Attr["Error"]; ok && !ret.OK {
				return errors.New(err.(string))
			}
			log.Logger.Debugf("Created victim subnet: victim_id=%d subnet=%s", victim.ID, subnet.Name)
			return nil
		})
	}

	if err := wg.Wait(); err != nil {
		return nil, nil, nil, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}}
	}

	return subnetMap, netAttachDefMap, endpoints, model.SuccessRetVal()
}

func StopVictim(ctx context.Context, victim model.Victim) model.RetVal {
	log.Logger.Debugf(
		"Deleting victim k8s resources: victim_id=%d team_id=%d challenge_id=%d exposed_endpoints=%d",
		victim.ID, victim.TeamID.V, victim.ChallengeID, len(victim.ExposedEndpoints),
	)
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
	tryDelete(DeleteSubnetCollection(ctx, labels))
	tryDelete(DeleteNetAttachDefCollection(ctx, globalNamespace, labels))
	tryDelete(DeleteVPCCollection(ctx, labels))
	tryDelete(DeleteConfigMapCollection(ctx, labels))
	tryDelete(DeleteNetworkPolicyCollection(ctx, labels))
	tryDelete(DeleteEndpointCollection(ctx, labels))
	tryDelete(DeleteServiceCollection(ctx, labels))
	tryDelete(DeletePodCollection(ctx, labels))
	tryDelete(DeleteVMCollection(ctx, labels))
	for _, subnet := range victim.Spec.NetworkPlan.Subnets {
		tryDelete(DeleteIPCollection(ctx, map[string]string{"ovn.kubernetes.io/subnet": subnet.Name}))
	}
	tryDelete(WaitVictimPodsDeleted(ctx, victim))
	if firstErr.OK {
		log.Logger.Debugf("Deleted victim k8s resources: victim_id=%d", victim.ID)
	} else {
		log.Logger.Warningf("Victim k8s cleanup incomplete: victim_id=%d error=%s", victim.ID, firstErr.Msg)
	}
	return firstErr
}

func WaitVictimPodsDeleted(ctx context.Context, victim model.Victim) model.RetVal {
	labels := VictimLabels(victim, map[string]string{VictimPodTag: VictimPodTag})
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		podList, ret := ListPods(ctx, labels)
		if !ret.OK {
			return ret
		}
		if len(podList.Items) == 0 {
			return model.SuccessRetVal()
		}
		select {
		case <-ctx.Done():
			return model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": ctx.Err().Error()}}
		case <-ticker.C:
		}
	}
}
