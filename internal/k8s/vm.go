package k8s

import (
	"context"
	"fmt"
	"net/netip"
	"strconv"
	"strings"
	"time"

	"go.yaml.in/yaml/v4"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	v1 "kubevirt.io/api/core/v1"

	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
)

type CreateVMOptions struct {
	SubmitOnly  bool
	Name        string
	Labels      map[string]string
	Image       string
	Bootloader  string
	SecureBoot  bool
	CPUMillis   int64
	MemoryBytes int64
	UserData    model.CloudInitConfig
	Networks    []Network
}

type networkDataConfig struct {
	Version   int                            `yaml:"version"`
	Ethernets map[string]networkDataEthernet `yaml:"ethernets"`
}

type networkDataEthernet struct {
	Match     networkDataMatch   `yaml:"match"`
	SetName   string             `yaml:"set-name"`
	Addresses []string           `yaml:"addresses,omitempty"`
	Routes    []networkDataRoute `yaml:"routes,omitempty"`
}

// networkDataRoute 是 guest 内部的一条静态路由。
// guest 只有靶场网卡,没有集群网络可用,跨网段可达性完全依赖这些路由。
type networkDataRoute struct {
	Destination string `yaml:"destination"`
	Gateway     string `yaml:"gateway"`
}

type networkDataMatch struct {
	MACAddress string `yaml:"macaddress"`
}

