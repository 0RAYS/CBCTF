package k8s

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "kubevirt.io/api/core/v1"
)

type CloudInit struct {
	UserData    string
	NetworkData string
}

type CreateVMOptions struct {
	Name          string
	Labels        map[string]string
	Annotations   map[string]string
	VMLabels      map[string]string
	VMAnnotations map[string]string
	Image         string
	CloudInit     CloudInit
}

func CreateVM(ctx context.Context, options CreateVMOptions) (*v1.VirtualMachine, model.RetVal) {
	var (
		vm  *v1.VirtualMachine
		err error
	)
	vm = &v1.VirtualMachine{
		ObjectMeta: metav1.ObjectMeta{
			Name:        options.Name,
			Namespace:   globalNamespace,
			Labels:      options.Labels,
			Annotations: options.Annotations,
		},
		Spec: v1.VirtualMachineSpec{
			RunStrategy: new(v1.RunStrategyAlways),
			Template: &v1.VirtualMachineInstanceTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      options.VMLabels,
					Annotations: options.VMAnnotations,
				},
				Spec: v1.VirtualMachineInstanceSpec{
					Domain: v1.DomainSpec{
						Firmware: &v1.Firmware{
							Bootloader: &v1.Bootloader{
								BIOS: &v1.BIOS{
									UseSerial: new(false),
								},
								EFI: &v1.EFI{
									SecureBoot: new(false),
								},
							},
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
							Interfaces: []v1.Interface{},
						},
						Resources: v1.ResourceRequirements{
							Requests: corev1.ResourceList{},
							Limits:   corev1.ResourceList{},
						},
					},
					Networks: []v1.Network{},
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
									UserData:    options.CloudInit.UserData,
									NetworkData: options.CloudInit.NetworkData,
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
