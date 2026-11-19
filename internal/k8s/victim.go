package k8s

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"CBCTF/internal/config"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
	"CBCTF/internal/utils"
)

func VictimLabels(victim model.Victim, tags ...map[string]string) map[string]string {
	labels := map[string]string{
		"cbctf.io/namespace":   globalNamespace,
		"victim_id":            strconv.Itoa(int(victim.ID)),
		"user_id":              strconv.Itoa(int(victim.UserID)),
		"team_id":              strconv.Itoa(int(victim.TeamID.V)),
		"contest_id":           strconv.Itoa(int(victim.ContestID.V)),
		"challenge_id":         strconv.Itoa(int(victim.ChallengeID)),
		"contest_challenge_id": strconv.Itoa(int(victim.ContestChallengeID.V)),
	}
	if len(tags) > 0 {
		maps.Copy(labels, tags[0])
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
	// FRP without a frps server falls back to the plain exposure path, the
	// system setting stays untouched.
	if victim.Spec.FrpEnabled && len(config.Env.K8S.Frp.Frps) == 0 {
		log.Logger.Warningf(
			"FRP is enabled but no frps server is configured, starting victim %d without FRP exposure",
			victim.ID,
		)
		victim.Spec.FrpEnabled = false
	}
	labels := VictimLabels(victim, map[string]string{RoleLabel: VictimPodTag})
	victim.Resources.UIDs = make(map[string]string)
	if config.Env.K8S.CaptureEnabled {
		// The capture sidecar writes into the traffic directory, which has to
		// exist and belong to the capture user before the pods are created.
		if err := prepareCaptureTrafficDir(victim); err != nil {
			log.Logger.Warningf(
				"Failed to prepare traffic directory for victim %d: %s", victim.ID, err,
			)
		}
	}
	workloadCtx, ret := createVictimRoot(ctx, victim, "workloads")
	if !ret.OK {
		return victim, ret
	}
	networkCtx, ret := createVictimRoot(ctx, victim, "network")
	if !ret.OK {
		return victim, ret
	}

	subnetMap, netAttachDefMap, endpoints, ret := createVictimNetworkResources(networkCtx, &victim, labels)
	if !ret.OK {
		return victim, ret
	}
	victim.Endpoints = endpoints
	ctx = workloadCtx

	endpointsMutex := &sync.Mutex{}
	pods := append([]model.Pod(nil), victim.Pods...)
	wg := utils.NewGroup(ctx)
	for _, pod := range pods {
		wg.Go(func() error {
			ctx := wg.Context()
			// In VPC mode each pod gets its own service label, so several
			// NetworkPolicies can bind to different pods.
			podLabels := make(map[string]string, len(labels)+1)
			maps.Copy(podLabels, labels)
			if serviceName := podServiceName(pod.Spec); serviceName != "" {
				podLabels[ServiceLabel] = serviceName
			}

			networks := make([]Network, 0)
			for _, network := range pod.Spec.Networks {
				networks = append(networks, Network{
					Interface:    network.Attachment.Name,
					IPv4:         network.Attachment.IP,
					MAC:          network.Attachment.MAC,
					Subnet:       subnetMap[network.Definition.Name].Name,
					NetAttachDef: netAttachDefMap[network.Definition.Name].Name,
				})
			}

			// In VPC mode a pod carries a single container.
			if victim.Spec.NetworkPlan.Name != "" && pod.Spec.Containers[0].KubeVirt {
				container := pod.Spec.Containers[0]
				vm, ret := CreateVM(ctx, CreateVMOptions{
					SubmitOnly:  true,
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
				if !ret.OK {
					return ret
				}
				endpointsMutex.Lock()
				victim.Resources.UIDs[pod.Name] = string(vm.UID)
				endpointsMutex.Unlock()
				log.Logger.Debugf("Created victim vm: victim_id=%d vm=%s", victim.ID, pod.Name)
				return nil
			}
			capture := captureContainer(victim, fmt.Sprintf("pod-%s.pcap", pod.Name))
			var containers []corev1.Container
			var volumes []corev1.Volume
			if config.Env.K8S.CaptureEnabled {
				containers = append(containers, capture)
				volumes = append(volumes, corev1.Volume{
					Name: nfsVolumeName,
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: nfsVolumeClaim,
					},
				})
			}
			fileData, fileMounts := challengeFileMounts(pod.Spec.Containers)
			for containerIndex, container := range pod.Spec.Containers {
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
					VolumeMounts:    fileMounts[containerIndex],
					Resources: corev1.ResourceRequirements{
						Limits:   limit,
						Requests: workloadRequests(limit),
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
			if len(fileData) > 0 {
				cm, ret := CreateConfigMap(ctx, CreateConfigMapOptions{
					Name: victimResourceName("vol", victim.ID, pod.Spec.Key), Labels: labels, Data: fileData,
				})
				if !ret.OK {
					return ret
				}
				volumes = append(volumes, corev1.Volume{Name: "challenge-files", ConfigMap: &corev1.ConfigMapVolumeSource{Name: cm.Name}})
			}

			pOptions := CreatePodOptions{
				SubmitOnly:         true,
				Name:               pod.Name,
				Labels:             podLabels,
				Networks:           networks,
				Containers:         containers,
				Volumes:            volumes,
				ServiceAccountName: runtimeServiceAccount(),
			}

			p, ret := CreatePod(ctx, pOptions)
			if !ret.OK {
				return ret
			}
			endpointsMutex.Lock()
			victim.Resources.UIDs[pod.Name] = string(p.UID)
			endpointsMutex.Unlock()

			if len(pod.Spec.ServicePorts) > 0 {
				service, ret := CreateService(ctx, CreateServiceOptions{
					ClusterIP: victim.Spec.FrpEnabled,
					Name:      victimResourceName("svc", victim.ID, pod.Spec.Key),
					Ports:     pod.Spec.ServicePorts,
					Labels:    labels,
					Selector:  podLabels,
				})
				if !ret.OK {
					return ret
				}
				endpointsMutex.Lock()
				for _, port := range service.Spec.Ports {
					endpoint := model.Endpoint{
						Name:     findExposeDisplayName(pod.Spec.ServicePorts, port.Port, string(port.Protocol)),
						IP:       p.Status.HostIP,
						Port:     port.NodePort,
						Protocol: string(port.Protocol),
					}
					if victim.Spec.FrpEnabled {
						endpoint.IP, endpoint.Port = service.Spec.ClusterIP, port.Port
					} else {
						victim.Resources.NodePorts = append(victim.Resources.NodePorts, model.NodePortEndpoint{PodName: pod.Name, Endpoint: endpoint})
						continue
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
		return victim, model.FailRetVal(i18n.Common.UnknownError, "Victim", fmt.Sprintf("victim_id=%d", victim.ID), err)
	}

	victim.Resources.NetworkPlan = victim.Spec.NetworkPlan
	victim.Resources.PodNames = make([]string, 0, len(pods))
	for _, pod := range pods {
		victim.Resources.PodNames = append(victim.Resources.PodNames, pod.Name)
	}
	victim.ExposedEndpoints = append(model.Endpoints(nil), victim.Endpoints...)
	if victim.Spec.FrpEnabled {
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
		name := victimResourceName("np", victim.ID, labels[ServiceLabel])
		_, ret := CreateNetworkPolicy(wg.Context(), CreateNetworkPolicyOptions{
			Name:     name,
			Labels:   labels,
			Policies: policies,
		})
		if !ret.OK {
			return ret
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
			maps.Copy(matchLabels, labels)
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
			return nil, nil, nil, model.FailRetVal(i18n.Common.UnknownError, "VictimNetwork", fmt.Sprintf("victim_id=%d", victim.ID), err)
		}
		return subnetMap, netAttachDefMap, endpoints, model.SuccessRetVal()
	}

	vpc, ret := CreateVPC(ctx, CreateVPCOptions{
		Name:   victim.Spec.NetworkPlan.Name,
		Labels: labels,
	})
	if !ret.OK {
		_ = wg.Wait()
		return nil, nil, nil, ret
	}
	subnetCtx := withResourceOwner(wg.Context(), metav1.OwnerReference{APIVersion: "kubeovn.io/v1", Kind: "Vpc", Name: vpc.Name, UID: vpc.UID})

	for _, subnet := range victim.Spec.NetworkPlan.Subnets {
		if subnet == nil {
			continue
		}
		subnetMap[subnet.DefName] = subnet
		netAttachDefMap[subnet.DefName] = subnet.NetAttachDef

		wg.Go(func() error {
			ctx := wg.Context()
			_, ret := CreateNetAttachDef(ctx, CreateNetAttachDefOptions{
				Name:   subnet.NetAttachDef.Name,
				Labels: labels,
			})
			if !ret.OK {
				return ret
			}
			log.Logger.Debugf("Created victim net attach def: victim_id=%d net_attach_def=%s", victim.ID, subnet.NetAttachDef.Name)
			return nil
		})

		wg.Go(func() error {
			_, ret := CreateSubnet(subnetCtx, CreateSubnetOptions{
				Name:         subnet.Name,
				Labels:       labels,
				VPC:          victim.Spec.NetworkPlan.Name,
				CIDR:         subnet.CIDRBlock,
				Gateway:      subnet.Gateway,
				ExcludeIPs:   subnet.ExcludeIps,
				NetAttachDef: subnet.NetAttachDef.Name,
			})
			if !ret.OK {
				return ret
			}
			log.Logger.Debugf("Created victim subnet: victim_id=%d subnet=%s", victim.ID, subnet.Name)
			return nil
		})
	}

	if err := wg.Wait(); err != nil {
		return nil, nil, nil, model.FailRetVal(i18n.Common.UnknownError, "VictimNetwork", fmt.Sprintf("victim_id=%d", victim.ID), err)
	}

	return subnetMap, netAttachDefMap, endpoints, model.SuccessRetVal()
}

func StopVictim(ctx context.Context, victim model.Victim) model.RetVal {
	log.Logger.Debugf(
		"Deleting victim k8s resources: victim_id=%d team_id=%d challenge_id=%d exposed_endpoints=%d",
		victim.ID, victim.TeamID.V, victim.ChallengeID, len(victim.ExposedEndpoints),
	)
	batch := model.NewBatch(len(victim.Spec.NetworkPlan.Subnets) + len(victim.ExposedEndpoints))
	// The workload tree goes first, the network tree stays until every
	// workload (including FRP/capture sidecars) has disappeared.
	if err := deleteVictimRoot(ctx, victim, "workloads"); err != nil {
		return model.FailRetVal(i18n.K8S.DeleteError, "Victim", fmt.Sprintf("victim_id=%d", victim.ID), err)
	}
	if ret := WaitVictimPodsDeleted(ctx, victim); !ret.OK {
		return ret
	}
	if err := deleteVictimRoot(ctx, victim, "network"); err != nil {
		return model.FailRetVal(i18n.K8S.DeleteError, "VictimNetwork", fmt.Sprintf("victim_id=%d", victim.ID), err)
	}
	for _, subnet := range victim.Spec.NetworkPlan.Subnets {
		if ctx.Err() != nil {
			return batch.Result(ctx)
		}
		if subnet != nil {
			batch.Record(subnet.Name, "delete_subnet_ips", DeleteIPCollection(ctx, map[string]string{"ovn.kubernetes.io/subnet": subnet.Name}))
		} else {
			batch.Skip("", "no_subnet")
		}
	}
	if batch.Failed > 0 {
		return batch.Result(ctx)
	}
	if err := deleteVictimVPC(ctx, victim); err != nil {
		return model.FailRetVal(i18n.K8S.DeleteError, "VPC", victim.Spec.NetworkPlan.Name, err)
	}
	for _, endpoint := range victim.ExposedEndpoints {
		if ctx.Err() != nil {
			return batch.Result(ctx)
		}
		ret := redis.UnlockFrpsPort(endpoint.IP, endpoint.Port, endpoint.Protocol)
		if !ret.OK {
			log.Logger.Warningf("Failed to release FRPS reservation: victim_id=%d host=%s port=%d protocol=%s reason=%s", victim.ID, endpoint.IP, endpoint.Port, endpoint.Protocol, ret.Msg)
		}
		batch.Record(fmt.Sprintf("%s:%d/%s", endpoint.IP, endpoint.Port, endpoint.Protocol), "release_frps_port", ret)
	}
	result := batch.Result(ctx)
	if result.OK {
		log.Logger.Debugf("Deleted victim k8s resources: victim_id=%d", victim.ID)
	} else {
		log.Logger.Warningf("Victim k8s cleanup incomplete: victim_id=%d error=%s", victim.ID, result.Msg)
	}
	return result
}

func WaitVictimPodsDeleted(ctx context.Context, victim model.Victim) model.RetVal {
	labels := VictimLabels(victim)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		podList, ret := ListPods(ctx, labels)
		if !ret.OK {
			return ret
		}
		vmsGone := true
		for _, pod := range victim.Pods {
			if len(pod.Spec.Containers) == 0 || !pod.Spec.Containers[0].KubeVirt {
				continue
			}
			vm, err := cachedVM(ctx, pod.Name)
			if err != nil {
				return model.FailRetVal(i18n.K8S.GetError, "VirtualMachine", pod.Name, err)
			}
			if vm != nil {
				vmsGone = false
			} else {
				_, ret := GetVM(ctx, pod.Name)
				if !ret.OK && ret.Msg != i18n.K8S.NotFound {
					return ret
				}
				if ret.OK {
					vmsGone = false
				}
			}
		}
		if len(podList.Items) == 0 && vmsGone {
			// Cache absence alone cannot prove deletion after a recent create.
			direct, ret := listPodsDirect(ctx, labels)
			if !ret.OK {
				return ret
			}
			if len(direct.Items) == 0 {
				return model.SuccessRetVal()
			}
		}
		select {
		case <-ctx.Done():
			return model.FailRetVal(i18n.Common.UnknownError, "Victim", fmt.Sprintf("victim_id=%d", victim.ID), ctx.Err())
		case <-ticker.C:
		}
	}
}
