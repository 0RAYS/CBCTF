package k8s

import (
	"context"
	"strings"

	"k8s.io/apimachinery/pkg/labels"
	virtv1 "kubevirt.io/api/core/v1"
)

func cachedVMI(ctx context.Context, name string) (*virtv1.VirtualMachineInstance, error) {
	c, err := cachedObjects(ctx, "vmis")
	if err != nil {
		return nil, err
	}
	obj, exists, err := c.informer.GetStore().GetByKey(globalNamespace + "/" + name)
	if err != nil || !exists {
		return nil, err
	}
	return obj.(*virtv1.VirtualMachineInstance).DeepCopy(), nil
}

func listVMIs(ctx context.Context, selector map[string]string) ([]virtv1.VirtualMachineInstance, error) {
	c, err := cachedObjects(ctx, "vmis")
	if err != nil {
		return nil, err
	}
	items := c.informer.GetStore().List()
	result := make([]virtv1.VirtualMachineInstance, 0, len(items))
	for _, item := range items {
		vmi := item.(*virtv1.VirtualMachineInstance)
		if !labels.SelectorFromSet(selector).Matches(labels.Set(vmi.Labels)) {
			continue
		}
		result = append(result, *vmi)
	}
	return result, nil
}

// vmiInterfacesReady 判断实例是否运行,且每张规划内的网卡都接上了。
// 实例起来了但只挂上一张网卡是多网卡拓扑里最常见的故障。
// 虚拟机的地址配在 guest 内部,VMI 状态里的地址属于 virt-launcher 的附加网卡,
// 虚拟机只能按 guest 网卡的 MAC 校验。
func vmiInterfacesReady(vmi *virtv1.VirtualMachineInstance, expectedIPs, expectedMACs []string) bool {
	if vmi == nil || vmi.Status.Phase != virtv1.Running {
		return false
	}
	if len(expectedIPs) == 0 && len(expectedMACs) == 0 {
		return true
	}
	addresses := make(map[string]bool, len(vmi.Status.Interfaces))
	macs := make(map[string]bool, len(vmi.Status.Interfaces))
	for _, iface := range vmi.Status.Interfaces {
		if iface.IP != "" {
			addresses[iface.IP] = true
		}
		if iface.MAC != "" {
			macs[strings.ToLower(iface.MAC)] = true
		}
	}
	for _, ip := range expectedIPs {
		if !addresses[ip] {
			return false
		}
	}
	for _, mac := range expectedMACs {
		if !macs[strings.ToLower(mac)] {
			return false
		}
	}
	return true
}