// buildNetworkData 渲染 cloud-init 的 network-data。
// 地址必须与拓扑规划的静态地址一致,DHCP 会拿到网段里的任意地址。
func buildNetworkData(networks []Network) (string, error) {
	networkData := networkDataConfig{
		Version:   2,
		Ethernets: make(map[string]networkDataEthernet, len(networks)),
	}
	for _, network := range networks {
		ethernet := networkDataEthernet{
			Match: networkDataMatch{
				MACAddress: network.MAC,
			},
			SetName: network.Interface,
		}
		if network.CIDR != "" && network.IPv4 != "" {
			prefix, err := netip.ParsePrefix(network.CIDR)
			if err != nil {
				return "", err
			}
			ethernet.Addresses = []string{fmt.Sprintf("%s/%d", network.IPv4, prefix.Bits())}
		}
		for _, route := range network.Routes {
			ethernet.Routes = append(ethernet.Routes, networkDataRoute{
				Destination: route.Dst,
				Gateway:     route.Gw,
			})
		}
		networkData.Ethernets[network.Interface] = ethernet
	}
	data, err := yaml.Marshal(networkData)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func CreateVM(ctx context.Context, options CreateVMOptions) (*v1.VirtualMachine, model.RetVal) {
	userData, err := options.UserData.String()
	if err != nil {
		log.Logger.Warningf("Failed to render cloud-init user data: %s", err)
		return nil, model.FailRetVal(i18n.K8S.CreateError, "VirtualMachine", options.Name, err)
	}
	networkData, err := buildNetworkData(options.Networks)
	if err != nil {
		log.Logger.Warningf("Failed to render cloud-init network data: %s", err)
		return nil, model.FailRetVal(i18n.K8S.CreateError, "VirtualMachine", options.Name, err)
	}

	affinity, err := imageFailureAffinity(ctx, []string{options.Image})
	if err != nil {
		return nil, model.FailRetVal(i18n.K8S.CreateError, "VirtualMachine", options.Name, err)
	}
	annotations := make(map[string]string)
	interfaces := make([]v1.Interface, 0, len(options.Networks))
	networks := make([]v1.Network, 0, len(options.Networks))
	for _, network := range options.Networks {
		attachment := fmt.Sprintf("%s/%s", globalNamespace, network.NetAttachDef)
		interfaces = append(interfaces, v1.Interface{
			Name:   network.Interface,
			Bridge: new(v1.InterfaceBridge),
		})
		networks = append(networks, v1.Network{
			Name:   network.Interface,
			Multus: &v1.MultusNetwork{NetworkName: attachment},
		})
	}
	for key, value := range ovnAnnotations(options.Networks) {
		annotations[key] = value
	}
	bootloader := &v1.Bootloader{BIOS: &v1.BIOS{}}
	if strings.ToLower(options.Bootloader) == "efi" {
		bootloader = &v1.Bootloader{EFI: &v1.EFI{SecureBoot: &options.SecureBoot}}
	}
	limits := make(corev1.ResourceList)
	if options.CPUMillis > 0 {
		limits[corev1.ResourceCPU] = resource.MustParse(strconv.FormatInt(options.CPUMillis, 10) + "m")
	}
	if options.MemoryBytes > 0 {
		limits[corev1.ResourceMemory] = resource.MustParse(strconv.FormatInt(options.MemoryBytes, 10))
	}
	vm := &v1.VirtualMachine{
		OwnerReferences: resourceOwners(ctx),
		Name:            options.Name,
		Namespace:       globalNamespace,
		Labels:          options.Labels,
		Spec: v1.VirtualMachineSpec{
			RunStrategy: new(v1.RunStrategyAlways),
			Template: &v1.VirtualMachineInstanceTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      options.Labels,
					Annotations: annotations,
				},
				Spec: v1.VirtualMachineInstanceSpec{
					Affinity: affinity,
					Domain: v1.DomainSpec{
						Firmware: &v1.Firmware{
							Bootloader: bootloader,
						},
						Devices: v1.Devices{
							Disks: []v1.Disk{
								{
									Name:      "root",
									BootOrder: new(uint(1)),
									Disk: &v1.DiskTarget{
										Bus: v1.DiskBusVirtio,
									},
								},
								{
									Name:      "cloud-init",
									BootOrder: new(uint(2)),
									Disk: &v1.DiskTarget{
										Bus: v1.DiskBusVirtio,
									},
								},
							},
							Interfaces: interfaces,
						},
						Resources: v1.ResourceRequirements{
							Limits:   limits,
							Requests: workloadRequests(limits),
						},
					},
					Networks: networks,
					Volumes: []v1.Volume{
						{
							Name: "root",
							ContainerDisk: &v1.ContainerDiskSource{
								Image: options.Image,
							},
						},
						{
							Name: "cloud-init",
							CloudInitNoCloud: &v1.CloudInitNoCloudSource{
								UserData:    userData,
								NetworkData: networkData,
							},
						},
					},
				},
			},
		},
	}
	vm, err = virtClient.KubevirtV1().VirtualMachines(globalNamespace).Create(ctx, vm, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create VirtualMachine %s: %s", options.Name, err)
		return nil, model.FailRetVal(i18n.K8S.CreateError, "VirtualMachine", options.Name, err)
	}
	// API acceptance does not mean the VM is ready for connections.
	if options.SubmitOnly {
		return vm, model.SuccessRetVal()
	}
	err = wait.PollUntilContextCancel(ctx, time.Second, true, func(ctx context.Context) (bool, error) {
		current, getErr := cachedVM(ctx, options.Name)
		if getErr != nil {
			return false, getErr
		}
		if current == nil {
			return false, nil
		}
		if current.UID != vm.UID || current.DeletionTimestamp != nil {
			return false, fmt.Errorf("VM %s was deleted or replaced", options.Name)
		}
		vm = current
		return vm.Status.Ready, nil
	})
	if err != nil {
		return nil, model.FailRetVal(i18n.K8S.GetError, "VirtualMachine", options.Name, err)
	}
	return vm, model.SuccessRetVal()
}

func GetVM(ctx context.Context, name string) (*v1.VirtualMachine, model.RetVal) {
	vm, err := virtClient.KubevirtV1().VirtualMachines(globalNamespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, getFailed("VirtualMachine", name, err)
	}
	return vm, model.SuccessRetVal()
}

func DeleteVMCollection(ctx context.Context, labels map[string]string) model.RetVal {
	return deleteCollection("VirtualMachine", "VirtualMachines", labels, func(options metav1.ListOptions) error {
		return virtClient.KubevirtV1().VirtualMachines(globalNamespace).DeleteCollection(ctx, metav1.DeleteOptions{PropagationPolicy: new(metav1.DeletePropagationForeground)}, options)
	})
}
