package k8s

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"fmt"
	"strconv"
	"strings"

	corev1 "k8s.io/api/core/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "kubevirt.io/api/core/v1"
)

const NetworkDataTmpl = `
version: 2
ethernets:
%s
`

const NetworkDataEthernetTmpl = `
  %s:
    match:
      macaddress: "%s"
    set-name: %s
    dhcp4: true
`

type CreateVMOptions struct {
	Name        string
	Labels      map[string]string
	Image       string
	Bootloader  string
	SecureBoot  bool
	CPUMillis   int64
	MemoryBytes int64
	UserData    string
	Networks    []Network
}

func CreateVM(ctx context.Context, options CreateVMOptions) (*v1.VirtualMachine, model.RetVal) {
	var (
		vm  *v1.VirtualMachine
		err error
	)
	vm = &v1.VirtualMachine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      options.Name,
			Namespace: globalNamespace,
			Labels:    options.Labels,
		},
		Spec: v1.VirtualMachineSpec{
			RunStrategy: new(v1.RunStrategyAlways),
			Template: &v1.VirtualMachineInstanceTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: func() map[string]string {
						annotations := make(map[string]string)
						for i, network := range options.Networks {
							if i == 0 {
								annotations["ovn.kubernetes.io/logical_switch"] = network.Subnet
								annotations["ovn.kubernetes.io/ip_address"] = network.IPv4
								annotations["ovn.kubernetes.io/mac_address"] = network.MAC
								annotations["v1.multus-cni.io/default-network"] = fmt.Sprintf("%s/%s", globalNamespace, network.NetAttachDef)
							} else {
								annotations["k8s.v1.cni.cncf.io/networks"] += fmt.Sprintf(",%s/%s", globalNamespace, network.NetAttachDef)
								annotations["k8s.v1.cni.cncf.io/networks"] = strings.Trim(annotations["k8s.v1.cni.cncf.io/networks"], ",")
							}
							annotations[fmt.Sprintf("%s.%s.ovn.kubernetes.io/logical_switch", network.NetAttachDef, globalNamespace)] = network.Subnet
							annotations[fmt.Sprintf("%s.%s.ovn.kubernetes.io/ip_address", network.NetAttachDef, globalNamespace)] = network.IPv4
							annotations[fmt.Sprintf("%s.%s.ovn.kubernetes.io/mac_address", network.NetAttachDef, globalNamespace)] = network.MAC
							// 需要出网时, 设定网关
							if network.External {
								annotations[fmt.Sprintf("%s.%s.ovn.kubernetes.io/routes", network.NetAttachDef, globalNamespace)] = fmt.Sprintf("[{\"gw\":\"%s\"}]", network.Gateway)
							}
						}
						return annotations
					}(),
				},
				Spec: v1.VirtualMachineInstanceSpec{
					Domain: v1.DomainSpec{
						Firmware: &v1.Firmware{
							Bootloader: func() *v1.Bootloader {
								if strings.ToLower(options.Bootloader) == "efi" {
									return &v1.Bootloader{
										EFI: &v1.EFI{SecureBoot: &options.SecureBoot},
									}
								}
								return &v1.Bootloader{
									BIOS: &v1.BIOS{},
								}
							}(),
						},
						Devices: v1.Devices{
							Disks: []v1.Disk{
								{
									Name:      "root",
									BootOrder: new(uint(1)),
									DiskDevice: v1.DiskDevice{
										Disk: &v1.DiskTarget{
											Bus: v1.DiskBusVirtio,
										},
									},
								},
								{
									Name:      "cloud-init",
									BootOrder: new(uint(2)),
									DiskDevice: v1.DiskDevice{
										Disk: &v1.DiskTarget{
											Bus: v1.DiskBusVirtio,
										},
									},
								},
							},
							Interfaces: func() []v1.Interface {
								interfaces := make([]v1.Interface, 0)
								for _, network := range options.Networks {
									interfaces = append(interfaces, v1.Interface{
										Name: network.Interface,
										InterfaceBindingMethod: v1.InterfaceBindingMethod{
											Bridge: new(v1.InterfaceBridge),
										},
									})
								}
								return interfaces
							}(),
						},
						Resources: v1.ResourceRequirements{
							Limits: func() corev1.ResourceList {
								limit := make(corev1.ResourceList)
								if options.CPUMillis > 0 {
									limit[corev1.ResourceCPU] = resource.MustParse(strconv.FormatInt(options.CPUMillis, 10) + "m")
								}
								if options.MemoryBytes > 0 {
									limit[corev1.ResourceMemory] = resource.MustParse(strconv.FormatInt(options.MemoryBytes, 10))
								}
								return limit
							}(),
						},
					},
					Networks: func() []v1.Network {
						networks := make([]v1.Network, 0)
						for _, network := range options.Networks {
							networks = append(networks, v1.Network{
								Name: network.Interface,
								NetworkSource: v1.NetworkSource{
									Multus: &v1.MultusNetwork{
										NetworkName: fmt.Sprintf("%s/%s", globalNamespace, network.NetAttachDef),
									},
								},
							})
						}
						return networks
					}(),
					Volumes: []v1.Volume{
						{
							Name: "root",
							VolumeSource: v1.VolumeSource{
								ContainerDisk: &v1.ContainerDiskSource{
									Image: options.Image,
								},
							},
						},
						{
							Name: "cloud-init",
							VolumeSource: v1.VolumeSource{
								CloudInitNoCloud: &v1.CloudInitNoCloudSource{
									UserData: options.UserData,
									NetworkData: func() string {
										ethernets := make([]string, 0)
										for _, network := range options.Networks {
											ethernets = append(ethernets, fmt.Sprintf(NetworkDataEthernetTmpl, network.Interface, network.MAC, network.Interface))
										}
										return fmt.Sprintf(NetworkDataTmpl, strings.Join(ethernets, "\n"))
									}(),
								},
							},
						},
					},
				},
			},
		},
	}
	vm, err = virtClient.KubevirtV1().VirtualMachines(globalNamespace).Create(ctx, vm, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create virtual machine: %s", err)
		return nil, model.RetVal{Msg: i18n.K8S.CreateError, Attr: map[string]any{"Model": "VirtualMachine", "Error": err.Error()}}
	}
	return vm, model.SuccessRetVal()
}

