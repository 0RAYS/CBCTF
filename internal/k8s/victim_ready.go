package k8s

import (
	"context"
	"fmt"

	"CBCTF/internal/model"
)

// StartupFailure means the observed workload state is terminal. API and cache
// failures stay ordinary errors, so a transient read outage does not destroy
// workloads.
type StartupFailure struct{ Err error }

func (e *StartupFailure) Error() string { return e.Err.Error() }
func (e *StartupFailure) Unwrap() error { return e.Err }

// VictimReady evaluates readiness from the shared caches, because the start
// workers only submit objects.
func VictimReady(ctx context.Context, victim *model.Victim) (bool, error) {
	if !victim.Resources.Submitted {
		return false, nil
	}
	hosts := make(map[string]string)
	for _, record := range victim.Pods {
		uid := victim.Resources.UIDs[record.Name]
		if uid == "" {
			return false, &StartupFailure{Err: fmt.Errorf("missing submitted UID for %s", record.Name)}
		}
		if len(record.Spec.Containers) > 0 && record.Spec.Containers[0].KubeVirt {
			vm, err := cachedVM(ctx, record.Name)
			if err != nil {
				return false, err
			}
			if vm == nil {
				return false, nil
			}
			if string(vm.UID) != uid || vm.DeletionTimestamp != nil {
				return false, &StartupFailure{Err: fmt.Errorf("VM %s replaced or terminating", record.Name)}
			}
			vmi, err := cachedVMI(ctx, record.Name)
			if err != nil {
				return false, err
			}
			if vmi == nil {
				return false, nil
			}
			// 虚拟机的地址在 guest 内部,VMI 状态里的地址是 virt-launcher 网卡的地址,
			// 只能按 guest 网卡的 MAC 校验网卡是否按规划接上。
			expectedMACs := make([]string, 0, len(record.Spec.Networks))
			for _, network := range record.Spec.Networks {
				if network.Attachment.MAC != "" {
					expectedMACs = append(expectedMACs, network.Attachment.MAC)
				}
			}
			if !vmiInterfacesReady(vmi, nil, expectedMACs) {
				return false, nil
			}
			if !vm.Status.Ready {
				return false, nil
			}
			continue
		}
		pod, err := cachedPod(ctx, record.Name)
		if err != nil {
			return false, err
		}
		if pod == nil {
			return false, nil
		}
		if string(pod.UID) != uid {
			return false, &StartupFailure{Err: fmt.Errorf("pod %s replaced", record.Name)}
		}
		ready, err := podStartupComplete(pod)
		if err != nil {
			return false, &StartupFailure{Err: err}
		}
		if !ready {
			return false, nil
		}
		if err = verifyPodAddresses(pod, record.Spec.Networks); err != nil {
			return false, &StartupFailure{Err: err}
		}
		hosts[record.Name] = pod.Status.HostIP
	}
	if !victim.Spec.FrpEnabled {
		victim.Endpoints = nil
		for _, pending := range victim.Resources.NodePorts {
			endpoint := pending.Endpoint
			endpoint.IP = hosts[pending.PodName]
			if endpoint.IP == "" {
				return false, nil
			}
			victim.Endpoints = append(victim.Endpoints, endpoint)
		}
		victim.ExposedEndpoints = append(model.Endpoints(nil), victim.Endpoints...)
	}
	return true, nil
}
