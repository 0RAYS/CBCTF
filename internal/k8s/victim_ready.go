package k8s

import (
	"context"
	"fmt"

	"CBCTF/internal/model"
)

// VictimReady Readiness is evaluated from shared caches; start workers only submit objects.
func VictimReady(ctx context.Context, victim *model.Victim) (bool, error) {
	if !victim.Resources.Submitted {
		return false, nil
	}
	hosts := make(map[string]string)
	for _, record := range victim.Pods {
		uid := victim.Resources.UIDs[record.Name]
		if uid == "" {
			return false, fmt.Errorf("missing submitted UID for %s", record.Name)
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
				return false, fmt.Errorf("VM %s replaced or terminating", record.Name)
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
			return false, fmt.Errorf("pod %s replaced", record.Name)
		}
		ready, err := podStartupComplete(pod)
		if err != nil || !ready {
			return false, err
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