func GetVM(ctx context.Context, name string) (*v1.VirtualMachine, model.RetVal) {
	vm, err := virtClient.KubevirtV1().VirtualMachines(globalNamespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if apierror.IsNotFound(err) {
			return nil, model.RetVal{Msg: i18n.K8S.NotFound, Attr: map[string]any{"Model": "VirtualMachine"}}
		}
		log.Logger.Warningf("Failed to get VirtualMachine: %s", err)
		return nil, model.RetVal{Msg: i18n.K8S.GetError, Attr: map[string]any{"Model": "VirtualMachine", "Error": err.Error()}}
	}
	return vm, model.SuccessRetVal()
}

func ListVM(ctx context.Context, labels ...map[string]string) (*v1.VirtualMachineList, model.RetVal) {
	var options metav1.ListOptions
	if len(labels) > 0 {
		var selector string
		for k, v := range labels[0] {
			selector += fmt.Sprintf("%s=%s,", k, v)
		}
		options = metav1.ListOptions{
			LabelSelector: strings.TrimSuffix(selector, ","),
		}
	}
	vmList, err := virtClient.KubevirtV1().VirtualMachines(globalNamespace).List(ctx, options)
	if err != nil {
		log.Logger.Warningf("Failed to list VirtualMachines: %s", err)
		return nil, model.RetVal{Msg: i18n.K8S.GetError, Attr: map[string]any{"Model": "VirtualMachine", "Error": err.Error()}}
	}
	return vmList, model.SuccessRetVal()
}

func DeleteVM(ctx context.Context, name string) model.RetVal {
	err := virtClient.KubevirtV1().VirtualMachines(globalNamespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete VirtualMachine: %s", err)
		return model.RetVal{Msg: i18n.K8S.DeleteError, Attr: map[string]any{"Model": "VirtualMachine", "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

func DeleteVMCollection(ctx context.Context, labels ...map[string]string) model.RetVal {
	var options metav1.ListOptions
	if len(labels) > 0 {
		var selector string
		for k, v := range labels[0] {
			selector += fmt.Sprintf("%s=%s,", k, v)
		}
		options = metav1.ListOptions{
			LabelSelector: strings.TrimSuffix(selector, ","),
		}
	}
	err := virtClient.KubevirtV1().VirtualMachines(globalNamespace).DeleteCollection(ctx, metav1.DeleteOptions{}, options)
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete VirtualMachines: %s", err)
		return model.RetVal{Msg: i18n.K8S.DeleteError, Attr: map[string]any{"Model": "VirtualMachine", "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}
